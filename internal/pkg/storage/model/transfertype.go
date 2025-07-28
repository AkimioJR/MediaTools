package model

type TransferType uint8

const (
	TransferUnknown  TransferType = iota // 未知传输类型
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
		return "Unknown"
	}
}
