package mailer

type Mailer interface {
	To(name, address string) Mailer
	From(name, address string) Mailer
	ReplyTo(name, address string) Mailer
	AddCC(name, address string) Mailer
	AddBCC(name, address string) Mailer
	AddFile(name, path string) Mailer
	Subject(subject string) Mailer
	Html(value string) Mailer
	Template(name string, payload any) Mailer
	Text(value string) Mailer
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
