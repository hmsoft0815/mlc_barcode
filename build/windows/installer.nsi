; ─── MLC Barcode — NSIS Installer ───────────────────────────────────────────
; Requires NSIS 3.x with MUI2
;
; Build from project root (Linux cross-compile):
;   makensis -V3 -NOCD \
;     -DAPP_VERSION=1.3.0 \
;     -DEXE_SRC=bin/mlc-barcode-gui-windows-amd64.exe \
;     build/windows/installer.nsi
; ───────────────────────────────────────────────────────────────────────────────

!ifndef APP_VERSION
  !define APP_VERSION "1.3.0"
!endif

!ifndef EXE_SRC
  !define EXE_SRC "bin/mlc-barcode-gui-windows-amd64.exe"
!endif

; ─── General ──────────────────────────────────────────────────────────────────

Name               "MLC Barcode"
OutFile            "bin/mlc-barcode-${APP_VERSION}-windows-setup.exe"
InstallDir         "$PROGRAMFILES64\MLC Barcode"
InstallDirRegKey   HKLM "Software\MLC Barcode" "InstallDir"
RequestExecutionLevel admin
Unicode            true
SetCompressor      /SOLID lzma
SetCompressorDictSize 32

; ─── MUI2 ─────────────────────────────────────────────────────────────────────

!include "MUI2.nsh"
!include "LogicLib.nsh"
!include "FileFunc.nsh"

!define MUI_ABORTWARNING
!define MUI_ICON        "build/windows/icon.ico"
!define MUI_UNICON      "build/windows/icon.ico"

; Finish page — offer to launch the app
!define MUI_FINISHPAGE_RUN         "$INSTDIR\mlc-barcode.exe"
!define MUI_FINISHPAGE_RUN_TEXT    "$(LAUNCH_APP)"

; ─── Pages ────────────────────────────────────────────────────────────────────

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

; ─── Languages ────────────────────────────────────────────────────────────────

!insertmacro MUI_LANGUAGE "English"
!insertmacro MUI_LANGUAGE "German"

; ─── Translated strings ───────────────────────────────────────────────────────

LangString LAUNCH_APP         ${LANG_ENGLISH}        "Launch MLC Barcode"
LangString LAUNCH_APP         ${LANG_GERMAN}         "MLC Barcode starten"

LangString UNINSTALL_CONFIRM  ${LANG_ENGLISH}        "Remove MLC Barcode and all its components?"
LangString UNINSTALL_CONFIRM  ${LANG_GERMAN}         "MLC Barcode und alle Komponenten entfernen?"

; ─── Installation Section ─────────────────────────────────────────────────────

Section "MainSection" SEC01
  SetOutPath "$INSTDIR"
  SetOverwrite on

  ; Install GUI, CLI, and MCP Server binaries
  File "/oname=mlc-barcode.exe" "${EXE_SRC}"
  File "/oname=barcode.exe" "bin/barcode-windows-amd64.exe"
  File "/oname=mcp-barcode-server.exe" "bin/mcp-barcode-server-windows-amd64.exe"
  File "build/windows/icon.ico"

  ; Create uninstaller
  WriteUninstaller "$INSTDIR\Uninstall.exe"

  ; Registry entries for Add/Remove Programs
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\MLC Barcode" \
    "DisplayName" "MLC Barcode"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\MLC Barcode" \
    "UninstallString" '"$INSTDIR\Uninstall.exe"'
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\MLC Barcode" \
    "DisplayIcon" '"$INSTDIR\mlc-barcode.exe",0'
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\MLC Barcode" \
    "DisplayVersion" "${APP_VERSION}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\MLC Barcode" \
    "Publisher" "Michael Lechner"
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\MLC Barcode" \
    "NoModify" 1
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\MLC Barcode" \
    "NoRepair" 1

  ; Start Menu Shortcuts
  CreateDirectory "$SMPROGRAMS\MLC Barcode"
  CreateShortcut "$SMPROGRAMS\MLC Barcode\MLC Barcode.lnk" "$INSTDIR\mlc-barcode.exe" "" "$INSTDIR\icon.ico" 0
  CreateShortcut "$SMPROGRAMS\MLC Barcode\Uninstall.lnk" "$INSTDIR\Uninstall.exe"

  ; Desktop Shortcut
  CreateShortcut "$DESKTOP\MLC Barcode.lnk" "$INSTDIR\mlc-barcode.exe" "" "$INSTDIR\icon.ico" 0
SectionEnd

; ─── Uninstallation Section ───────────────────────────────────────────────────

Section "Uninstall"
  Delete "$DESKTOP\MLC Barcode.lnk"
  Delete "$SMPROGRAMS\MLC Barcode\MLC Barcode.lnk"
  Delete "$SMPROGRAMS\MLC Barcode\Uninstall.lnk"
  RMDir "$SMPROGRAMS\MLC Barcode"

  Delete "$INSTDIR\mlc-barcode.exe"
  Delete "$INSTDIR\barcode.exe"
  Delete "$INSTDIR\mcp-barcode-server.exe"
  Delete "$INSTDIR\icon.ico"
  Delete "$INSTDIR\Uninstall.exe"
  RMDir "$INSTDIR"

  DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\MLC Barcode"
  DeleteRegKey HKLM "Software\MLC Barcode"
SectionEnd
