package log

import (
	"fmt"
	"os"
)

func FatalOnError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stdout, "FATAL: %s\n", err.Error())
		os.Exit(1)
	}
}

func Info(format string, args ...interface{}) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stdout, "INFO: %s\n", format)
	} else {
		fmt.Fprintf(os.Stdout, "INFO: %s\n", fmt.Sprintf(format, args...))
	}
}

func Error(format string, args ...interface{}) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stdout, "ERROR: %s\n", format)
	} else {
		fmt.Fprintf(os.Stdout, "ERROR: %s\n", fmt.Sprintf(format, args...))
	}
}
