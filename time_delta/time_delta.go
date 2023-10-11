package timedelta

import "sync"

// we have compressed the data spatially, but now we can compress temporally, by finding pixels that don't change a lot
// betwee frames
func TimeDelta(frames [][]byte) [][]byte {
	// Delta is mostly zeros!
	// compress repeated values with run length encoding.
	// store the number of times a value repeats, then the value
	
	encoded := make([][]byte, len(frames))
	wg := sync.WaitGroup{}

	for i := range frames {
		wg.Add(1)
		go processFrame(frames, i, &wg, &encoded)
	}
	wg.Wait()

	return encoded
}

func processFrame(frames [][]byte, frame_ind int, wg *sync.WaitGroup, encoded *[][]byte) {
	defer wg.Done()
	curr_frame := frames[frame_ind]

	if frame_ind == 0 {
		(*encoded)[frame_ind] = curr_frame
		return
	}
	prev_frame := frames[frame_ind-1]

	delta := make([]byte, len(curr_frame))
	for j := 0; j < len(delta); j++ { // subtract frames to find the deltas
		delta[j] = curr_frame[j] - prev_frame[j]
	}

	var rle []byte
	for j := 0; j < len(delta); j++ {
		// Count the number of times the current value repeats.
		var count byte
		for count = 0; count < 255 && j+int(count) < len(delta) && delta[j+int(count)] == delta[j]; count++ {
		}

		// Store the count and value.
		rle = append(rle, count)
		rle = append(rle, delta[j])

		j += int(count)
	}

	// Save the RLE frame.
	(*encoded)[frame_ind] = rle
}
