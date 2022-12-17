package repository

type ImgRepository interface {
	Save(id string, resolution int, bytes []byte) error
}

type Repository struct {
	ImgRepository
}

func NewRepository() *Repository {
	return &Repository{&TempImgRepo{}}
}

type TempImgRepo struct {
}

func (r *TempImgRepo) Save(id string, resolution int, bytes []byte) error {
	return nil
}
