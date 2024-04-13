package log

func Close() {
	for _, logger := range loggers {
		logger.Close()
	}
}
