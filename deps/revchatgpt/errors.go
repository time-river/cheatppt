package revchatgpt

func (err *ChatGPTError) Error() string {
	return err.StatusText
}
