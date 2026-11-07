param(
  [string]$Apk = "$env:USERPROFILE\mlcb-test\app.apk",
  [string]$Avd = "mlcb-test",
  [string]$Package = "eu.mlcgo.mlcbarcode",
  [string]$Out = "$env:USERPROFILE\mlcb-test",
  [string]$CameraImage = "",
  [string]$ExpectText = "",
  [switch]$Stop,
  [switch]$Docs
)
$sdk = "$env:LOCALAPPDATA\Android\Sdk"
$adb = "$sdk\platform-tools\adb.exe"
$emu = "$sdk\emulator\emulator.exe"
function Fail($msg) { Write-Output "FAIL: $msg"; exit 1 }

# The serial of the emulator running $Avd - other emulators may be running.
function Find-Serial {
  foreach ($line in (& $adb devices)) {
    if ($line -match "^(emulator-\d+)\s+device") {
      $s = $Matches[1]
      $name = (& $adb -s $s emu avd name 2>$null | Select-Object -First 1)
      if ($name -and $name.Trim() -eq $Avd) { return $s }
    }
  }
  return $null
}

& $adb start-server 2>$null | Out-Null
# A fresh adb server sees running emulators as "offline" for a few seconds;
# deciding then would start a second instance of the AVD.
for ($i = 0; $i -lt 15; $i++) {
  if (-not ((& $adb devices) -match "\soffline$")) { break }
  Start-Sleep 1
}
$serial = Find-Serial
# The back camera shows -CameraImage (a code to scan) or the emulator's test
# pattern; it is fixed at emulator start, so a different camera restarts it.
$camera = if ($CameraImage) { "imagefile:$CameraImage" } else { "emulated" }
$camFile = "$Out\camera.txt"
# The image's timestamp is part of the key: a changed file needs a restart too.
$cameraKey = if ($CameraImage) { "$camera @ " + (Get-Item $CameraImage).LastWriteTimeUtc.Ticks } else { $camera }
$lastCamera = if (Test-Path $camFile) { Get-Content $camFile } else { "" }
# qemu processes of this AVD (the emulator.exe launcher may linger; it is harmless).
function Emu-Running {
  [bool](Get-CimInstance Win32_Process -Filter "Name LIKE 'qemu-system%'" |
    Where-Object { $_.CommandLine -match "-avd $Avd(\s|$)" })
}
function Stop-Emu {
  if ($serial) { & $adb -s $serial emu kill | Out-Null }
  # Wait for the process, not just for adb: a new instance started while the
  # old one is still shutting down loses its adb connection and quits.
  for ($i = 0; $i -lt 60 -and (Emu-Running); $i++) { Start-Sleep 1 }
}
function Start-Emu {
  Set-Content $camFile $cameraKey
  # Win32_Process.Create (survives the SSH session ending) inherits the WMI
  # service's environment: without the user's paths the emulator sends the
  # wrong adb key and finds no AVDs.
  $envs = "set USERPROFILE=$env:USERPROFILE&& set HOMEDRIVE=$env:HOMEDRIVE&& set HOMEPATH=$env:HOMEPATH&& set LOCALAPPDATA=$env:LOCALAPPDATA&& set APPDATA=$env:APPDATA&& set TEMP=$env:TEMP&& set TMP=$env:TMP"
  $cmd = "cmd /c `"$envs&& `"$emu`" -avd $Avd -no-window -no-audio -no-boot-anim -no-snapshot-save -gpu swiftshader_indirect -camera-back $camera > `"$Out\emulator.log`" 2>&1`""
  $r = Invoke-CimMethod -ClassName Win32_Process -MethodName Create -Arguments @{ CommandLine = $cmd }
  if ($r.ReturnValue -ne 0) { Fail "emulator start ($($r.ReturnValue))" }
  Write-Output "emulator starting ($Avd, camera $camera)"
}
function Wait-Boot([int]$seconds) {
  for ($i = 0; $i -lt $seconds / 2; $i++) {
    if (-not $script:serial) { $script:serial = Find-Serial }
    if ($script:serial -and ((& $adb -s $script:serial shell getprop sys.boot_completed 2>$null) -match "1")) { return $true }
    # Started and gone again (lost its adb connection during boot): give up early.
    if ($i -gt 15 -and -not (Emu-Running)) { return $false }
    Start-Sleep 2
  }
  return $false
}

if ($Stop) {
  if ($serial -or (Emu-Running)) { Stop-Emu; Write-Output "emulator $Avd stopped" }
  else { Write-Output "emulator $Avd not running" }
  exit 0
}
if ($serial -and $lastCamera -ne $cameraKey) { Stop-Emu; $serial = $null }
if (-not $serial -and (Emu-Running)) { Stop-Emu }   # running but not reachable

$booted = $false
if ($serial) { $booted = Wait-Boot 240 }
else {
  for ($attempt = 1; $attempt -le 2 -and -not $booted; $attempt++) {
    Start-Emu
    $booted = Wait-Boot 240
    if (-not $booted) { $serial = $null; Stop-Emu }
  }
}
if (-not $booted) { & $adb devices; Fail "emulator $Avd did not boot" }
Write-Output ("emulator $Avd ready ($serial), API " + (& $adb -s $serial shell getprop ro.build.version.sdk))

