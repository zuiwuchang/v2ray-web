package utils

type Pair[First any, Second any] struct {
	First  First
	Second Second
}

func MakePair[First any, Second any](first First, second Second) Pair[First, Second] {
	return Pair[First, Second]{
		First:  first,
		Second: second,
	}
}
