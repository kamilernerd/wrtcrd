# wrtcrd
Remote desktop server and client using Webrtc protocol.

## WIP
This project is a work in progress and mainly used to have some fun programming
and to expand horizons into different fields and topics but the end solution
should have most of the todo list implemented.

## Features
- Cross platform server (windows, linux, darwin) and client (web) (tested only
  on darwin)
- Multiple display support at full resolution
- Audio output (Opus) 48KHz Mono
- Mouse input (Left, right, middle buttons, scroll)
- Keyboard input
- Webrtc 60fps (x264)

## TODO
- Adjustable framerate
- Adjustable frame down-scaling (Who needs to send FULL HD / 4k over wire?)
- Look for performance/stability improvements
- Add support for audio codecs AAC, Vorbis
- Add support for video codecs VP8/9 and AV1
- Logs
- Authentication
- Refactoring
- TLS

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
```$ brew install portaudio```

Then install the opus codec
```$ brew install pkg-config opus opusfile```

That should be everything you need to run this.

Now install the dependencies and run this

```$ GOOS=darwin GOARCH=<cpu architecture> CGO_ENABLED=1 go run main.go```

### Linux
Not tested so there's no installation steps

### Windows
Not tested so there's no installation steps
