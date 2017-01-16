package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {

	//USE THIS TO STREAM AN AUDIO FILE LOOPED FROM DISK
	//------------------
	source := []string{"-stream_loop", "-1", "-i", "source/jorkkasten.mp3"}

	//USE THIS TO STREAM YOUR AUDIO INPUT ON LINUX
	//------------------
	//source := []string{"-f", "alsa", "-i", "hw:0"}

	//USE THIS TO STREAM YOUR AUDIO INPUT ON WINDOWS
	//Run this to list your devices "ffmpeg -list_devices true -f dshow -i dummy"
	//Replace the soundcard with your selection
	//------------------
	//source := []string{"-f", "dshow", "-i", "audio=\"Microphone (Cirrus Logic CS4206A (AB 71))\""}

	//USE THIS TO STREAM YOUR AUDIO INPUT ON MAC
	//Run ffmpeg -f avfoundation -list_devices true -i ""
	//Figure out which number is your audio input
	//------------------
	//source := []string{"-f", "avfoundation", "-i", ":0"}

	var wg sync.WaitGroup

	wg.Add(1)
	go RunWeb(&wg)

	wg.Add(1)
	go RunFFmpeg(&wg, source)

	wg.Wait()
}

func RunWeb(wg *sync.WaitGroup) (err error) {
	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("static", true)))
	r.Use(static.Serve("/media", static.LocalFile("media", true)))
	err = r.Run()

	wg.Done()

	return err

}

func RunFFmpeg(wg *sync.WaitGroup, source []string) (err error) {

	s := uniuri.New()

	ffmpegArgs := append([]string{
		"-y",
		"-re",
		"-protocol_whitelist", "file,http,https,tcp,tls",
	}, source...)

	ffmpegArgs = append(ffmpegArgs, []string{
		"-acodec", "aac",
		"-b:a", "128k",
		"-ar", "44100",
		"-strict", "-2",
		"-flags", "+global_header",
		"-f", "hls",
		"-hls_list_size", "20",
		"-hls_time", "6",
		"-hls_segment_filename", fmt.Sprintf("media/hls/%s-fl%s.ts", s, "%03d"),
		"-hls_flags", "delete_segments",
		"media/hls/out.m3u8",
	}...)

	cmd := exec.Command("ffmpeg", ffmpegArgs...)

	stderr, _ := cmd.StderrPipe()
	in := bufio.NewScanner(stderr)

	if err = cmd.Start(); err != nil {
		return err
	}

	for in.Scan() {
		fmt.Println(in.Text())
	}

	err = cmd.Wait()
	wg.Done()
	return err
}
