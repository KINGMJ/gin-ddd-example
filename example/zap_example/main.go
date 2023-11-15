package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// 创建 Logger
	logger := zap.New(getZapCore())
	// 使用 Logger 记录日志
	for i := 0; i < 1000; i++ {
		logger.Info("This is an info log", zap.String("key", "value"))
		logger.Error("This is an error log", zap.Error(errors.New("something went wrong")))
	}
}

// 日志配置
var LOG = struct {
	Env        string // 环境：dev/prod
	Format     string // 格式：json 或者其他格式
	RootDir    string // 日志文件存放位置
	Filename   string // 日志文件名
	MaxSize    int    // 日志文件最大大小(M)
	MaxBackups int    // 旧文件的最大个数
	MaxAge     int    // 旧文件的最大保留天数
	Compress   bool   // 是否压缩

}{"dev", "json", "logs", "app.log", 1, 3, 28, true}

func getZapCore() zapcore.Core {
	level := zap.NewAtomicLevel()
	level.SetLevel(zapcore.InfoLevel)

	var encoder zapcore.Encoder
	// 调整编码器默认配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format("[" + "2006-01-02 15:04:05.000" + "]"))
	}
	encoderConfig.EncodeLevel = func(l zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(LOG.Env + "." + l.String())
	}

	// 设置编码器
	if LOG.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	return zapcore.NewCore(encoder, getLogWriter(), level)
}

func getLogWriter() zapcore.WriteSyncer {
	today := time.Now().Format("2006-01-02")
	file := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s-%s.log", LOG.RootDir, LOG.Filename, today),
		MaxSize:    LOG.MaxSize,
		MaxBackups: LOG.MaxBackups,
		MaxAge:     LOG.MaxAge,
		Compress:   LOG.Compress,
	}
	return zapcore.AddSync(file)
}
