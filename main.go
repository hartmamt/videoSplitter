package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type videoMetaData struct {
	Width int
	Height int
	FrameRate string
	Duration float64
	DurationTS int
	FrameCount int
}

func main(){

	// Get Metadata
	a, err := ffmpeg.Probe("/Users/hartmamt/Projects/pg/videoSplitter/in/face3.mp4")
	if err != nil {
		panic(err)
	}

	fmt.Println(a)
	metaData := videoMetaData{
		Width: int(gjson.Get(a, "streams.0.width").Int()),
		Height: int(gjson.Get(a, "streams.0.height").Int()),
		FrameRate: gjson.Get(a, "streams.0.height").String(),
		Duration: gjson.Get(a, "streams.0.duration").Float(),
		DurationTS: int(gjson.Get(a, "streams.0.duration_ts").Int()),
		FrameCount: int(gjson.Get(a, "streams.0.nb_frames").Int()),
	}

	fmt.Printf("%v", metaData)

	// Break Into Frames
	err = ffmpeg.
		Input("/Users/hartmamt/Projects/pg/videoSplitter/in/face3.mp4").
		Filter("fps", ffmpeg.Args{"30/1"}).
		Output("/Users/hartmamt/Projects/pg/videoSplitter/out/test-00%d.jpg", ffmpeg.KwArgs{"start_number": 0}).
		OverWriteOutput().
		Run()

	if err!=nil {
		fmt.Println(err)
	}

	// Reassemble Video
	err = ffmpeg.
		Input(
			"/Users/hartmamt/Projects/pg/videoSplitter/out/*.jpg",
			ffmpeg.KwArgs{
				"pattern_type":"glob",
				"framerate": 30,
			},
		).
		Output("/Users/hartmamt/Projects/pg/videoSplitter/out/move.mp4").
		Run()

	if err!=nil {
		fmt.Println(err)
	}

}