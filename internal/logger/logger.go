package logger

import (
	"github.com/bernardoms/StarWarsPlanetAPI-GO/config"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

type Logger struct {
	Logger *logrus.Logger
}

func NewLogger(config config.LoggerConfig) *Logger {
	log := new(Logger)
	log.Logger = logrus.New()
	switch strings.ToLower(config.Level) {

	case "info":
		log.Logger.SetLevel(logrus.InfoLevel)
		break
	case "debug":
		log.Logger.SetLevel(logrus.DebugLevel)
	case "error":
		log.Logger.SetLevel(logrus.ErrorLevel)
	default:
		log.Logger.SetLevel(logrus.InfoLevel)
	}

	log.Logger.SetOutput(os.Stdout)
	log.Logger.SetFormatter(&logrus.JSONFormatter{})
	return log
}

func (l *Logger) LogWithFields(req *http.Request, level string, fields map[string]interface{}, message string) {

	if fields == nil {
		fields = make(map[string]interface{})
	}

	if req != nil {
		fields["method"] = req.Method
		fields["path"] = req.URL.Path
		fields["queryParam"] = req.URL.Query()
		fields["header"] = req.Header
	}

	switch strings.ToLower(level) {
	case "info":
		l.Logger.WithFields(fields).Info(message)
		break
	case "error":
		l.Logger.WithFields(fields).Error(message)
		break
	case "warn":
		l.Logger.WithFields(fields).Warn(message)
		break
	case "debug":
		l.Logger.WithFields(fields).Debug(message)
		break
	default:
		l.Logger.WithFields(fields).Info(message)
	}
}
