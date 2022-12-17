package repository

type ImgRepository interface {
	GetImage(id string, resolution int) ([]byte, error)
}

type Repository struct {
	ImgRepository
}

func NewRepository() *Repository {
	return &Repository{&TempImgRepo{}}
}

type TempImgRepo struct {
}

func (r *TempImgRepo) GetImage(id string, resolution int) ([]byte, error) {
	return nil, nil
}
