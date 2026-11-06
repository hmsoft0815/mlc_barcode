param(
  [string]$Apk = "$env:USERPROFILE\mlcb-test\app.apk",
  [string]$Avd = "mlcb-test",
  [string]$Package = "eu.mlcgo.mlcbarcode",
  [string]$Out = "$env:USERPROFILE\mlcb-test",
  [switch]$Stop
)
$sdk = "$env:LOCALAPPDATA\Android\Sdk"
$adb = "$sdk\platform-tools\adb.exe"
$emu = "$sdk\emulator\emulator.exe"
function Fail($msg) { Write-Output "FAIL: $msg"; exit 1 }

# The serial of the emulator running $Avd — other emulators may be running.
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
$serial = Find-Serial
if ($Stop) {
  if ($serial) { & $adb -s $serial emu kill | Out-Null; Write-Output "emulator $Avd stopped" }
  else { Write-Output "emulator $Avd not running" }
  exit 0
}

if (-not $serial) {
  # Win32_Process.Create (survives the SSH session ending) inherits the WMI
  # service's environment: without the user's paths the emulator sends the
  # wrong adb key and finds no AVDs.
  $envs = "set USERPROFILE=$env:USERPROFILE&& set HOMEDRIVE=$env:HOMEDRIVE&& set HOMEPATH=$env:HOMEPATH&& set LOCALAPPDATA=$env:LOCALAPPDATA&& set APPDATA=$env:APPDATA&& set TEMP=$env:TEMP&& set TMP=$env:TMP"
  $cmd = "cmd /c `"$envs&& `"$emu`" -avd $Avd -no-window -no-audio -no-boot-anim -no-snapshot-save -gpu swiftshader_indirect -camera-back emulated > `"$Out\emulator.log`" 2>&1`""
  $r = Invoke-CimMethod -ClassName Win32_Process -MethodName Create -Arguments @{ CommandLine = $cmd }
  if ($r.ReturnValue -ne 0) { Fail "emulator start ($($r.ReturnValue))" }
  Write-Output "emulator starting ($Avd)"
}

$booted = $false
for ($i = 0; $i -lt 120; $i++) {
  if (-not $serial) { $serial = Find-Serial }
  if ($serial -and ((& $adb -s $serial shell getprop sys.boot_completed 2>$null) -match "1")) { $booted = $true; break }
  Start-Sleep 2
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
Write-Output "OK: app running (pid $appPid), screenshot $Out\screen.png"
