package config

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type LogLevel uint32

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warning"
	case LevelError:
		return "error"
	default:
		return "unknown"
	}
}

func (l *LogLevel) MarshalSting() string {
	return l.String()
}

func (l *LogLevel) UnmarshalString(str string) error {
	switch str {
	case "debug":
		*l = LevelDebug
	case "info":
		*l = LevelInfo
	case "warning":
		*l = LevelWarn
	case "error":
		*l = LevelError
	default:
		return fmt.Errorf("unknown log level: %s", str)
	}
	return nil
}

func (l LogLevel) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

func (l *LogLevel) UnmarshalText(text []byte) error {
	return l.UnmarshalString(string(text))
}

func (l LogLevel) MarshalJSON() ([]byte, error) {
	return []byte(`"` + l.String() + `"`), nil
}

func (l *LogLevel) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	return l.UnmarshalString(str)
}

func (l LogLevel) MarshalYAML() (any, error) {
	return l.String(), nil
}

func (l *LogLevel) UnmarshalYAML(node *yaml.Node) error {
	var str string
	if err := node.Decode(&str); err != nil {
		return err
	}
	return l.UnmarshalString(str)
}

func (l LogLevel) ToLogrusLevel() logrus.Level {
	switch l {
	case LevelDebug:
		return logrus.DebugLevel
	case LevelInfo:
		return logrus.InfoLevel
	case LevelWarn:
		return logrus.WarnLevel
	case LevelError:
		return logrus.ErrorLevel
	default:
		return logrus.InfoLevel // 默认返回 Info 级别
	}
}
