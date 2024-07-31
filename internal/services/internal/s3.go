package internal

const MOCKED_URL1 = "https://picsum.photos/100"
const MOCKED_URL2 = "https://picsum.photos/50"

// No real file storage is implemented, it's kinda "mocked" methods to simulate working with s3
// No bucket logic implemented as well, assuming that it should be present in real project.

// UploadFile responsible for uploading file to s3
func (s *s3Service) UploadFile(fileData []byte, fileName string) (fileURL string, err error) {
	return MOCKED_URL1, nil
}

// GeneratePreSignedURL responsible for reserving url for a file
func (s *s3Service) GeneratePreSignedURL(fileName string) (preSignedURL string, err error) {
	return MOCKED_URL2, nil
}

// DeleteFile responsible for removing file from storage
func (s *s3Service) DeleteFile(fileURL string) error {
	return nil
}
