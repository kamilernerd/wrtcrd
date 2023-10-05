# wrtcrd
Remote desktop server and client using Webrtc protocol.

## WIP
This project is a work in progress and mainly used to have some fun programming
and to expand horizons into different fields and topics but the end solution
should have most of the todo list implemented.

## Features
- Cross platform server (windows, linux, darwin) and client (web)
- Multiple display support at full resolution
- Audio output (Opus) 48KHz Mono
- Mouse input
- Keyboard input
- Webrtc 60fps (x264)

## TODO
- Adjustable framerate
- Adjustable frame down-scaling (Who needs to send FULL HD / 4k over wire?)
- Actually implement the mouse and keyboard support
- Look for performance/stability improvements
- Add support for audio codecs AAC, Vorbis
- Add support for video codecs VP8/9 and AV1
- Logs
- Authentication
- General server improvements
- SFU for party usage (Does this really make any sense?)
- Very important thing is to go away from using libraries that do most of the magic behind and create own implementations.
- TLS!!!!

## How to set it up
This project is not really tested on other platforms than darwin using arm64 so
a lot of things are a subject to change and might still not work as expected on
other platforms so keep that in mind...

Well this is a bit weird... You need to install x264 and Opus codecs somehow
(FFMPEG?).
Depending on your platform things will be a bit different.

### Darwin
It turns out that there is no simple way around CoreAudio to create a virtual
audio output device so we could send audio out. Therefore you will have to
install BlackHole and create multi-output device then set it up to output to your speakers/other device and
blackhole. Then set input device to blackhole as well (microphone input). Now sound should work.

Install libportaudio using brew

That should be everything you need to run this.

Now install the dependencies and run this

```GOOS=darwin GOARCH=<cpu architecture> CGO_ENABLED=1 go run main.go```

There is a bigger list of dependencies you should have but they are in the readme
files of the dependencies.

### Linux
Not tested so there's no installation steps

### Windows
Not tested so there's no installation steps
