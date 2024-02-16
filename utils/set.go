package utils

func SliceToSet[T comparable](d []T) map[T]struct{} {
	resp := make(map[T]struct{})

	for i := range d {
		resp[d[i]] = struct{}{}
	}
	return resp
}

func SetToSlice[T comparable](d map[T]struct{}) []T {
	resp := make([]T, 0, len(d))

	for k := range d {
		resp = append(resp, k)
	}
	return resp
}
