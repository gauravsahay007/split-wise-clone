package repository

import (
	"database/sql"

	"github.com/gauravsahay007/split-wise-clone/models"
	"github.com/lib/pq"
)

type Repo struct {
	DB *sql.DB
}

func (r *Repo) SaveUser(name string) (models.User, error) {
	var u models.User
	query := "INSERT INTO users(name) VALUES($1) RETURNING id, name"
	err := r.DB.QueryRow(query, name).Scan(&u.ID, &u.Name)
	return u, err
}

func (r *Repo) SaveExpense(exp models.Expense) error {
	//Start transaction
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var expID int
	query := "INSERT INTO expenses(group_id, paid_by, amount) VALUES($1, $2, $3) RETURNING id"
	err = tx.QueryRow(query, exp.GroupID, exp.PaidBy, exp.Amount).Scan(&expID)
	if err != nil {
		return err
	}

	for _, uid := range exp.UserIDs {
		query := "INSERT INTO participants(expense_id, user_id) VALUES($1, $2)"
		_, err = tx.Exec(query, expID, uid)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	return nil
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

func (r *Repo) SaveGroup(name string) (models.Group, error) {
	var g models.Group
	err := r.DB.QueryRow("INSERT INTO groups(name) VALUES($1) RETURNING id, name", name).Scan(&g.ID, &g.Name)
	return g, err
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
	defer rows.Close()

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
