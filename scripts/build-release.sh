#!/usr/bin/env bash
# Build all GUI release artifacts (macOS DMG, Windows GUI exe + NSIS setup + MSI).
# Runs on macOS with: mingw-w64, makensis, msitools (wixl) installed.
# Usage: scripts/build-release.sh <version>   e.g. 1.0.1
set -euo pipefail

V="${1:?version required, e.g. 1.0.1}"
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
export PATH="/opt/homebrew/bin:/usr/local/bin:$PATH"
LDFLAGS="-s -w -X ncp-nuke/pkg/version.Version=${V}"

rm -rf dist && mkdir -p dist

echo "[1/4] macOS universal app + DMG"
CGO_ENABLED=1 GOARCH=arm64 go build -ldflags "$LDFLAGS" -o /tmp/d-arm64 ./desktop
CGO_ENABLED=1 GOARCH=amd64 go build -ldflags "$LDFLAGS" -o /tmp/d-amd64 ./desktop
lipo -create -output /tmp/d-univ /tmp/d-arm64 /tmp/d-amd64
APP="dist/NCP Nuke.app"
mkdir -p "$APP/Contents/MacOS"
cp /tmp/d-univ "$APP/Contents/MacOS/ncp-nuke-desktop"; chmod +x "$APP/Contents/MacOS/ncp-nuke-desktop"
cat > "$APP/Contents/Info.plist" <<PLIST
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
<key>CFBundleName</key><string>NCP Nuke</string>
<key>CFBundleDisplayName</key><string>NCP Nuke</string>
<key>CFBundleIdentifier</key><string>com.tbit.ncp-nuke</string>
<key>CFBundleVersion</key><string>${V}</string>
<key>CFBundleShortVersionString</key><string>${V}</string>
<key>CFBundlePackageType</key><string>APPL</string>
<key>CFBundleExecutable</key><string>ncp-nuke-desktop</string>
<key>LSMinimumSystemVersion</key><string>10.15</string>
<key>LSApplicationCategoryType</key><string>public.app-category.developer-tools</string>
<key>NSHighResolutionCapable</key><true/>
<key>NSHumanReadableCopyright</key><string>(c) 2026 TBIT</string>
<key>NSAppTransportSecurity</key><dict><key>NSAllowsLocalNetworking</key><true/></dict>
</dict></plist>
PLIST
codesign --force --deep --sign - "$APP" >/dev/null 2>&1 || true
STAGE="$(mktemp -d)"; cp -R "$APP" "$STAGE/"; ln -s /Applications "$STAGE/Applications"
hdiutil create -volname "NCP Nuke" -srcfolder "$STAGE" -ov -format UDZO "dist/NCP-Nuke_${V}_macOS_universal.dmg" >/dev/null
rm -rf "$STAGE" "$APP"

echo "[2/4] Windows GUI exe (mingw cross-compile)"
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ \
  go build -ldflags "-H windowsgui $LDFLAGS" -o "dist/NCP-Nuke_${V}_windows_amd64.exe" ./desktop

echo "[3/4] NSIS setup.exe"
NS="$(mktemp -d)"; cp "dist/NCP-Nuke_${V}_windows_amd64.exe" "$NS/NCP-Nuke.exe"
cat > "$NS/installer.nsi" <<'NSI'
Unicode true
!include "MUI2.nsh"
Name "NCP Nuke"
OutFile "ncp-setup.exe"
InstallDir "$PROGRAMFILES64\NCP Nuke"
InstallDirRegKey HKLM "Software\TBIT\NCPNuke" "InstallDir"
RequestExecutionLevel admin
!define MUI_ABORTWARNING
!define MUI_FINISHPAGE_RUN "$INSTDIR\NCP-Nuke.exe"
!define MUI_FINISHPAGE_RUN_TEXT "NCP Nuke 실행"
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_LANGUAGE "Korean"
!insertmacro MUI_LANGUAGE "English"
!define UNINSTKEY "Software\Microsoft\Windows\CurrentVersion\Uninstall\NCPNuke"
Section "NCP Nuke"
  SetOutPath "$INSTDIR"
  File "NCP-Nuke.exe"
  CreateDirectory "$SMPROGRAMS\NCP Nuke"
  CreateShortcut "$SMPROGRAMS\NCP Nuke\NCP Nuke.lnk" "$INSTDIR\NCP-Nuke.exe"
  WriteUninstaller "$INSTDIR\uninstall.exe"
  WriteRegStr HKLM "Software\TBIT\NCPNuke" "InstallDir" "$INSTDIR"
  WriteRegStr HKLM "${UNINSTKEY}" "DisplayName" "NCP Nuke"
  WriteRegStr HKLM "${UNINSTKEY}" "Publisher" "TBIT"
  WriteRegStr HKLM "${UNINSTKEY}" "DisplayIcon" "$INSTDIR\NCP-Nuke.exe"
  WriteRegStr HKLM "${UNINSTKEY}" "UninstallString" "$INSTDIR\uninstall.exe"
  WriteRegDWORD HKLM "${UNINSTKEY}" "NoModify" 1
  WriteRegDWORD HKLM "${UNINSTKEY}" "NoRepair" 1
