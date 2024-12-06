package types

func Nothing[T any]() T {
	var zero T
	return zero
}