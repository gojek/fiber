package fiber

type Request interface {
	Payload() interface{}
	Header() map[string][]string
	Clone() (Request, error)
	OperationName() string

	Transform(backend Backend) (Request, error)
}
