package suggest

type Suggest struct {
	SuggesterName string
	Text          string
	Field         string
	PayloadKeys   map[string]string
}
