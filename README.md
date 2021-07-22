# patreon-dl: Patreon Image Downloader
`patreon-dl` lets you download images of creators you support.

## Installation

Download a release from https://github.com/PrivateGER/patreon-dl/releases.
  - Linux (64bit): patreon-dl_x.x.x_Linux_amd64
  - Linux (32bit): patreon-dl_x.x.x_Linux_i386
  - Windows (64bit): patreon-dl_x.x.x_Windows_x86_64.exe 
  - Windows (32bit): patreon-dl_x.x.x_Windows_i386.exe
  - Mac (x64): patreon-dl_x.x.x_Darwin_x86_64
  - Mac (ARM): patreon-dl_x.x.x_Darwin_ARM64

You can also compile patreon-dl yourself using `go build` or build for all operating systems using `compile_allarch.sh`.

## Usage

Run `patreon-dl` and follow the instructions.

### Example:

Open https://patreon.com/creatornamehere/posts in your browser. Now start `patreon-dl` and open the browser console. Paste the line of code shown by `patreon-dl` into the console and execute it by pressing ENTER. 

After a few seconds of loading, depending on how many posts the creator has, `patreon-dl` will start downloading all images into the `images` folder.

## Security

`patreon-dl` is not distinguishable from normal use for Patreon. There is no risk of getting banned or punished for the use of this tool.

All releases of `patreon-dl` are signed with a PGP key with the fingerprint `CAE625C962F94C67` and can be verified with `gpg --verify <name>.asc`. 
