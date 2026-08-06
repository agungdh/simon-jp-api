package service

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrForbidden      = errors.New("forbidden")
	ErrModuleNotFound = errors.New("module not found")
	ErrNoFileAttached = errors.New("no file attached")
	ErrFileNotFound   = errors.New("file not found")
)
