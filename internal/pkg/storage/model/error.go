package model

import "errors"

var (
	ErrNoImplement  = errors.New("filesystem not implement")
	ErrNoSupport    = errors.New("filesystem not support this operation")
	ErrFileNotExist = errors.New("file does not exist")
)
