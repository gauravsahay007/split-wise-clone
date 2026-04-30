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

func (r *Repo) GetUserGroups(userID int) ([]models.Group, error) {
	query := `
	SELECT 
    g.id AS group_id,
    g.name,
    g.created_at,
    COUNT(gm2.user_id) AS total_members
	FROM group_members gm1
	JOIN groups g ON gm1.group_id = g.id
	JOIN group_members gm2 ON gm1.group_id = gm2.group_id
	WHERE gm1.user_id = $1
	GROUP BY g.id, g.name, g.created_at;
	`
	var res []models.Group
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		rows.Close()
	}()

	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.ID, &group.Name, &group.CreatedAt, &group.Count); err != nil {
			return nil, err
		}
		res = append(res, group)
	}
	return res, nil
}

func (r *Repo) GetUserIDByEmail(email string) (int, error) {
	query := "SELECT id FROM users where email=$1;"
	var res int
	err := r.DB.QueryRow(query, email).Scan(&res)
	if err != nil {
		return -1, err
	}
	return res, nil
}

func (r *Repo) AddFriend(userId int, friendIds []string) error {
	query := `
	INSERT INTO friends (user_id, friend_user_id)
	SELECT LEAST($1, u.id), GREATEST($1, u.id)
	FROM users u WHERE u.email = ANY($2::text[])
	AND u.id <> $1
	ON CONFLICT (user_id, friend_user_id) DO NOTHING;
	`

	_, err := r.DB.Exec(query, userId, pq.Array(friendIds))
	return err
}

func (r *Repo) GetFriendsList(userId int) ([]models.User, error) {
	query := `
		SELECT u.name, u.email, u.profile_pic FROM (select 
		CASE 
			WHEN user_id = $1 THEN friend_user_id
			ELSE user_id
		END
		FROM friends WHERE user_id=$1 OR friend_user_id=9) AS f JOIN users u ON f.user_id = u.id;
	`

	rows, err := r.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}

	var res []models.User

	defer func() {
		rows.Close()
	}()

	for rows.Next() {
		var friend models.User
		err := rows.Scan(&friend.Name, &friend.Email, &friend.ProfilePic)
		if err != nil {
			return nil, err
		}

		res = append(res, friend)
	}

	return res, nil
}

func (r *Repo) SearchFriendsInAGroup(userId int, searchString string, groupId int) ([]models.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.profile_pic
		FROM friends f
		JOIN users u 
		  ON u.id = CASE 
					 WHEN f.user_id = $1 THEN f.friend_user_id
					 ELSE f.user_id
				   END
		WHERE $1 IN (f.user_id, f.friend_user_id)
		  AND (
			u.name ILIKE '%' || $2 || '%'
			OR u.email ILIKE '%' || $2 || '%'
		  )
		  AND NOT EXISTS (
			SELECT 1 
			FROM group_members gm
			WHERE gm.group_id = $3
			  AND gm.user_id = u.id
		  );
	`

	rows, err := r.DB.Query(query, userId, searchString, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.User

	for rows.Next() {
		var friend models.User
		if err := rows.Scan(
			&friend.ID,
			&friend.Name,
			&friend.Email,
			&friend.ProfilePic,
		); err != nil {
			return nil, err
		}
		res = append(res, friend)
	}

	return res, nil
}

func (r *Repo) GetGroupMembers(userId, groupId int) ([]models.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.profile_pic
		FROM group_members gm
		JOIN users u ON u.id = gm.user_id
		JOIN group_members gm_check 
		  ON gm_check.group_id = gm.group_id 
		 AND gm_check.user_id = $1
		WHERE gm.group_id = $2;
	`

	rows, err := r.DB.Query(query, userId, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.ProfilePic,
		); err != nil {
			return nil, err
		}
		res = append(res, user)
	}

	return res, nil
}

func (r *Repo) IsUserInGroup(userID, groupID int) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM group_members 
			WHERE user_id = $1 AND group_id = $2
		)
	`
	err := r.DB.QueryRow(query, userID, groupID).Scan(&exists)
	return exists, err
}

func (r *Repo) GetGroupExpenses(groupID int) ([]models.Expense, error) {
	query := `
	SELECT 
		e.id,
		e.paid_by,
		e.amount,
		e.group_id,
		e.description,
		e.category,
		e.receipt_image,
		e.split_type,
		e.created_at,
		p.user_id,
		p.share_amount
	FROM expenses e
	LEFT JOIN participants p ON e.id = p.expense_id
	WHERE e.group_id = $1
	ORDER BY e.created_at DESC;
	`

	rows, err := r.DB.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expMap := make(map[int]*models.Expense)

	for rows.Next() {
		var (
			exp         models.Expense
			userID      *int
			shareAmount *float64
		)

		err := rows.Scan(
			&exp.ID,
			&exp.PaidBy,
			&exp.Amount,
			&exp.GroupID,
			&exp.Description,
			&exp.Category,
			&exp.ReceiptImage,
			&exp.SplitType,
			&exp.CreatedAt,
			&userID,
			&shareAmount,
		)
		if err != nil {
			return nil, err
		}

		if _, ok := expMap[exp.ID]; !ok {
			exp.UserIDs = []int{}
			exp.Shares = []models.ExpenseShare{}
			expMap[exp.ID] = &exp
		}

		if userID != nil {
			expMap[exp.ID].UserIDs = append(expMap[exp.ID].UserIDs, *userID)

			if shareAmount != nil {
				expMap[exp.ID].Shares = append(
					expMap[exp.ID].Shares,
					models.ExpenseShare{
						UserID: *userID,
						Amount: *shareAmount,
					},
				)
			}
		}
	}

	result := []models.Expense{}
	for _, v := range expMap {
		result = append(result, *v)
	}

	return result, nil
}
