package log

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

const (
	// DefaultLogFilename is the default log filename.
	DefaultLogFilename = "debug.log"
)

func init() {

	// Default logger is silent, logs to ioutil.Discard
	log.SetOutput(ioutil.Discard)
	log.SetPrefix("")
	log.SetFlags(log.Ldate | log.Ltime)

}

// Setup configures the global logger
func Setup(verboseFlag bool, logFlag string) {

	// If logFlag is not empty then log to the file specified.
	if logFlag != "" {

		logFile, err := os.OpenFile(logFlag, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			fmt.Printf("Error: Could not open log file: %s\n.", logFlag)
			os.Exit(1)
		}

		log.SetOutput(logFile)

	}

	// If verboseFlag is true then log to Stdout.
	if verboseFlag {

		log.SetOutput(os.Stdout)

	}

	// If logFlag is not empty and verboseFlag is true then log to both logFile and Stdout.
	if logFlag != "" && verboseFlag {

		logFile, err := os.OpenFile(logFlag, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			fmt.Printf("Error: Could not open log file: %s\n.", logFlag)
			os.Exit(1)
		}

		mw := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(mw)

	}

}
