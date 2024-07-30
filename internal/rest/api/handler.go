package api

const TEMP_IMAGE = "https://via.placeholder.com/300x200"

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}
