package internal

import (
	_interface "github.com/saur4ig/file-storage/internal/services/interface"
)

type s3Service struct {
}

func NewS3Service() _interface.FileStorage {
	return &s3Service{}
}
