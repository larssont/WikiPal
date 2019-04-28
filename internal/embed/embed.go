package embed

//Message struct for embedded messages
type Message struct {
	Color       int
	Description string
	Title       string
	Fields      map[string]string
}
