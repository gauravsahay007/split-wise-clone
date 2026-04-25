package business

import (
	"errors"

	"github.com/gauravsahay007/split-wise-clone/models"
	"github.com/gauravsahay007/split-wise-clone/repository"
	"github.com/gauravsahay007/split-wise-clone/utils"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Repo *repository.Repo
}

func (s *Service) CreateUser(name string, password string) (models.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}
	return s.Repo.SaveUser(name, string(hashed)) //convert byte slice to string for DB storage
}

func (s *Service) CreateExpense(exp models.Expense) error {
	return s.Repo.SaveExpense(exp)
}

func (s *Service) simplifyDebts(netBalances map[int]float64) []models.Balance {
	type score struct {
		userID int
		amount float64
	}

	// Ignore tiny floating-point errors—only treat amounts greater than ₹0.01 as real debts or credits.
	var debtors, creditors []score
	for id, amt := range netBalances {
		if amt < -0.01 {
			debtors = append(debtors, score{userID: id, amount: -amt})
		} else if amt > 0.01 {
			creditors = append(creditors, score{userID: id, amount: amt})
		}
	}

	var results []models.Balance
	i, j := 0, 0

	// Match debtors with creditors greedily
	for i < len(debtors) && j < len(creditors) {
		// debtor: 50 creditor: 40
		settleAmount := debtors[i].amount // 50

		if creditors[j].amount < settleAmount {
			settleAmount = creditors[j].amount //creditor got what it needed
		}

		results = append(results, models.Balance{
			FromUser: debtors[i].userID,
			ToUser:   creditors[j].userID,
			Amount:   settleAmount,
		})

		debtors[i].amount -= settleAmount   // 50-40 = 10
		creditors[j].amount -= settleAmount //40-40 = 0

		// Move to next person if their balance is settled
		if debtors[i].amount < 0.01 {
			i++
		}
		if creditors[j].amount < 0.01 {
			j++
		}
	}

	return results
}

func (s *Service) GetBalances(groupID int) ([]models.Balance, error) {
	expenses, err := s.Repo.GetExpensesByGroup(groupID)
	if err != nil {
		return nil, err
	}

	netBalances := make(map[int]float64)
	for _, exp := range expenses {
		netBalances[exp.PaidBy] += exp.Amount
		share := exp.Amount / float64(len(exp.UserIDs))
		for _, uid := range exp.UserIDs {
			netBalances[uid] -= share
		}
	}
	return s.simplifyDebts(netBalances), nil
}

func (s *Service) AddMemberToGroup(groupID int, userID int) error {
	return s.Repo.AddUserToGroup(groupID, userID)
}

func (s *Service) Authenticate(id int, password string) (string, error) {
	user, err := s.Repo.GetUserByID(id)
	if err != nil {
		return "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)

	if err != nil {
		return "", errors.New("Invalid Credentials")
	}

	return utils.GenerateToken(user.ID)
}
