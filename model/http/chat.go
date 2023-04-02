package msg

type Model struct {
	Owner    string `json:"owner"`
	Id       string `json:"Id"`
	MaxToken int    `json:"maxToken"`
}

type ListModelsResponse struct {
	Message string  `json:"message"`
	Models  []Model `json:"data"`
}
