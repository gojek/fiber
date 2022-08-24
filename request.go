package fiber

type Request interface {
	Payload() interface{}
	Header() map[string][]string
	Clone() (Request, error)
	OperationName() string
	Protocol() string

	Transform(backend Backend) (Request, error)
}

const (
	HTTP string = "HTTP"
	GRPC string = "GRPC"
)
