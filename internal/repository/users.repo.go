package repository

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/oktaviandwip/musalabel-backend/config"
	models "github.com/oktaviandwip/musalabel-backend/internal/models"
)

type RepoUsersIF interface {
	CreateUser(data *models.User) (*config.Result, error)
	GetPassByEmail(email string) (*models.User, error)
	UpdateProfile(user *models.User) (*config.Result, error)
	UpdatePassword(user *models.User) (*config.Result, error)
	UpdateCheckoutUser(user *models.User) (*config.Result, error)
}

type RepoUsers struct {
	*sqlx.DB
}

func NewUser(db *sqlx.DB) *RepoUsers {
	return &RepoUsers{db}
}

// Create User
func (r *RepoUsers) CreateUser(data *models.User) (*config.Result, error) {
	q := `INSERT INTO users (
					email, 
					password, 
					phone_number,
					role,
					image
				)
				VALUES(
					:email,
					:password,
					:phone_number,
					:role,
					:image
				)`

	_, err := r.NamedExec(q, data)
	if err != nil {
		if strings.Contains(err.Error(), "users_email_key") {
			return nil, errors.New("email sudah digunakan")
		}
		return nil, err
	}

	return &config.Result{Message: "1 data user created"}, nil
}

// Login
func (r *RepoUsers) GetPassByEmail(email string) (*models.User, error) {
	var result models.User

	q := `SELECT * FROM users WHERE email = $1`

	if err := r.Get(&result, q, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("email belum terdaftar")
		}
		return nil, err
	}

	return &result, nil
}

// Update Profile
func (r *RepoUsers) UpdateProfile(data *models.User) (*config.Result, error) {
	q := `
		UPDATE users
		SET
			image = $1,
			email = COALESCE(NULLIF($2, ''), email),
			phone_number = COALESCE(NULLIF($3, ''), phone_number),
			address1 = COALESCE(NULLIF($4, ''), address1),
			address2 = COALESCE(NULLIF($5, ''), address2),
			address3 = COALESCE(NULLIF($6, ''), address3),
			full_name = COALESCE(NULLIF($7, ''), full_name),
			username = COALESCE(NULLIF($8, ''), username),
			birthday = COALESCE(CAST(NULLIF($9, '') AS DATE), birthday),
			gender = COALESCE(NULLIF($10, ''), gender),
			updated_at = NOW()
		WHERE
			id = CAST($11 AS UUID)
		RETURNING id, image, email, phone_number, role, address1, address2, address3, full_name, username, birthday, gender;
	`

	args := []interface{}{
		data.Image,
		data.Email,
		data.Phone_number,
		data.Address1,
		data.Address2,
		data.Address3,
		data.Full_name,
		data.Username,
		data.Birthday,
		data.Gender,
		data.Id,
	}

	var user models.User
	err := r.QueryRowx(q, args...).StructScan(&user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user tidak ditemukan")
		}
		return nil, err
	}

	return &config.Result{Data: user}, nil
}

// Update Password
func (r *RepoUsers) UpdatePassword(data *models.User) (*config.Result, error) {
	q := `
		UPDATE users
		SET
			password = COALESCE(NULLIF(:password, ''), password),
			updated_at = NOW()
		WHERE
			email = :email
	`

	_, err := r.NamedExec(q, data)
	if err != nil {
		return nil, err
	}

	return &config.Result{Message: "Password berhasil di-update!"}, nil
}

// Update Address
func (r *RepoUsers) UpdateCheckoutUser(data *models.User) (*config.Result, error) {
	q := `
		UPDATE users
		SET
			full_name = COALESCE(NULLIF(:full_name, ''), full_name),
			phone_number = COALESCE(NULLIF(:phone_number, ''), phone_number),
			address1 = COALESCE(NULLIF(:address1, ''), address1),
			address2 = COALESCE(NULLIF(:address2, ''), address2),
			address3 = COALESCE(NULLIF(:address3, ''), address3),
			updated_at = NOW()
		WHERE
			id = CAST(:id AS UUID)
	`

	_, err := r.NamedExec(q, data)
	if err != nil {
		return nil, err
	}

	return &config.Result{Message: "Alamat berhasil di-update!"}, nil
}
