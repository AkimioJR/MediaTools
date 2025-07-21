package themoviedb

import "fmt"

type TMDBError struct {
	err error
	msg string
}

func (e *TMDBError) Error() string {
	if e.err == nil {
		return fmt.Sprintf("TMDB 错误: %s", e.msg)
	}
	return fmt.Sprintf("TMDB 错误: %s - %s", e.msg, e.err.Error())
}

func NewTMDBError(err error, msg string) *TMDBError {
	return &TMDBError{
		err: err,
		msg: msg,
	}
}

type ErrorResponse struct {
	StatusCode    int32  `json:"status_code"`
	StatusMessage string `json:"status_message"`
	Success       bool   `json:"success"`
}
