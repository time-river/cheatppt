package model

type AclDetail struct {
	Username  uint
	ModelName string
	Provider  string
}

func Allow(model string) bool {
	return true
}
