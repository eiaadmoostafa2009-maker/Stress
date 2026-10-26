package reporter

import (
	"fmt"

	"stress/internal/model"
)

func(r *reporter) PrintReqReport(result model.RequestResult) {
	fmt.Fprintln(r.writer, "\nLoadLab Test Report")
	fmt.Fprintln(r.writer, "-------------------")
	fmt.Fprintf(r.writer, "StatusCode:       %d\n", result.StatusCode)
	fmt.Fprintf(r.writer, "Latency:   %d\n", result.Latency)
	fmt.Fprintf(r.writer, "Failed:       %t\n", result.Failed)
	fmt.Fprintf(r.writer, "Error:  %s\n", result.Error)
}