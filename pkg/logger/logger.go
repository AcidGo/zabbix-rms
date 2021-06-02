package logger

import (
    "fmt"
    "io"
    "strings"
    "time"

    "github.com/sirupsen/logrus"
    nested "github.com/antonfisher/nested-logrus-formatter"
    rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

type ContextLogger = logrus.Entry

var logger *logrus.Logger

func init() {
    logger = logrus.New()
    logger.SetFormatter(&nested.Formatter{
        HideKeys:           true,
        NoColors:           true,
        TimestampFormat:    time.RFC3339,
    })
    logger.SetLevel(logrus.DebugLevel)
}

func LogFileSetting(filePath string) (error) {
    r, err := rotatelogs.New(filePath + "-%Y%m%d")
    if err != nil {
        return err
    }

    mw := io.MultiWriter(r)
    logger.SetOutput(mw)
    return nil
}

func LogLevelSetting(level string) (error) {
    var l logrus.Level
    var err error

    switch strings.ToLower(level) {
        case "panic":
            l = logrus.PanicLevel
        case "fatal":
            l = logrus.FatalLevel
        case "error":
            l = logrus.ErrorLevel
        case "warn":
            l = logrus.WarnLevel
        case "debug":
            l = logrus.DebugLevel
        case "trace":
            l = logrus.TraceLevel
        default:
            err = fmt.Errorf("not mean the log level %s", level)
    }

    if err != nil {
        return err
    }

    logger.SetLevel(l)
    return nil
}

func ReportCallerSetting(b bool) {
    logger.SetReportCaller(b)
}

func FitContext(mName string) (*logrus.Entry) {
    cl := logger.WithFields(logrus.Fields{
        "module":   mName,
    })
    return cl
}