package fiber

import "github.com/gojek/fiber/protocol"

type Request interface {
	Payload() interface{}
	Header() map[string][]string
	Clone() (Request, error)
	OperationName() string
	Protocol() protocol.Protocol

	Transform(backend Backend) (Request, error)
}
