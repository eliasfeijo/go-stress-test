package stress

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type MapStatusRequests map[int]int

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
	FastestTime int64
	// The slowest time taken to complete a request
	SlowestTime int64
	// The percentage of requests that succeeded
	PercentageSucceeded float64
	// The percentage of requests that failed
	PercentageFailed float64
	// The percentage of requests that timed out
	PercentageTimedOut float64
	// The status code of the requests
	StatusRequests MapStatusRequests
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
		StatusRequests:      make(MapStatusRequests),
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
	// Verify TLS
	VerifyTls bool
	// Mutex
	mu sync.Mutex
}

// NewStress creates a new stress test
func NewStress(url string, method string, concurrency int, requests int, timeout int, verifyTls bool, verbose bool) *Stress {
	report := NewStressReport()
	return &Stress{
		URL:         url,
		Method:      method,
		Concurrency: concurrency,
		Requests:    requests,
		Timeout:     timeout,
		Verbose:     verbose,
		Report:      report,
		VerifyTls:   verifyTls,
		mu:          sync.Mutex{},
	}
}

// Run the stress test
func (s *Stress) Run() error {
	fmt.Println("Running stress test...")
	s.run()
	return nil
}

// PrintReport prints the report
func (s *Stress) PrintReport() {
	fmt.Println("--- Report ---")
	fmt.Println("Requests:", s.Report.Requests)
	fmt.Println("Failed:", s.Report.Failed)
	fmt.Println("Succeeded:", s.Report.Succeeded)
	fmt.Println("TimedOut:", s.Report.TimedOut)
	fmt.Println("TotalTime:", s.Report.TotalTime, "ms")
	fmt.Println("AverageTime:", s.Report.AverageTime, "ms")
	fmt.Println("FastestTime:", s.Report.FastestTime, "ms")
	fmt.Println("SlowestTime:", s.Report.SlowestTime, "ms")
	fmt.Println("PercentageSucceeded:", s.Report.PercentageSucceeded, "%")
	fmt.Println("PercentageFailed:", s.Report.PercentageFailed, "%")
	fmt.Println("PercentageTimedOut:", s.Report.PercentageTimedOut, "%")
	fmt.Println("--- Requests per status code ---")
	for status, requests := range s.Report.StatusRequests {
		fmt.Println("Status", fmt.Sprint(status)+":", requests, "requests")
	}
}

// Run the stress test
func (s *Stress) run() {
	start := time.Now()

	// Create a wait group to wait for all requests to finish
	var wg sync.WaitGroup

	// Start the specified number of concurrent goroutines
	for i := 0; i < s.Concurrency; i++ {
		wg.Add(1)
		i := i

		go func() {
			defer wg.Done()
			// Run the specified number of requests
			for j := 0; j < s.Requests/s.Concurrency; j++ {
				s.runRequest(i + 1)
			}
		}()
	}

	// Wait for all requests to finish
	wg.Wait()

	// Time elapsed in seconds
	elapsed := time.Since(start).Milliseconds()

	// Set the total time taken
	s.Report.TotalTime = float64(elapsed)

	// Calculate AverageTime
	s.Report.AverageTime = s.Report.TotalTime / float64(s.Report.Requests)

	// Calculate PercentageSucceeded
	s.Report.PercentageSucceeded = float64(s.Report.Succeeded) / float64(s.Report.Requests) * 100

	// Calculate PercentageFailed
	s.Report.PercentageFailed = float64(s.Report.Failed) / float64(s.Report.Requests) * 100

	// Calculate PercentageTimedOut
	s.Report.PercentageTimedOut = float64(s.Report.TimedOut) / float64(s.Report.Requests) * 100

	fmt.Println("Finished stress test")
}

func (s *Stress) runRequest(concurrencyGroup int) {
	start := time.Now()

	// Create a new HTTP transport
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !s.VerifyTls},
	}

	// Create a new HTTP client
	client := &http.Client{
		Timeout:   time.Duration(s.Timeout) * time.Second,
		Transport: tr,
	}

	// Make the request
	req, err := http.NewRequest(s.Method, s.URL, nil)
	if err != nil {
		panic(err)
	}

	// Get the response
	res, err := client.Do(req)

	elapsed := time.Since(start).Milliseconds()

	// If Verbose, print the request
	if s.Verbose {
		fmt.Print(fmt.Sprint(concurrencyGroup) + " | " + fmt.Sprint(s.Report.Requests+1) + " " + s.Method + " " + s.URL)
		fmt.Println(" Time:", elapsed, "ms, Status:", res.StatusCode)
	}

	s.updateReport(res, err, elapsed)
}

// Update the report
func (s *Stress) updateReport(res *http.Response, err error, elapsed int64) {
	// Lock the mutex
	s.mu.Lock()
	defer s.mu.Unlock()

	if err != nil {
		fmt.Println(err)
		// Error making the request
		if strings.Contains(err.Error(), "connection refused") {
			// If the error is a connection refused, panic
			panic(err)
		}
		// If the error is a timeout, increment TimedOut
		if err.Error() == http.ErrHandlerTimeout.Error() {
			s.Report.TimedOut++
		}
		// Increment Failed
		s.Report.Failed++
	} else {
		// No error making the request
		if res.StatusCode != 200 {
			// If the status code is not 200, increment Failed
			s.Report.Failed++
		} else {
			// If the status code is 200, increment Succeeded
			s.Report.Succeeded++
		}
		// If the status code is not in StatusRequests, add it
		if _, ok := s.Report.StatusRequests[res.StatusCode]; !ok {
			s.Report.StatusRequests[res.StatusCode] = 0
		}
		// Increment the status code
		s.Report.StatusRequests[res.StatusCode]++ // Increment Succeeded
	}

	// Increment Requests
	s.Report.Requests++

	// If the time taken is faster than FastestTime, set FastestTime
	if elapsed < s.Report.FastestTime || s.Report.FastestTime == 0 {
		s.Report.FastestTime = elapsed
	}

	// If the time taken is slower than SlowestTime, set SlowestTime
	if elapsed > s.Report.SlowestTime {
		s.Report.SlowestTime = elapsed
	}
}
