#!/bin/bash

VERSION="0.0.2"

echo "Building Linux..."
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.buildVersion=$VERSION" -o "build/patreon-dl_${VERSION}_Linux_amd64"
GOOS=linux GOARCH=386 go build -ldflags "-X main.buildVersion=$VERSION" -o "build/patreon-dl_${VERSION}_Linux_i386"

echo "Building Windows..."
GOOS=windows GOARCH=386 go build -ldflags "-X main.buildVersion=$VERSION" -o "build/patreon-dl_${VERSION}_Windows_x86_64.exe"
GOOS=windows GOARCH=386 go build -ldflags "-X main.buildVersion=$VERSION" -o "build/patreon-dl_${VERSION}_Windows_i386.exe"

echo "Building Darwin..."
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.buildVersion=$VERSION" -o "build/patreon-dl_${VERSION}_Darwin_x86_64"
GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.buildVersion=$VERSION" -o "build/patreon-dl_${VERSION}_Darwin_ARM64"

echo "Creating release signatures..."
cd build || exit
gpg --detach-sign --armor --local-user privateger "patreon-dl_${VERSION}_Linux_amd64"
gpg --detach-sign --armor --local-user privateger "patreon-dl_${VERSION}_Linux_i386"
gpg --detach-sign --armor --local-user privateger "patreon-dl_${VERSION}_Windows_x86_64.exe"
gpg --detach-sign --armor --local-user privateger "patreon-dl_${VERSION}_Windows_i386.exe"
gpg --detach-sign --armor --local-user privateger "patreon-dl_${VERSION}_Darwin_x86_64"
gpg --detach-sign --armor --local-user privateger "patreon-dl_${VERSION}_Darwin_ARM64"

