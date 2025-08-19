package storage

import (
	"encoding/json"
	"strings"
)

type TransferType uint8

const (
	TransferUnknown  TransferType = iota // 未知转移类型
	TransferCopy                         // 复制
	TransferMove                         // 移动
	TransferLink                         // 硬链接
	TransferSoftLink                     // 软链接
)

func (t TransferType) String() string {
	switch t {
	case TransferCopy:
		return "Copy"
	case TransferMove:
		return "Move"
	case TransferLink:
		return "Link"
	case TransferSoftLink:
		return "SoftLink"
	default:
		return "UnknownTransferType"
	}
}

func (t TransferType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *TransferType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*t = ParseTransferType(s)
	return nil
}

func ParseTransferType(s string) TransferType {
	switch strings.ToLower(s) {
	case "copy":
		return TransferCopy
	case "move":
		return TransferMove
	case "link":
		return TransferLink
	case "softlink":
		return TransferSoftLink
	default:
		return TransferUnknown
	}
}

func (t TransferType) MarshalYAML() (any, error) {
	return t.String(), nil
}

func (t *TransferType) UnmarshalYAML(unmarshal func(any) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	*t = ParseTransferType(s)
	return nil
}