SectionEnd
Section "Uninstall"
  Delete "$INSTDIR\NCP-Nuke.exe"
  Delete "$INSTDIR\uninstall.exe"
  Delete "$SMPROGRAMS\NCP Nuke\NCP Nuke.lnk"
  RMDir "$SMPROGRAMS\NCP Nuke"
  RMDir "$INSTDIR"
  DeleteRegKey HKLM "${UNINSTKEY}"
  DeleteRegKey HKLM "Software\TBIT\NCPNuke"
SectionEnd
NSI
( cd "$NS" && makensis -V2 -DVERSION="$V" installer.nsi >/dev/null )
cp "$NS/ncp-setup.exe" "dist/NCP-Nuke_${V}_windows_amd64_setup.exe"
rm -rf "$NS"

echo "[4/4] MSI (wixl)"
MSI="$(mktemp -d)"; cp "dist/NCP-Nuke_${V}_windows_amd64.exe" "$MSI/NCP-Nuke.exe"
PG="$(uuidgen)"; UG="$(uuidgen)"; CG="$(uuidgen)"; SG="$(uuidgen)"
cat > "$MSI/m.wxs" <<EOF
<?xml version="1.0" encoding="utf-8"?>
<Wix xmlns="http://schemas.microsoft.com/wix/2006/wi">
  <Product Name="NCP Nuke" Id="$PG" UpgradeCode="$UG" Language="1033" Version="$V" Manufacturer="TBIT">
    <Package InstallerVersion="200" Compressed="yes" Comments="NCP Nuke Installer" InstallScope="perMachine"/>
    <MajorUpgrade DowngradeErrorMessage="A newer version is already installed."/>
    <Media Id="1" Cabinet="app.cab" EmbedCab="yes"/>
    <Directory Id="TARGETDIR" Name="SourceDir"><Directory Id="ProgramFiles64Folder"><Directory Id="INSTALLDIR" Name="NCP Nuke">
      <Component Id="MainExe" Guid="$CG" Win64="yes"><File Id="ncpnukeexe" Source="NCP-Nuke.exe" KeyPath="yes"/></Component>
    </Directory></Directory>
    <Directory Id="ProgramMenuFolder"><Component Id="StartMenu" Guid="$SG" Win64="yes">
      <Shortcut Id="sc" Name="NCP Nuke" Target="[INSTALLDIR]NCP-Nuke.exe" WorkingDirectory="INSTALLDIR"/>
      <RegistryValue Root="HKCU" Key="Software\\TBIT\\NCPNuke" Name="installed" Type="integer" Value="1" KeyPath="yes"/>
    </Component></Directory></Directory>
    <Feature Id="Main" Title="NCP Nuke" Level="1"><ComponentRef Id="MainExe"/><ComponentRef Id="StartMenu"/></Feature>
  </Product>
</Wix>
EOF
( cd "$MSI" && wixl --arch x64 -o m.msi m.wxs >/dev/null )
cp "$MSI/m.msi" "dist/NCP-Nuke_${V}_windows_amd64.msi"
rm -rf "$MSI"

( cd dist && shasum -a 256 NCP-Nuke_* > checksums.txt )
echo "DONE:"; ls -1 dist
