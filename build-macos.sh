#!/bin/bash
set -e

APP_NAME="Totally Goated"
BUNDLE_ID="com.totallygoated.app"
EXECUTABLE="totally-goated"
APP_DIR="build/${APP_NAME}.app"

echo "=== Building Totally Goated for macOS ==="

# Clean previous build
rm -rf "build/${APP_NAME}.app"
mkdir -p "${APP_DIR}/Contents/MacOS"
mkdir -p "${APP_DIR}/Contents/Resources"

# Build the Go binary for the current architecture
echo "Building for $(uname -m)..."
CGO_ENABLED=1 go build -o "${APP_DIR}/Contents/MacOS/${EXECUTABLE}" .

# Copy Info.plist
cp Info.plist "${APP_DIR}/Contents/"

# Generate .icns icon from goat.png
echo "Generating app icon..."
ICONSET_DIR="build/AppIcon.iconset"
mkdir -p "${ICONSET_DIR}"

sips -z 16 16     assets/textures/goat.png --out "${ICONSET_DIR}/icon_16x16.png"      > /dev/null 2>&1
sips -z 32 32     assets/textures/goat.png --out "${ICONSET_DIR}/icon_16x16@2x.png"   > /dev/null 2>&1
sips -z 32 32     assets/textures/goat.png --out "${ICONSET_DIR}/icon_32x32.png"      > /dev/null 2>&1
sips -z 64 64     assets/textures/goat.png --out "${ICONSET_DIR}/icon_32x32@2x.png"   > /dev/null 2>&1
sips -z 128 128   assets/textures/goat.png --out "${ICONSET_DIR}/icon_128x128.png"    > /dev/null 2>&1
sips -z 256 256   assets/textures/goat.png --out "${ICONSET_DIR}/icon_128x128@2x.png" > /dev/null 2>&1
sips -z 256 256   assets/textures/goat.png --out "${ICONSET_DIR}/icon_256x256.png"    > /dev/null 2>&1
sips -z 512 512   assets/textures/goat.png --out "${ICONSET_DIR}/icon_256x256@2x.png" > /dev/null 2>&1
sips -z 512 512   assets/textures/goat.png --out "${ICONSET_DIR}/icon_512x512.png"    > /dev/null 2>&1
sips -z 1024 1024 assets/textures/goat.png --out "${ICONSET_DIR}/icon_512x512@2x.png" > /dev/null 2>&1

iconutil -c icns "${ICONSET_DIR}" -o "${APP_DIR}/Contents/Resources/AppIcon.icns"
rm -rf "${ICONSET_DIR}"

# Ad-hoc code sign so macOS Gatekeeper doesn't flag the app as "damaged"
echo "Code signing..."
codesign --force --deep --sign - "${APP_DIR}"

echo ""
echo "=== Build complete! ==="
echo "App bundle: ${APP_DIR}"
echo ""
echo "You can run it with:  open \"${APP_DIR}\""
echo ""
echo "If the app shows 'damaged' on another Mac, run this on that Mac:"
echo "  xattr -cr \"${APP_DIR}\""
