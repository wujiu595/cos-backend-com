package locales

type Content struct {
	Message       string   `json:"message"`
	MessageId     string   `json:"messageId"`
	MessageParams []string `json:"mesageParams"`
}

var (
	Message = func(message string, others ...string) (res Content) {
		res.Message = message
		if len(others) > 0 {
			res.MessageId = others[0]
		}
		if len(others) > 1 {
			res.MessageParams = others[1:]
		}
		return
	}

	Required = func(v1 string) (res Content) {
		return
	}
)
