package services

import (
	"github.com/pabloantipan/hobe-api-boilerplate/internal/models"
	"github.com/pabloantipan/hobe-api-boilerplate/internal/repositories/interfaces"
)

type PlayerService struct {
	repo interfaces.PlayerRepository
}

func NewPlayerService(repo interfaces.PlayerRepository) PlayerServiceInterface {
	return &PlayerService{repo: repo}
}

type PlayerServiceInterface interface {
	Create(player models.Player) (models.Player, error)
	GetByID(id string) (models.Player, error)
	GetAll() ([]models.Player, error)
	Update(player models.Player) (models.Player, error)
	Delete(id string) error
}

func (s *PlayerService) Create(player models.Player) (models.Player, error) {
	return s.repo.Create(player)
}

func (s *PlayerService) GetByID(id string) (models.Player, error) {
	return s.repo.GetByID(id)
}

func (s *PlayerService) GetAll() ([]models.Player, error) {
	return s.repo.GetAll()
}

func (s *PlayerService) Update(player models.Player) (models.Player, error) {
	return s.repo.Update(player)
}

func (s *PlayerService) Delete(id string) error {
	return s.repo.Delete(id)
}
