package mailer

type Client interface {
	To(name, address string) Client
	From(name, address string) Client
	ReplyTo(name, address string) Client
	AddCC(name, address string) Client
	AddBCC(name, address string) Client
	AddFile(name, path string) Client
	Subject(subject string) Client
	Html(value string) Client
	Template(name string, payload any) Client
	Text(value string) Client
	Send() error
}

type Address struct {
	Name    string
	Address string
}

type File struct {
	Name string
	Path string
}

type Template struct {
	Payload any
	Name    string
}
