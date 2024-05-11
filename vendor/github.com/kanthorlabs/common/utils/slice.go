package utils

func ChunkNext[T int | int32 | int64 | uint | uint32 | uint64](prev, end, step T) T {
	if prev+step > end {
		return end
	}
	return prev + step
}