# Fresh install: no leftover data, no signature clash with older builds.
& $adb -s $serial uninstall $Package 2>$null | Out-Null
$inst = & $adb -s $serial install $Apk 2>&1
if ($LASTEXITCODE -ne 0) { Fail "install: $inst" }
& $adb -s $serial shell pm grant $Package android.permission.CAMERA
& $adb -s $serial logcat -c
& $adb -s $serial shell am start -W -n "$Package/com.wails.app.MainActivity" | Select-String "Status|TotalTime"
Start-Sleep 10

$crash = & $adb -s $serial logcat -d | Select-String -Pattern "FATAL EXCEPTION|AndroidRuntime: |Fatal signal|E/WailsBridge"
$appPid = & $adb -s $serial shell pidof $Package
& $adb -s $serial shell screencap -p /sdcard/mlcb.png
& $adb -s $serial pull /sdcard/mlcb.png "$Out\screen.png" 2>$null | Out-Null
& $adb -s $serial logcat -d > "$Out\logcat.txt"
if ($crash) { $crash | Select-Object -First 20; Fail "crash in logcat" }
if (-not $appPid) { Fail "app not running" }

# --- UI helpers: uiautomator sees the WebView's text once accessibility is up.
function Get-UI {
  & $adb -s $serial shell uiautomator dump /sdcard/mlcb-ui.xml 2>$null | Out-Null
  & $adb -s $serial pull /sdcard/mlcb-ui.xml "$Out\ui.xml" 2>$null | Out-Null
  [xml](Get-Content -Raw -Encoding UTF8 "$Out\ui.xml")
}
function Find-Node([string]$pattern, [string]$class = "") {
  for ($try = 0; $try -lt 10; $try++) {
    $ui = Get-UI
    foreach ($n in $ui.SelectNodes("//node")) {
      $label = ($n.text + " " + $n.'content-desc').Trim()
      if ($class -and $n.class -ne $class) { continue }
      if ($pattern -and $label -notmatch $pattern) { continue }
      if ($n.bounds -match '\[(\d+),(\d+)\]\[(\d+),(\d+)\]') {
        return @{ X = ([int]$Matches[1] + [int]$Matches[3]) / 2; Y = ([int]$Matches[2] + [int]$Matches[4]) / 2; Label = $label }
      }
    }
    Start-Sleep 1
  }
  return $null
}
function Tap([string]$pattern, [string]$class = "") {
  $n = Find-Node $pattern $class
  if (-not $n) { Fail "not found on screen: $pattern $class" }
  & $adb -s $serial shell input tap ([int]$n.X) ([int]$n.Y)
  Start-Sleep 2
}
function Shot([string]$name) {
  & $adb -s $serial shell screencap -p /sdcard/mlcb-shot.png
  & $adb -s $serial pull /sdcard/mlcb-shot.png "$Out\docs\$name.png" 2>$null | Out-Null
  Write-Output "  $name"
}

# With a code in front of the camera the app must have read it.
if ($ExpectText) {
  if (-not (Find-Node $ExpectText)) { Fail "scan result missing (expected '$ExpectText')" }
  Write-Output "scan result shows '$ExpectText'"
}

if ($Docs) {
  # Screenshots for the README, per language: scan result, history,
  # create (with a QR code), about. Tab labels in both languages.
  New-Item -ItemType Directory -Force "$Out\docs" | Out-Null
  Remove-Item "$Out\docs\*.png" -ErrorAction SilentlyContinue
  $L = @{
    en = @{ Scan = 'Scan$'; History = 'History$'; Create = 'Create$'; About = 'About$'; Lang = '^English$'; Make = 'Create QR code' }
    de = @{ Scan = 'Scannen$'; History = 'Verlauf$'; Create = 'Erzeugen$'; About = '.ber$'; Lang = '^Deutsch$'; Make = 'QR-Code erzeugen' }
  }
  foreach ($lang in @("en", "de")) {
    $t = $L[$lang]
    $about = if ($lang -eq "en") { $L.de.About + '|' + $L.en.About } else { $L.en.About + '|' + $L.de.About }
    Tap $about "android.widget.Button"
    Tap $t.Lang
    Tap $t.Scan "android.widget.Button"
    # A fresh scan, so the camera picture shows behind the result.
    Tap "Scan again|Weiter scannen"
    Start-Sleep 5
    Shot "android-$lang-scan"
    Tap $t.History "android.widget.Button"
    Shot "android-$lang-history"
    Tap $t.Create "android.widget.Button"
    Tap "" "android.widget.EditText"
    & $adb -s $serial shell input text "https://mlcgo.eu/products/mlc-barcode/"
    & $adb -s $serial shell input keyevent 111   # ESC: hide the keyboard
    Start-Sleep 1
    Tap $t.Make
    Start-Sleep 2
    Shot "android-$lang-create"
    Tap $t.About "android.widget.Button"
    Shot "android-$lang-about"
  }
  Write-Output "docs screenshots: $Out\docs"
}
Write-Output "OK: app running (pid $appPid), screenshot $Out\screen.png"
