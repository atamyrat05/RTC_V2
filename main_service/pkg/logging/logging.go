package package_log

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		_, err := w.Write([]byte(line))
		if err != nil {
			return err
		}
	}
	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}


type Logger struct {
	*logrus.Entry
}

func GetLogger(filePath string, fileName string) *Logger {
	l := logrus.New()
	l.SetReportCaller(true)
	l.SetFormatter(&logrus.JSONFormatter{})

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err := os.MkdirAll(filePath, 0755)
		if err != nil {
			panic(err)
		}
	}

	allFile, err := os.OpenFile(filePath+"/"+fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}

	l.SetOutput(io.Discard)
	l.AddHook(&writerHook{
		Writer:    []io.Writer{allFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	l.SetLevel(logrus.TraceLevel)

	return &Logger{logrus.NewEntry(l)}
}
