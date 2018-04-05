package log

import (
	"bytes"
	"fmt"

	"github.com/sirupsen/logrus"
)

type Formatter struct {
	QuoteEmptyFields bool
	File             string
	Line             int
	TimestampFormat  string
	LogFormat        string
}

const (
	defaultTimestampFormat = "2006-01-02_15:04:05.000000"
	defaultLogFormat       = "[%lvl%]: %time% - %msg%"
)

func (f *Formatter) levelString(level logrus.Level) string {
	switch level {
	case logrus.DebugLevel:
		return "DEBG"
	case logrus.InfoLevel:
		return "INFO"
	case logrus.WarnLevel:
		return "WARN"
	case logrus.ErrorLevel:
		return "ERRO"
	case logrus.FatalLevel:
		return "FATL"
	case logrus.PanicLevel:
		return "PANC"
	}

	return "UNKN"
}

func (f *Formatter) appendKeyValue(b *bytes.Buffer, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}
	b.WriteString(stringVal)
}

// TODO log format
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	output := f.LogFormat
	if output == "" {
		output = defaultLogFormat
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	f.appendKeyValue(b, entry.Time.Format(timestampFormat))
	f.appendKeyValue(b, f.levelString(entry.Level))

	for _, key := range keys {
		f.appendKeyValue(b, entry.Data[key])
	}

	if entry.Message != "" {
		f.appendKeyValue(b, entry.Message)
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}
