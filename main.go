package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/tidwall/gjson"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type videoMetaData struct {
	Width float64
	Height float64
	FrameRate string
	Duration float64
	DurationTS int
	FrameCount int
}

func main(){

	// Get Metadata
	a, err := ffmpeg.Probe("/Users/hartmamt/Projects/pg/videoSplitter/in/bezos_vogels.mp4")
	if err != nil {
		panic(err)
	}

	fmt.Println(a)
	metaData := videoMetaData{
		Width: gjson.Get(a, "streams.0.width").Float(),
		Height: gjson.Get(a, "streams.0.height").Float(),
		FrameRate: gjson.Get(a, "streams.0.height").String(),
		Duration: gjson.Get(a, "streams.0.duration").Float(),
		DurationTS: int(gjson.Get(a, "streams.0.duration_ts").Int()),
		FrameCount: int(gjson.Get(a, "streams.0.nb_frames").Int()),
	}

	fmt.Printf("%v", metaData)

	/*
	          "BoundingBox": {
	            "Width": 0.02242584154009819,
	            "Height": 0.05018870159983635,
	            "Left": 0.23596711456775665,
	            "Top": 0.43446972966194153
	          },
	 */

	/*
		ffmpeg -i derpdog.mp4 -filter_complex \
		 "[0:v]crop=200:200:60:30,boxblur=10[fg]; \
		  [0:v][fg]overlay=60:30[v]" \
		-map "[v]" -map 0:a -c:v libx264 -c:a copy -movflags +faststart derpdogblur.mp4
	*/

	box := rekognition.BoundingBox{
		Height: aws.Float64(0.4016371965408325),
		Left: aws.Float64(0.42521291971206665),
		Top: aws.Float64(0.17050546407699585),
		Width: aws.Float64(0.1513211727142334),
	}



	// Break Into Frames
	err = ffmpeg.
		Input("/Users/hartmamt/Projects/pg/videoSplitter/in/bezos_vogels.mp4").
		Filter("fps", ffmpeg.Args{"30/1"}).
		Output("/Users/hartmamt/Projects/pg/videoSplitter/out/test-00%d.jpg", ffmpeg.KwArgs{"start_number": 0}).
		OverWriteOutput().
		Run()

	if err!=nil {
		fmt.Println(err)
	}

	overlayImage := ffmpeg.
		Input("/Users/hartmamt/Projects/pg/videoSplitter/out/test-001.jpg").
		Crop(
			int(*box.Left * metaData.Width),
			int(*box.Top * metaData.Height),
			int(*box.Width * metaData.Width),
			int(*box.Height * metaData.Height)).
		Filter("boxblur", ffmpeg.Args{"10"})

	err = ffmpeg.
		Input("/Users/hartmamt/Projects/pg/videoSplitter/out/test-001.jpg").
		//DrawBox(int(*box.Left * metaData.Width), int(*box.Top * metaData.Height), int(*box.Width * metaData.Width), int(*box.Height * metaData.Height), "red", 2).
		//Filter("fps", ffmpeg.Args{"30/1"}).
		Overlay(overlayImage,"", ffmpeg.KwArgs{
			"x": int(*box.Left * metaData.Width),
			"y": int(*box.Top * metaData.Height),
		}).
		Output("/Users/hartmamt/Projects/pg/videoSplitter/new-001.jpg").
		OverWriteOutput().
		Run()

	if err != nil {

		fmt.Println(err.Error())
	}

	//func ComplexFilterExample(testInputFile, testOverlayFile, testOutputFile string) *ffmpeg.Stream {
	//	split := ffmpeg.Input(testInputFile).VFlip().Split()
	//	split0, split1 := split.Get("0"), split.Get("1")
	//	overlayFile := ffmpeg.Input(testOverlayFile).Crop(10, 10, 158, 112)
	//	return ffmpeg.Concat([]*ffmpeg.Stream{
	//	split0.Trim(ffmpeg.KwArgs{"start_frame": 10, "end_frame": 20}),
	//	split1.Trim(ffmpeg.KwArgs{"start_frame": 30, "end_frame": 40})}).
	//	Overlay(overlayFile.HFlip(), "").
	//	DrawBox(50, 50, 120, 120, "red", 5).
	//	Output(testOutputFile).
	//	OverWriteOutput()
	//}

	//err = ffmpeg.
	//	Input("/Users/hartmamt/Projects/pg/videoSplitter/out/test-001.jpg").
	//	//DrawBox(int(*box.Left * metaData.Width), int(*box.Top * metaData.Height), int(*box.Width * metaData.Width), int(*box.Height * metaData.Height), "red", 2).
	//	//Filter("fps", ffmpeg.Args{"30/1"}).
	//	Overlay(overlayFile,"").
	//	Output("/Users/hartmamt/Projects/pg/videoSplitter/out/new-001.jpg").
	//	OverWriteOutput().
	//	Run()

	if err!=nil {
		fmt.Println("after overlay")
		fmt.Println(err)
	}

	//split := Input(TestInputFile1).VFlip().Split()
	//split0, split1 := split.Get("0"), split.Get("1")
	//overlayFile := Input(TestOverlayFile).Crop(10, 10, 158, 112)
	//err := Concat([]*Stream{
	//	split0.Trim(KwArgs{"start_frame": 10, "end_frame": 20}),
	//	split1.Trim(KwArgs{"start_frame": 30, "end_frame": 40})}).
	//	Overlay(overlayFile.HFlip(), "").
	//	DrawBox(50, 50, 120, 120, "red", 5).
	//	Output(TestOutputFile1).
	//	OverWriteOutput().
	//	Run()



	/*
	ffmpeg -i derpdog.mp4 -filter_complex \
	 "[0:v]crop=200:200:60:30,boxblur=10[fg]; \
	  [0:v][fg]overlay=60:30[v]" \
	-map "[v]" -map 0:a -c:v libx264 -c:a copy -movflags +faststart derpdogblur.mp4
	 */


	//
	//// Reassemble Video
	err = ffmpeg.
		Input(
			"/Users/hartmamt/Projects/pg/videoSplitter/out/test*.jpg",
			ffmpeg.KwArgs{
				"pattern_type":"glob",
				"framerate": 30000/1001,
			},
		).
		Output("/Users/hartmamt/Projects/pg/videoSplitter/out/move.mp4").
		OverWriteOutput().
		Run()
	//
	//if err!=nil {
	//	fmt.Println(err)
	//}

}