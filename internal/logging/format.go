package logging

import (
	"MediaTools/constants"
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Formater struct{}

func (f *Formater) Format(entry *logrus.Entry) ([]byte, error) {
	// 根据日志级别设置颜色
	var colorCode uint8
	switch entry.Level {
	case logrus.DebugLevel:
		colorCode = constants.ColorBlue
	case logrus.InfoLevel:
		colorCode = constants.ColorGreen
	case logrus.WarnLevel:
		colorCode = constants.ColorYellow
	case logrus.ErrorLevel:
		colorCode = constants.ColorRed
	default:
		colorCode = constants.ColorGray
	}

	// 设置文本Buffer
	var b *bytes.Buffer
	if entry.Buffer == nil {
		b = &bytes.Buffer{}
	} else {
		b = entry.Buffer
	}
	// 时间格式化
	formatTime := entry.Time.Format(time.DateTime)

	fmt.Fprintf(
		b,
		"\033[3%dm【%s】\033[0m %s | %s:%d - %s\n", // 长度需要算是上控制字符的长度
		colorCode,
		strings.ToUpper(entry.Level.String()),
		formatTime,
		entry.Caller.Function,
		entry.Caller.Line,
		entry.Message,
	)
	return b.Bytes(), nil
}
