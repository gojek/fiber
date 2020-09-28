package fiber

//
type Dispatcher interface {
	Do(request Request) Response
}
