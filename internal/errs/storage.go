package errs

import "errors"

var (
	ErrStorageProvideNoImplement = errors.New("storage provider not implement")
	ErrStorageProvideNoSupport   = errors.New("storage provider not support this operation")
	ErrStorageProviderNotFound   = errors.New("storage provider not found")

	ErrFileNotFound  = errors.New("file not found")
	ErrNotADirectory = errors.New("not a directory")
)
