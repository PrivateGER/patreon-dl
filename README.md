# patreon-dl: Patreon Image Downloader
[![Downloads](https://img.shields.io/github/downloads/PrivateGER/patreon-dl/total?style=for-the-badge&color=blue)](https://github.com/PrivateGER/patreon-dl/releases/latest)
[![AUR Package](https://img.shields.io/aur/version/patreon-dl?style=for-the-badge&color=blue)](https://aur.archlinux.org/packages/patreon-dl/)

`patreon-dl` lets you download images of creators you support.

# DEPRECATED.
Use [gallery-dl](https://github.com/mikf/gallery-dl) for the same functionality (and more).

## Installation

**If you are running Arch:** patreon-dl is available as [AUR package!](https://aur.archlinux.org/packages/patreon-dl/)

Otherwise:
Download a release from https://github.com/PrivateGER/patreon-dl/releases.
  - Linux (64bit): patreon-dl_x.x.x_Linux_x86_64
  - Linux (32bit): patreon-dl_x.x.x_Linux_i386
  - Linux (ARM64): patreon-dl_x.x.x_Linux_arm64
  - Windows (64bit): patreon-dl_x.x.x_Windows_x86_64.exe 
  - Windows (32bit): patreon-dl_x.x.x_Windows_i386.exe
  - Mac (x64): patreon-dl_x.x.x_Darwin_x86_64
  - Mac (ARM): patreon-dl_x.x.x_Darwin_arm64

The Mac builds are unsupported. I don't own a Mac and don't have access to one either. If someone finds a way to get it working, please open an issue and I'll add it to the README.

You can also compile patreon-dl yourself using `go build` or build for all operating systems using [gorelaser](https://github.com/goreleaser/goreleaser) with the included configuration file.

## Usage

Run `patreon-dl` and follow the instructions.

### Example:

Open https://patreon.com/creatornamehere/posts in your browser. Now start `patreon-dl` and open the browser console. Paste the line of code shown by `patreon-dl` into the console and execute it by pressing ENTER. 

After a few seconds of loading, depending on how many posts the creator has, `patreon-dl` will start downloading all images into the `images` folder.

## Security

`patreon-dl` is not distinguishable from normal use for Patreon. There is no risk of getting banned or punished for the use of this tool.

Release tags of `patreon-dl` after v1.0.2 are signed with a PGP key with the fingerprint `CAE625C962F94C67`.
