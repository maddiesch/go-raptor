package test

func Ptr[T any](x T) *T {
	return &x
}
