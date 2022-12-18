package services

import (
	"fmt"
	"github.com/vadimpk/image-resizing-server-go/internal/optimizer/repository"
	"image"
	"log"
)

type Service interface {
	Optimize(body []byte, headers map[string]interface{})
}

type Resizer interface {
	Resize(body []byte, imgType string, scales []float64) ([]image.Image, error)
}

type Optimizer struct {
	Resizer
	repo repository.Repository
}

func NewOptimizer(repo repository.Repository) *Optimizer {
	return &Optimizer{repo: repo, Resizer: &NFNTResizer{}}
}

func (o *Optimizer) Optimize(body []byte, headers map[string]interface{}) {

	id := fmt.Sprintf("%v", headers["id"])
	imgType := fmt.Sprintf("%v", headers["img-type"])
	scales := []float64{1, .75, .5, .25}

	log.Printf("started optimizing %s", id)

	imgs, err := o.Resize(body, imgType, scales)
	if err != nil {
		log.Printf("couldn't resize img with id %v: [%s]\n", headers["id"], err)
	}

	for i, img := range imgs {
		err := o.repo.Save(id, img, int(scales[i]*100), imgType)
		if err != nil {
			log.Printf("couldn't save image with id %s to repository: [%s]\n", id, err)
		}
	}

	log.Printf("finished optimizing %s", id)
}
