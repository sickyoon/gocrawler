package crawler

// Scheduler -
type Scheduler struct {
	reqCh      chan Request
	numWorkers int
}

// New -
func New(numWorkers int, bufSize int) *Scheduler {
	// share memory by communicating
	reqCh := make(chan Request, bufSize)
	return &Scheduler{
		reqCh:      reqCh,
		numWorkers: numWorkers,
	}
}

// StartWorkers -
func (s *Scheduler) StartWorkers() error {
	for i := 0; i < s.numWorkers; i++ {
		w, err := NewWorker(i, s.reqCh)
		if err != nil {
			return err
		}
		go w.Start()
	}
	return nil
}

// Schedule -
func (s *Scheduler) Schedule(req Request) {
	s.reqCh <- req
}
