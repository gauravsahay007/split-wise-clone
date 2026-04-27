package business

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gauravsahay007/split-wise-clone/auth"
	"github.com/gauravsahay007/split-wise-clone/models"
	"github.com/gauravsahay007/split-wise-clone/repository"
	"github.com/gauravsahay007/split-wise-clone/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Repo *repository.Repo
}

func (s *Service) CreateUser(name string, password string, email string, profilePic string) (models.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}
	if profilePic == "" {
		profilePic = "https://cdn-icons-png.flaticon.com/512/4140/4140048.png"
	}
	hashedStr := string(hashed)
	return s.Repo.SaveUser(name, &hashedStr, email, profilePic) //convert byte slice to string for DB storage
}

func (s *Service) CreateGroup(name string, creatorID int) (models.Group, error) {
	return s.Repo.SaveGroup(name, creatorID)
}

func (s *Service) CreateExpense(exp models.Expense) error {
	if exp.Description == "" {
		exp.Description = "Uncategorized Expense"
	}

	if exp.Category == "" {
		exp.Category = "General"
	}

	if exp.ReceiptImage == "" {
		exp.ReceiptImage = "https://cdn-icons-png.flaticon.com/512/3135/3135679.png"
	}
	if exp.SplitType == "manual" {
		totalShares := 0.0
		for _, share := range exp.Shares {
			totalShares += share.Amount
		}
		// Safety check: shares must equal total amount
		if totalShares != exp.Amount {
			return fmt.Errorf("sum of shares (%v) does not equal total amount (%v)", totalShares, exp.Amount)
		}
	}

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

func (s *Service) GetUserOverallSummary(userID int) (map[string]float64, error) {
	paid, errP := s.Repo.GetTotalPaidByUser(userID)
	owed, errO := s.Repo.GetTotalOwedByUser(userID)

	if errP != nil || errO != nil {
		return nil, fmt.Errorf("failed to calculate financial summary")
	}

	return map[string]float64{
		"total_owed_to_you": paid,
		"total_you_owe":     owed,
		"net_balance":       paid - owed,
	}, nil
}

type GoogleUserInterface struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

func (s *Service) GoogleCallback(code string) (map[string]interface{}, error) {
	token, err := auth.GoogleConfig().Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	client := auth.GoogleConfig().Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch user info")
	}

	var googleUser GoogleUserInterface
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, err
	}

	if !googleUser.VerifiedEmail {
		return nil, errors.New("Email not verified by Google")
	}

	googleID := googleUser.ID
	email := googleUser.Email
	name := googleUser.Name
	picture := googleUser.Picture

	user, err := s.Repo.GetUserByProvider("google", googleID)
	if err != nil {
		return nil, err
	}

	if user != nil {
		token, err := utils.GenerateToken(user.ID)
		if err != nil {
			return nil, err
		}

		return gin.H{"token": token}, nil
	}

	user, err = s.Repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		err = s.Repo.AddAuthIdentity(user.ID, "google", googleID)
		if err != nil {
			return nil, err
		}
		token, err := utils.GenerateToken(user.ID)
		if err != nil {
			return nil, err
		}

		return gin.H{"token": token}, nil
	}

	newUser, err := s.Repo.SaveUser(name, nil, email, picture)
	if err != nil {
		return nil, err
	}

	err = s.Repo.AddAuthIdentity(newUser.ID, "google", googleID)
	if err != nil {
		return nil, err
	}

	newToken, err := utils.GenerateToken(newUser.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{"token": newToken}, nil
}
