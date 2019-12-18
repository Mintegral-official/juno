package log

type Logger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	InfoF(args ...interface{})
	WarnF(args ...interface{})
}
