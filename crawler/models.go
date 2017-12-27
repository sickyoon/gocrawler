package crawler

import "io"

// Request -
type Request interface {
	URL() string
	Method() string
	Process(chan Request, io.Reader) error
}
