package business

import (
	"github.com/gauravsahay007/split-wise-clone/models"
	"github.com/gauravsahay007/split-wise-clone/repository"
)

type Service struct {
	Repo *repository.Repo
}

func (s *Service) CreateUser(name string) (models.User, error) {
	return s.Repo.SaveUser(name)
}

func (s *Service) CreateExpense(exp models.Expense) error {
	return s.Repo.SaveExpense(exp)
}
