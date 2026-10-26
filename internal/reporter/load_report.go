package reporter

import (
	"fmt"

	"stress/internal/model"
)

func(r *reporter) PrintLoadReport(result model.TestResult) {
	fmt.Fprintln(r.writer, "\nLoadLab Test Report")
	fmt.Fprintln(r.writer, "-------------------")
	fmt.Fprintf(r.writer, "Target:       %s\n", result.Url)
	fmt.Fprintf(r.writer, "Duration:     %s\n", result.Duration)
	fmt.Fprintf(r.writer, "Requests:     %d\n", result.TotalRequests)
	fmt.Fprintf(r.writer, "Successful:   %d\n", result.SuccessCount)
	fmt.Fprintf(r.writer, "Failed:       %d\n", result.FailureCount)
	fmt.Fprintf(r.writer, "Avg latency:  %d\n", result.AverageLatency)
	fmt.Fprintf(r.writer, "Requests/sec: %.2f\n", result.ReqPerSec)
}
