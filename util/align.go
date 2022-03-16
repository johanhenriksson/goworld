package util

func Align(offset, alignment int) int {
	count := offset / alignment
	diff := offset % alignment
	if diff > 0 {
		count++
	}
	return count * alignment
}
