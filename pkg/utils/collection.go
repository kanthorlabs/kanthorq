package utils

func MergeMaps[T any](dest map[string]T, maps ...map[string]T) map[string]T {
	for _, m := range maps {
		for k, v := range m {
			dest[k] = v
		}
	}
	return dest
}
