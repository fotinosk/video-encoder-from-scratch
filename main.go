package main

import (
	"codec/deflate"
	yuv "codec/rgb_to_yuv"
	timedelta "codec/time_delta"
	"codec/utils"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	start := time.Now()
	var width, height int
	flag.IntVar(&width, "width", 384, "width of the video")
	flag.IntVar(&height, "height", 216, "height of the video")
	flag.Parse()

	frames := make([][]byte, 0)

	for {
		frame := make([]byte, width*height*3)
		if _, err := io.ReadFull(os.Stdin, frame); err != nil {
			break
		}
		frames = append(frames, frame)
	}

	fmt.Printf("Initial file: %d, %d\n", len(frames), utils.FastSize(frames))

	yuv_frames := yuv.ConvertRGBtoYUV(frames, width, height)
	fmt.Printf("YUV file: %d\n", utils.FastSize(yuv_frames))
	fmt.Printf("Compression ratio: %0.2f%% \n", float32(utils.FastSize(yuv_frames))/float32(utils.FastSize(frames)))

	time_compressed := timedelta.TimeDelta(yuv_frames)

	fmt.Printf("TD file: %d\n", utils.Size(time_compressed))
	fmt.Printf("Compression ratio: %0.2f%% \n", float32(utils.Size(time_compressed))/float32(utils.FastSize(frames)))

	// TODO: fix deflate
	deflated := deflate.Deflate(time_compressed)
	fmt.Printf("Deflated file: %d\n", deflated.Len())

	elapsed := time.Since(start)
	fmt.Printf("Processing took %s", elapsed)
}
