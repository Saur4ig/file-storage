package _interface

// FileStorage is an interface that defines methods to store, generateUlr and remove files
type FileStorage interface {
	// UploadFile uploads a file to the storage.
	UploadFile(fileData []byte, fileName string) (fileURL string, err error)

	// GeneratePreSignedURL generates a pre-signed URL for the specified file operation.
	GeneratePreSignedURL(fileName string) (preSignedURL string, err error)

	// DeleteFile deletes a file from the storage.
	DeleteFile(fileURL string) error
}
