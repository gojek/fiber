package fiber

type Request interface {
	Payload() []byte
	Header() map[string][]string
	Clone() (Request, error)
	OperationName() string

	Transform(backend Backend) (Request, error)
}
