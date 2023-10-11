package yuv

import (
	"sync"
)

// convert rgb values to YUV420 values
// Red, Green, Blue -> Luma (Y), Chroma (UV)

func ConvertRGBtoYUV(frames [][]byte, width int, height int) [][]byte {
	//save the brightness for each pixel but the chroma for every 2x2 square
	// go from R,G,B,R,G,B,R,G,B,R,G,B... to Y,Y,Y,Y,Y,U,U,U,U,U,U,V,V,V,V,V,V, for each frame

	wg := sync.WaitGroup{}
	yuvFrames2 := make([][]byte, len(frames))

	for i, frame := range frames {
		wg.Add(1)
		go processFrame(frame, i, width, height, &wg, &yuvFrames2)
	}
	wg.Wait()

	return yuvFrames2
}

func processFrame(frame []byte, ind int, width int, height int, wg *sync.WaitGroup, yuvFrames2 *[][]byte) {
	defer wg.Done()
	Y := make([]byte, width*height)
	U := make([]float64, width*height)
	V := make([]float64, width*height)
	for j := 0; j < width*height; j++ {
		r, g, b := float64(frame[3*j]), float64(frame[3*j+1]), float64(frame[3*j+2])

		y := +0.299*r + 0.587*g + 0.114*b
		u := -0.169*r - 0.331*g + 0.449*b + 128
		v := 0.499*r - 0.418*g - 0.0813*b + 128

		Y[j] = uint8(y)
		U[j] = u
		V[j] = v
	}

	uDownsampled := make([]byte, width*height/4)
	vDownsampled := make([]byte, width*height/4)
	for x := 0; x < height; x += 2 {
		for y := 0; y < width; y += 2 {
			u := (U[x*width+y] + U[x*width+y+1] + U[(x+1)*width+y] + U[(x+1)*width+y+1]) / 4
			v := (V[x*width+y] + V[x*width+y+1] + V[(x+1)*width+y] + V[(x+1)*width+y+1]) / 4

			uDownsampled[x/2*width/2+y/2] = uint8(u)
			vDownsampled[x/2*width/2+y/2] = uint8(v)
		}
	}

	yuvFrame := make([]byte, len(Y)+len(uDownsampled)+len(vDownsampled))

	copy(yuvFrame, Y)
	copy(yuvFrame[len(Y):], uDownsampled)
	copy(yuvFrame[len(Y)+len(uDownsampled):], vDownsampled)

	(*yuvFrames2)[ind] = yuvFrame
}
