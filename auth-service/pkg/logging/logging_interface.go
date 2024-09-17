package logging

type Logger interface {
	LogInfo(message string)
	LogError(message string)
	LogWarn(message string)
	LogDebug(message string)
	LogPanic(message string)
}
