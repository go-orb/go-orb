package log

var GlobalLogger Logger

func init() {
	GlobalLogger = newZero()
	c := NewConfig()
	c.SetLevel("trace")
	if err := GlobalLogger.Init(c, nil); err != nil {
		panic(err)
	}
}

func Trace() Event { return GlobalLogger.Trace() }
func Debug() Event { return GlobalLogger.Debug() }
func Info() Event  { return GlobalLogger.Info() }
func Warn() Event  { return GlobalLogger.Warn() }
func Err() Event   { return GlobalLogger.Err() }
func Fatal() Event { return GlobalLogger.Fatal() }
func Panic() Event { return GlobalLogger.Panic() }
