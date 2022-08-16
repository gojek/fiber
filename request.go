package fiber

type Request interface {
	Payload() interface{}
	Header() map[string][]string
	Clone() (Request, error)
	OperationName() string
	Protocol() Protocol

	Transform(backend Backend) (Request, error)
}

type Protocol string

func (p Protocol) String() string {
	return string(p)
}

const (
	HTTP Protocol = "HTTP"
	GRPC Protocol = "GRPC"
)
