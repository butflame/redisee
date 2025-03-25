package reporter

import "fmt"

type Reporter interface {
	Report()
}

type StdoutReporter struct {
}

func (s *StdoutReporter) Report() {
	// todo
}

func Report() {
	fmt.Println("hello world")
}
