package crawler

// Request -
type Request interface {
	URL() string
	Method() string
	Process(chan Request, []byte)
}

// TestRequest -
type TestRequest struct {
	url    string
	method string
}

// URL -
func (r TestRequest) URL() string {
	return r.url
}

// Method -
func (r TestRequest) Method() string {
	return r.method
}

// Process data
func (r TestRequest) Process(reqCh chan Request, body []byte) {

}
