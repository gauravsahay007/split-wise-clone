package repository

import (
	"database/sql"
	"errors"

	"github.com/gauravsahay007/split-wise-clone/models"
	"github.com/lib/pq"
)

type Repo struct {
	DB *sql.DB
}

func (r *Repo) SaveUser(name string, password *string, email string, profilePic string) (models.User, error) {
	var u models.User
	query := `INSERT INTO users (name, password, email, profile_pic) 
	          VALUES ($1, $2, $3, $4) RETURNING id, name, email, profile_pic`

	err := r.DB.QueryRow(query, name, password, email, profilePic).Scan(
		&u.ID, &u.Name, &u.Email, &u.ProfilePic,
	)
	return u, err
}

func (r *Repo) GetUserByID(id int) (models.User, error) {
	var u models.User
	err := r.DB.QueryRow("SELECT id, name, email, profile_pic FROM users WHERE id = $1", id).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Password,
	)
	return u, err
}

func (r *Repo) GetUserWithHashedPassword(id int) (models.User, error) {
	var u models.User
	err := r.DB.QueryRow("SELECT id, name, password FROM users WHERE id = $1", id).Scan(
		&u.ID,
		&u.Name,
		&u.Password,
	)
	return u, err
}

func (r *Repo) SaveExpense(exp models.Expense) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	// 2. Insert Expense
	var expID int
	query := `INSERT INTO expenses(group_id, paid_by, amount, description, category, receipt_image, split_type) 
	          VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err = tx.QueryRow(query, exp.GroupID, exp.PaidBy, exp.Amount, exp.Description,
		exp.Category, exp.ReceiptImage, exp.SplitType).Scan(&expID)
	if err != nil {
		return err
	}

	// 3. Insert Participants
	if exp.SplitType == "manual" {
		for _, share := range exp.Shares {
			_, err = tx.Exec("INSERT INTO participants(expense_id, user_id, share_amount) VALUES($1, $2, $3)",
				expID, share.UserID, share.Amount)
			if err != nil {
				return err
			}
		}
	} else {
		if len(exp.UserIDs) == 0 {
			return errors.New("no participants")
		}
		shareAmt := exp.Amount / float64(len(exp.UserIDs))
		for _, uid := range exp.UserIDs {
			_, err = tx.Exec("INSERT INTO participants(expense_id, user_id, share_amount) VALUES($1, $2, $3)",
				expID, uid, shareAmt)
			if err != nil {
				return err
			}
		}
	}

	// 4. Commit everything if all steps succeeded
	return tx.Commit()
}

// // GetAllExpenses returns all expense records with participant IDs
// func (r *Repo) GetAllExpenses() ([]models.Expense, error) {
// 	query := `SELECT e.id, e.paid_by, e.amount, array_agg(p.user_id)
// 		FROM expenses e
// 		JOIN participants p ON e.id = p.expense_id
// 		GROUP BY e.id, e.paid_by, e.amount;`

// 	rows, err := r.DB.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	var expenses []models.Expense
// 	for rows.Next() {
// 		var e models.Expense
// 		var tempUserIDs pq.Int64Array // Use pq's specialized type for scanning
// 		// Scan Postgres INT[] into Go []int
// 		err := rows.Scan(&e.ID, &e.PaidBy, &e.Amount, &tempUserIDs)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// Convert []int64 to your model's []int
// 		e.UserIDs = make([]int, len(tempUserIDs))
// 		for i, v := range tempUserIDs {
// 			e.UserIDs[i] = int(v)
// 		}
// 		expenses = append(expenses, e)
// 	}
// 	return expenses, nil
// }

func (r *Repo) SaveGroup(name string, creatorID int) (models.Group, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return models.Group{}, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var g models.Group

	query := "INSERT INTO groups(name, created_by) VALUES($1, $2) RETURNING id, name"
	err = tx.QueryRow(query, name, creatorID).Scan(&g.ID, &g.Name)
	if err != nil {
		return models.Group{}, err
	}

	memberQuery := "INSERT INTO group_members(group_id, user_id) VALUES($1, $2)"
	_, err = tx.Exec(memberQuery, g.ID, creatorID)
	if err != nil {
		return models.Group{}, err
	}

	err = tx.Commit()
	if err != nil {
		return models.Group{}, err
	}

	return g, nil
}

func (r *Repo) GetExpensesByGroup(groupID int) ([]models.Expense, error) {
	query := `
		SELECT e.id, e.paid_by, e.amount, array_agg(p.user_id)
		FROM expenses e
		JOIN participants p ON e.id = p.expense_id
		WHERE e.group_id = $1
		GROUP BY e.id, e.paid_by, e.amount`

	rows, err := r.DB.Query(query, groupID)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			return
		}
	}()

	var expenses []models.Expense
	for rows.Next() {
		var e models.Expense
		var tempIDs pq.Int64Array
		if err := rows.Scan(&e.ID, &e.PaidBy, &e.Amount, &tempIDs); err != nil {
			return nil, err
		}
		e.UserIDs = make([]int, len(tempIDs))
		for i, v := range tempIDs {
			e.UserIDs[i] = int(v)
		}
		expenses = append(expenses, e)
	}
	return expenses, nil
}

func (r *Repo) AddUserToGroup(groupID int, userID int) error {
	_, err := r.DB.Exec(
		"INSERT INTO group_members (group_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		groupID, userID,
	)
	return err
}

// GetTotalPaidByUser calculates the sum of all expenses paid by this user
func (r *Repo) GetTotalPaidByUser(userID int) (float64, error) {
	var total float64
	query := "SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE paid_by = $1"
	err := r.DB.QueryRow(query, userID).Scan(&total)
	return total, err
}

// GetTotalOwedByUser calculates the sum of all shares assigned to this user
func (r *Repo) GetTotalOwedByUser(userID int) (float64, error) {
	var total float64
	query := "SELECT COALESCE(SUM(share_amount), 0) FROM participants WHERE user_id = $1"
	err := r.DB.QueryRow(query, userID).Scan(&total)
	return total, err
}
