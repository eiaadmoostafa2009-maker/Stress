package reporter

import (
	"io"
	"os"
	"stress/internal/model"
)

type Reporter interface{
    PrintLoadReport(result model.TestResult)
	PrintReqReport(result model.RequestResult)
}

type reporter struct{
	writer io.Writer
}

func NewReporter(writer io.Writer) Reporter{
	if writer == nil {
		writer = os.Stdout
	}
	return &reporter{
		writer: writer,
	}
}