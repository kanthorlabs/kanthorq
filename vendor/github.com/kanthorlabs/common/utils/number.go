package utils

func Min[T int | int32 | int64 | uint | uint32 | uint64 | float32 | float64](x, y T) T {
	if x > y {
		return y
	}
	return x
}

func Max[T int | int32 | int64 | uint | uint32 | uint64 | float32 | float64](x, y T) T {
	if x < y {
		return y
	}
	return x
}
