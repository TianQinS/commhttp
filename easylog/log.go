package easylog

import (
	"fmt"

	"github.com/TianQinS/commhttp/config"
	"github.com/TianQinS/commhttp/utils/mail"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// normal config.
const (
	DEFAULT_COMPRESSED = false
)

var (
	ErrorWriter zapcore.WriteSyncer
	LogMgr      = make(Logs)
	conf        = &config.Conf.Log
)

type Logs map[string]*zap.Logger

// Obtain or generate new logger by key.
func (this Logs) Get(key string) *zap.Logger {
	if logger, ok := this[key]; ok {
		return logger
	}
	if logger, err := this.Register(
		key,
		fmt.Sprintf("%s/%s.log", conf.LogDir, key),
		conf.LogNormMbytes,
		conf.LogNormMaxDays,
		conf.LogNormMaxBackups,
		true,
		DEFAULT_COMPRESSED); err == nil {
		return logger
	} else {
		mail.SendError(err.Error())
	}
	return nil
}

func (this Logs) errorHook(entry zapcore.Entry) error {
	if entry.Level >= zap.ErrorLevel {
		mail.SendError(entry.Message)
	}
	return nil
}

func (this Logs) Register(key, filePath string, maxSize, maxAge, maxBackups int, localTime, compress bool) (*zap.Logger, error) {
	var logger *zap.Logger
	hook := lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackups,
		LocalTime:  localTime,
		Compress:   compress,
	}
	fileWriter := zapcore.AddSync(&hook)

	if config.Conf.Debug {
		core := zapcore.NewTee(
			zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
				fileWriter,
				zap.DebugLevel),
			zapcore.RegisterHooks(zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
				ErrorWriter,
				zap.ErrorLevel), this.errorHook),
		)
		logger = zap.New(core).WithOptions(zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	} else {
		core := zapcore.NewTee(
			zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
				fileWriter,
				zap.InfoLevel),
			zapcore.RegisterHooks(zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
				ErrorWriter,
				zap.ErrorLevel), this.errorHook),
		)
		logger = zap.New(core).WithOptions(zap.AddCaller())
	}
	this[key] = logger
	return logger, nil
}

func init() {
	hook := lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/error.log", conf.LogDir),
		MaxSize:    conf.LogErrMbytes,
		MaxAge:     conf.LogErrMaxDays,
		MaxBackups: conf.LogErrMaxBackups,
		LocalTime:  true,
		Compress:   true,
	}

	ErrorWriter = zapcore.AddSync(&hook)
}
