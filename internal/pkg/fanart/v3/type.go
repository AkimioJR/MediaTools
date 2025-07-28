package fanart

import (
	"fmt"
)

type ErrorResponse struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"error message"`
}

type FanartError struct {
	err error
	msg string
}

func (e *FanartError) Error() string {
	if e.err == nil {
		return fmt.Sprintf("Fanart 错误: %s", e.msg)
	}
	return fmt.Sprintf("Fanart 错误: %s - %s", e.msg, e.err.Error())
}

func NewFanartError(msg string, err error) *FanartError {
	return &FanartError{
		err: err,
		msg: msg,
	}
}

type BaseInfo struct {
	ID    string `json:"id"`
	URL   string `json:"url"`
	Lang  string `json:"lang"`
	Likes string `json:"likes"`
}

type SeasonInfo struct {
	BaseInfo
	Season string `json:"season"`
}

type DiscInfo struct {
	BaseInfo
	Disc     string `json:"disc"`
	DiscType string `json:"disc_type"`
}

func (info BaseInfo) getLang() string {
	return info.Lang
}

type img interface {
	getLang() string
}
