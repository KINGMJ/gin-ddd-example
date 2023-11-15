package logs

import (
	"fmt"
	"gin-ddd-example/pkg/config"
	"gin-ddd-example/pkg/utils"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

var (
	level   zapcore.Level // zap 日志等级
	options []zap.Option  // zap 配置项
)

func InitLog(conf config.Config) {
	// 创建日志根目录
	createRootDir(&conf)
	// 设置日志等级
	setLogLevel(&conf)
	if conf.LogsConf.ShowLine {
		options = append(options, zap.AddCaller())
	}
	// 初始化 zap
	Log = zap.New(getZapCore(&conf), options...)
}

// 创建日志根目录
//
//	@param conf
func createRootDir(conf *config.Config) {
	if ok, _ := utils.PathExists(conf.LogsConf.RootDir); !ok {
		_ = os.Mkdir(conf.LogsConf.RootDir, os.ModePerm)
	}
}

// 设置日志等级
//
//	@param conf
func setLogLevel(conf *config.Config) {
	switch conf.LogsConf.Level {
	case "debug":
		level = zap.DebugLevel
		options = append(options, zap.AddStacktrace(level))
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
		options = append(options, zap.AddStacktrace(level))
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}
}

func getZapCore(config *config.Config) zapcore.Core {
	var encoder zapcore.Encoder
	// 调整编码器默认配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format("[" + "2006-01-02 15:04:05.000" + "]"))
	}
	encoderConfig.EncodeLevel = func(l zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		str := fmt.Sprintf("[%s.%s]", config.AppConf.Env, l.String())
		encoder.AppendString(str)
	}

	// 设置编码器
	if config.LogsConf.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	return zapcore.NewCore(encoder, getLogWriter(config), level)
}

func getLogWriter(config *config.Config) zapcore.WriteSyncer {
	// 按天进行日志分割
	today := time.Now().Format("2006-01-02")
	file := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s-%s.log", config.LogsConf.RootDir, config.LogsConf.Filename, today),
		MaxSize:    config.LogsConf.MaxSize,
		MaxBackups: config.LogsConf.MaxBackups,
		MaxAge:     config.LogsConf.MaxAge,
		Compress:   config.LogsConf.Compress,
	}
	return zapcore.AddSync(file)
}
