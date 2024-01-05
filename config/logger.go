// config/logger.go

package config

import (
	"fmt"
	"log"
)

const (
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorReset  = "\033[0m"
)

type Logger interface {
	Log(message string)
}

type ServiceLogger struct {
	Service     string
	ColorPrefix string
}

func (l *ServiceLogger) Log(message string) {
	log.Printf("%s%s%s\t%s", l.ColorPrefix, fmt.Sprintf("[%s]", l.Service), ColorReset, message)
}
