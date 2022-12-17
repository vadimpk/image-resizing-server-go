package services

import (
	"github.com/vadimpk/image-resizing-server-go/internal/api/repository"
	"log"
	"time"
)

type Service interface {
	Optimize(body []byte, headers map[string]interface{})
}

type Optimizer struct {
	repo *repository.Repository
}

func NewOptimizer(repo *repository.Repository) *Optimizer {
	return &Optimizer{repo: repo}
}

func (o *Optimizer) Optimize(body []byte, headers map[string]interface{}) {
	log.Println("optimizing...")
	time.Sleep(10 * time.Second)
}
