package repository

import (
	"auth_service/internal/dto/request"
	"auth_service/internal/dto/response"
	"context"
	"database/sql"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/dbtx"
	"github.com/gogaruda/pkg/utils"
)

type ProfileRepository interface {
	IsUserExists(userID string) (bool, error)
	Create(profile request.ProfileRequest) error
	Update(profile request.ProfileRequest) error
	Me(ctx context.Context, userID string) (*response.ProfileResponse, error)
}

type profileRepository struct {
	db *sql.DB
}

func NewProfileRepository(db *sql.DB) ProfileRepository {
	return &profileRepository{db}
}

func (r *profileRepository) IsUserExists(userID string) (bool, error) {
	var id string
	err := r.db.QueryRow(`SELECT id FROM profiles WHERE user_id = ?`, userID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, apperror.New(apperror.CodeDBError, "query profil gagal", err)
	}

	return true, nil
}

func (r *profileRepository) Create(profile request.ProfileRequest) error {
	_, err := r.db.Exec(`INSERT INTO profiles(id, user_id, full_name, address, gender) VALUES(?, ?, ?, ?, ?)`,
		utils.NewULID(), profile.UserID, profile.FullName, profile.Address, profile.Gender)
	if err != nil {
		return apperror.New(apperror.CodeDBError, "query insert profiles gagal", err)
	}

	return nil
}

func (r *profileRepository) Update(profile request.ProfileRequest) error {
	_, err := r.db.Exec(`UPDATE profiles SET full_name = ?, address = ?, gender = ? WHERE user_id = ?`,
		profile.FullName, profile.Address, profile.Gender, profile.UserID)
	if err != nil {
		return apperror.New(apperror.CodeDBError, "query update profiles gagal", err)
	}

	return nil
}

func (r *profileRepository) Me(ctx context.Context, userID string) (*response.ProfileResponse, error) {
	var profile response.ProfileResponse
	var address, gender sql.NullString
	err := dbtx.WithTxContext(ctx, r.db, func(ctx context.Context, tx *sql.Tx) error {
		err := tx.QueryRowContext(ctx, `SELECT id, user_id, full_name, address, gender FROM profiles WHERE user_id = ?`, userID).
			Scan(&profile.ID, &profile.UserID, &profile.FullName, &address, &gender)
		if err != nil {
			if err == sql.ErrNoRows {
				return apperror.New("PROFILE_NOT_FOUND", "profile tidak ditemukan", err, 404)
			}

			return apperror.New(apperror.CodeDBError, "query profile gagal", err)
		}

		if address.Valid {
			profile.Address = &address.String
		} else {
			profile.Address = nil
		}

		if gender.Valid {
			profile.Gender = &gender.String
		} else {
			profile.Gender = nil
		}

		var user response.UserResponse
		err = tx.QueryRowContext(ctx, `SELECT id, username, email FROM users WHERE id = ?`, profile.UserID).
			Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				return apperror.New(apperror.CodeUserNotFound, "user tidak ditemukan", err)
			}

			return apperror.New(apperror.CodeDBError, "query users gagal", err)
		}

		profile.User = user

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &profile, nil
}
