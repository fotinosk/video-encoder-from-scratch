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
	fmt.Printf("Compression ratio: %0.2f%% \n", 100*float32(utils.FastSize(yuv_frames))/float32(utils.FastSize(frames)))

	start := time.Now()
	time_compressed := timedelta.TimeDelta(yuv_frames)
	time_rle := time.Since(start)

	fmt.Printf("TD file: %d\n", utils.Size(time_compressed))
	fmt.Printf("Compression ratio: %0.2f%% \n", 100*float32(utils.Size(time_compressed))/float32(utils.FastSize(frames)))
	fmt.Printf("LRE Processing took %s\n", time_rle)

	start_delfated := time.Now()
	deflated := deflate.Deflate(yuv_frames)
	time_deflated := time.Since(start_delfated)
	fmt.Printf("Deflate Processing took %s\n", time_deflated)
	fmt.Printf("Deflated ratio: %0.2f%%\n", 100*float32(deflated.Len()) / float32(utils.FastSize(frames)))
}
