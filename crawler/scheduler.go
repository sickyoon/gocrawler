package crawler

// Scheduler -
type Scheduler struct {
	reqCh      chan Request
	numWorkers int
}

// NewScheduler -
func NewScheduler(numWorkers int, bufSize int) *Scheduler {
	// share memory by communicating
	reqCh := make(chan Request, bufSize)
	return &Scheduler{
		reqCh: reqCh,
	}
}

// StartWorkers -
func (s *Scheduler) StartWorkers() {
	for i := 0; i < s.numWorkers; i++ {
		go NewWorker(i, s.reqCh)
	}
}

// Schedule -
func (s *Scheduler) Schedule(req Request) {
	s.reqCh <- req
}
