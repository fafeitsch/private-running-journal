package filebased

type Service struct {
	path string
}

func NewService(path string) *Service {
	return &Service{path: path}
}

type track struct {
	Name string `json:"name"`
}
