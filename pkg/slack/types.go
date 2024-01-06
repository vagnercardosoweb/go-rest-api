package slack

import "sync"

type ColorName string

const (
	ColorError   ColorName = "error"
	ColorWarning ColorName = "warning"
	ColorSuccess ColorName = "success"
	ColorInfo    ColorName = "info"
)

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type Client struct {
	fields      []Field
	channel     string
	username    string
	memberIds   []string
	environment string
	color       ColorName
	token       string
	mu          *sync.Mutex
}
