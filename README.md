# simple-ffmpeg-hls-audio-project
Simple demo showing of how to use FFmpeg to stream audio using HLS. Comes with a small HTTP server.

###HOW TO RUN?

- Install FFmpeg 3.x
- go get github.com/coral/simple-ffmpeg-hls-audio-project
- go run run.go
- http://localhost:8080

###HOW TO CHANGE AUDIO SOURCE?
In run.go there are some examples on how to capture your audio input instead of looping an audio file.