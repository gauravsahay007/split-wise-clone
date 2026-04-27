package repository

import (
	"database/sql"

	"github.com/gauravsahay007/split-wise-clone/models"
)

func (r *Repo) GetUserByProvider(provider, providerID string) (*models.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.profile_pic
		FROM users u
		JOIN auth_identities ai ON u.id = ai.user_id
		WHERE ai.provider = $1 AND ai.provider_id = $2
	`

	var user models.User
	err := r.DB.QueryRow(query, provider, providerID).Scan(&user.ID, &user.Name, &user.Email, &user.ProfilePic)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (r *Repo) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, name, email, profile_pic FROM users WHERE email=$1`

	var user models.User
	err := r.DB.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.ProfilePic)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (r *Repo) AddAuthIdentity(userID int, provider, providerID string) error {
	query := `
		INSERT INTO auth_identities(user_id, provider, provider_id)
		VALUES($1, $2, $3)
		ON CONFLICT (provider, provider_id) DO NOTHING
	`

	_, err := r.DB.Exec(query, userID, provider, providerID)
	return err
}
