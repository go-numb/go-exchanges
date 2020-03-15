package v1

type Requester interface {
	Path() string
	Method() string
	Query() string
	Payload() []byte
}
