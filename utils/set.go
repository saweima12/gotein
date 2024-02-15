package utils

func SliceToSet[T comparable](d []T) map[T]struct{} {
	resp := make(map[T]struct{})

	for i := range d {
		resp[d[i]] = struct{}{}
	}
	return resp
}
