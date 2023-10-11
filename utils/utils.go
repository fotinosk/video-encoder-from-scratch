package utils

func Size(frames [][]byte) int {
	var size int
	for _, frame := range frames {
		size += len(frame)
	}
	return size
}

func FastSize(frames [][]byte) int {
	return len(frames[0]) * len(frames)
}