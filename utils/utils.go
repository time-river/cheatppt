package utils

func Must[T interface{}](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}
