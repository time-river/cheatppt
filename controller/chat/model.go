package chat

type Model struct {
	Id          int    `json:"id"`
	ModelName   string `json:"modelName"`
	DisplayName string `json:"displayName"`
	IsChatGPT   bool   `json:"isChatGPT"`
}

// @Default: `models[]` index instead of `Model.Id`
type ModelSetting struct {
	Default int     `json:"default"`
	Models  []Model `json:"models"`
}

func GetModelSetting(level int) ModelSetting {
	return ModelSetting{
		Default: 0,
		Models: []Model{
			{Id: 0, ModelName: "gpt-3.5-turbo", DisplayName: "gpt-3.5-turbo", IsChatGPT: false},
			{Id: 1, ModelName: "gpt-4", DisplayName: "gpt-4", IsChatGPT: false},
			{Id: 10000, ModelName: "text-davinci-002-render-sha", DisplayName: "ChatGPT-3 (unstable)", IsChatGPT: true},
			{Id: 10001, ModelName: "gpt-4", DisplayName: "ChatGPT-4 (unstable)", IsChatGPT: true},
		},
	}
}
