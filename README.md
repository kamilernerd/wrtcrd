# wrtcrd
Remote desktop server and client using Webrtc protocol.

## Features
- Cross platform server (windows, linux, darwin) and client (web)
- Multiple display support
- Audio output (Opus)
- Mouse input
- Keyboard input
- Webrtc 60fps (x264)

## TODO
- Actually implement the mouse and keyboard support
- Look for performance improvements
- Add support for other codecs like AAC, Vorbis for audio
- Add support for other codecs like VP8/9 and AV1???
- Logging
- Some kind of authentication
- General server improvements
- SFU????
- Very important thing is to go away from using libraries that do most of the magic behind and create own implementations.

## How to setup
This project is not really tested on other platforms than darwin using arm64 so
a lot of things are a subject to change and might still not work as expected on
other platforms so keep that in mind...

Well this is a bit weird... You need to install x264 and Opus codecs somehow
(FFMPEG?).
Depending on your platform things will be a bit different.

### Darwin
It turns out that there is no simple way around CoreAudio to create a virtual
audio output device so we could send audio out. Therfore you will need to
install BlackHole and set it up to output to your speakers/other device and
blackhole. Then set input device to blackhole as well. Now sound should work.

Install libportaudio using brew

That should be everything you need to run this.

Not install the dependencies and run this

```GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go run main.go```

There a bigger list of dependencies you should have but they are in the readme
files of the dependencies.

### WIP
Everything actually :)

PS: Respect the licenses
