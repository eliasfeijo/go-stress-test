package stress

import "fmt"

// StressReport struct
type StressReport struct {
	// The number of requests made
	Requests int
	// The number of requests that failed
	Failed int
	// The number of requests that succeeded
	Succeeded int
	// The number of requests that timed out
	TimedOut int
	// The total time taken to complete all requests
	TotalTime float64
	// The average time taken to complete all requests
	AverageTime float64
	// The fastest time taken to complete a request
	FastestTime float64
	// The slowest time taken to complete a request
	SlowestTime float64
	// The percentage of requests that succeeded
	PercentageSucceeded float64
	// The percentage of requests that failed
	PercentageFailed float64
	// The percentage of requests that timed out
	PercentageTimedOut float64
}

// NewStressReport creates a new stress report
func NewStressReport() *StressReport {
	return &StressReport{
		Requests:            0,
		Failed:              0,
		Succeeded:           0,
		TimedOut:            0,
		TotalTime:           0,
		AverageTime:         0,
		FastestTime:         0,
		SlowestTime:         0,
		PercentageSucceeded: 0,
		PercentageFailed:    0,
		PercentageTimedOut:  0,
	}
}

// Stress interface
type IStress interface {
	// Run the stress test
	Run() error
	// Get the report
	PrintReport()
}

// Stress test
type Stress struct {
	// The URL to test
	URL string
	// The HTTP method to use
	Method string
	// The number of concurrent requests to make
	Concurrency int
	// The number of requests to make
	Requests int
	// The timeout in seconds
	Timeout int
	// Verbose output
	Verbose bool
	// The report
	Report *StressReport
}

// NewStress creates a new stress test
func NewStress(url string, method string, concurrency int, requests int, timeout int, verbose bool) *Stress {
	fmt.Println("Creating stress test..." + url)
	report := NewStressReport()
	return &Stress{
		URL:         url,
		Method:      method,
		Concurrency: concurrency,
		Requests:    requests,
		Timeout:     timeout,
		Verbose:     verbose,
		Report:      report,
	}
}

// Run the stress test
func (s *Stress) Run() error {
	fmt.Println("Running stress test...")
	return nil
}

// PrintReport prints the report
func (s *Stress) PrintReport() {
	fmt.Println("Report:")
	fmt.Println("Requests:", s.Report.Requests)
	fmt.Println("Failed:", s.Report.Failed)
	fmt.Println("Succeeded:", s.Report.Succeeded)
	fmt.Println("TimedOut:", s.Report.TimedOut)
}
