package services

import (
	_interface "github.com/saur4ig/file-storage/internal/services/interface"
	"github.com/saur4ig/file-storage/internal/services/internal"
)

func NewS3Service() _interface.FileStorage {
	return internal.NewS3Service()
}
