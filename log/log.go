package log

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"bxs/config"
)

var Logger *zap.Logger

func init() {
	InitLoggerForTest()
}

func InitLoggerForTest() {
	Logger, _ = zap.NewDevelopment()
	return
}

func InitLogger() {
	if !config.G.Log.Async {
		Logger, _ = zap.NewDevelopment()
		return
	}

	buffer := &zapcore.BufferedWriteSyncer{
		Size:          config.G.Log.AsyncBufferSizeByByte,
		FlushInterval: time.Second * time.Duration(config.G.Log.AsyncFlushIntervalBySecond),
		WS:            os.Stdout,
	}
	writeSyncer := zapcore.AddSync(buffer)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logLevel, err := zapcore.ParseLevel(config.G.Log.Level)
	if err != nil {
		panic(fmt.Sprintf("log level error: %v, level:[%s]", err, config.G.Log.Level))
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writeSyncer,
		logLevel,
	)

	Logger = zap.New(core)
}
