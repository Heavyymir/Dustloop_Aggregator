package display

import(
	"fmt"
	"sync"
	"time"
)

// Spinner manages a terminal progress animation
type Spinner struct {
	message string
	stop	chan struct{}
	wg		sync.WaitGroup
}

// Startspinner starts a spinner animation on the current terminal line
func StartSpinner(message string) *Spinner {
	s := &Spinner {
		message:	message,
		stop:		make(chan struct{}),
	}

	frames := []string{"|", "/", "-", "\\"}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		i := 0
		for {
			select {
				case<-s.stop:
					// Clear the spinner line completely when it stops
					// \r moves to line start, padding with spaces removes lingering text
					fmt.Printf("\r%-80s\r", "")
					return
				default:
					frame := frames[i%len(frames)]
					// \r brings cursor to start without a newline
					fmt.Printf("\r[%s] %s", frame, s.message)
					i++
					time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	return s
}

// Stop cleanly terminates the spinner and clears the line
func (s *Spinner) Stop() {
	close(s.stop)
	s.wg.Wait()
}
