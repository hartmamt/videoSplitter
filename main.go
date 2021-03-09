package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"os"
	"strconv"
	"strings"
)

type videoMetaData struct {
	Width float64
	Height float64
	FrameRate float64
	Duration float64
	DurationTS int
	FrameCount int
}

func main(){

	// Get Metadata
	a, err := ffmpeg.Probe(
		"/Users/hartmamt/Projects/pg/videoSplitter/in/durst.mp4",
		ffmpeg.KwArgs{},
		)
	if err != nil {
		panic(err)
	}


	/*
	{
	            "BoundingBox": {
	                "Width": 0.06047183275222778,
	                "Height": 0.16532716155052185,
	                "Left": 0.4213453233242035,
	                "Top": 0.39640137553215027
	            },


	"BoundingBox": {
		"Width": 0.12551015615463257,
			"Height": 0.2838538885116577,
			"Left": 0.6826857328414917,
			"Top": 0.2881534993648529
	},
*/

	//b, err := ioutil.ReadFile("/Users/hartmamt/Projects/pg/videoSplitter/in/durst.json") // just pass the file name
	//if err != nil {
	//	fmt.Print(err)
	//}
//#(nets.#(=="fb"))#
//(Face.Confidence.#(>50.0))#.#
//{name.first,age,"the_murphys":friends.#(last="Murphy")#.first}
//	result := gjson.Get(string(b), "{Persons.#.Timestamp,\"faces\":Persons.#.Person.Face}") // print the content as 'bytes'
//
//	result.ForEach(func(key, value gjson.Result) bool {
//		//fmt.Println(key)
//		println(value.String())
//		return false
//		return true // keep iterating
//	})

//	return
	//str := string(b) // convert content to a 'string'
	/*
	ffprobe video.mp4 -select_streams v -show_entries frame=coded_picture_number,pkt_pts_time -of csv=p=0:nk=1 -v 0
	 */

	//fmt.Println(a)
	frameRate := gjson.Get(a, "streams.0.r_frame_rate").String()
	calculatedFrameRateA, _ := strconv.ParseFloat(strings.Split(frameRate,"/")[0],64)
	calculatedFrameRateB, _ := strconv.ParseFloat(strings.Split(frameRate,"/")[1], 64)

	metaData := videoMetaData{
		Width: gjson.Get(a, "streams.0.width").Float(),
		Height: gjson.Get(a, "streams.0.height").Float(),
		FrameRate: calculatedFrameRateA/calculatedFrameRateB,
		Duration: gjson.Get(a, "streams.0.duration").Float(),
		DurationTS: int(gjson.Get(a, "streams.0.duration_ts").Int()),
		FrameCount: int(gjson.Get(a, "streams.0.nb_frames").Int()),
	}

	fmt.Printf("%v", metaData)

	//
	//
	//err = ffmpeg.Input("/Users/hartmamt/Projects/pg/videoSplitter/in/bezos_vogels.mp4").
	//	Filter("aresample",ffmpeg.Args{"async=1000"}).
	//	Output("/Users/hartmamt/Projects/pg/videoSplitter/test.mp3").
	//	OverWriteOutput().
	//	Run()
	// Break Into Frames

	err = ffmpeg.
		Input("/Users/hartmamt/Projects/pg/videoSplitter/in/durst.mp4").
		//Filter("fps", ffmpeg.Args{fmt.Sprintf("%f", metaData.FrameRate)}).
		Filter("select", ffmpeg.Args{"'between(n\\,1\\,100)'"}).
		Output("/Users/hartmamt/Projects/pg/videoSplitter/out/test-%04d.jpg", ffmpeg.KwArgs{
			"start_number": 0,
			//"select": "between(n\\,1\\,100)",
		}).
		OverWriteOutput().
		Run()

	fmt.Println(err)

	return
	/*
	get timestamps for frames
	 ffprobe -f lavfi -i "movie=durst.mp4,fps=fps=29.97[out0]" -show_frames -show_entries frame=pkt_pts_time -of csv=p=0 > frames.txt
	 */
	//
	//if err!=nil {
	//	fmt.Println(err)
	//}

	videoPath := "/Users/hartmamt/Projects/pg/videoSplitter/in/durst.mp4"
	videoOutput := "/Users/hartmamt/Projects/pg/videoSplitter/out/"
	combinedVideoPath := videoOutput +  "combined.mp4"
//https://itectec.com/superuser/how-to-extract-the-timestamps-associated-with-frames-ffmpeg-extracts-from-a-video-with-the-r-option/

	/*
	image 0176
	  {
	            "Timestamp": 5038,
	            "Person": {
	                "Index": 1,
	                "BoundingBox": {
	                    "Width": 0.34453123807907104,
	                    "Height": 0.737500011920929,
	                    "Left": 0.08437500149011612,
	                    "Top": 0.17916665971279144
	                }
	            }
	        },
	0.035
	5028 * ? = 176
	 */

	//// Reassemble Video
	video := ffmpeg.Input(
			videoOutput + "*.jpg",
			ffmpeg.KwArgs{
				"pattern_type":"glob",
				"framerate": metaData.FrameRate,
			},
		)
	fmt.Println(video)
	audio := ffmpeg.Input(videoPath)

	err = ffmpeg.Concat(
			[]*ffmpeg.Stream{video, audio,},
			ffmpeg.KwArgs{
				"v":1,
				"a":1,
			},
		).
		Output(combinedVideoPath).
		OverWriteOutput().
		Run()

	if err!=nil {
		fmt.Fprint(os.Stderr, "request cancelled\n")
		fmt.Println(err)
	}



}