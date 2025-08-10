package logger

import (
	"log"
	"os"
)

func InitLogger() {
	logFile, err := os.OpenFile("eth_trading.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
}
