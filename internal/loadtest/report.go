package loadtest

import (
	"fmt"
	"io"

	"stress/internal/model"
)

func PrintReport(w io.Writer, result model.TestResult) {
	fmt.Fprintln(w, "\nLoadLab Test Report")
	fmt.Fprintln(w, "-------------------")
	fmt.Fprintf(w, "Target:       %s\n", result.Url)
	fmt.Fprintf(w, "Duration:     %s\n", result.Duration)
	fmt.Fprintf(w, "Requests:     %d\n", result.TotalRequests)
	fmt.Fprintf(w, "Successful:   %d\n", result.SuccessCount)
	fmt.Fprintf(w, "Failed:       %d\n", result.FailureCount)
	fmt.Fprintf(w, "Avg latency:  %d\n", result.AverageLatency)
	fmt.Fprintf(w, "Requests/sec: %.2f\n", result.ReqPerSec)
}
