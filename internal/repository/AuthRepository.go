package repository

import (
	"auth_service/internal/model"
	"context"
	"database/sql"
	"fmt"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/dbtx"
	"github.com/gogaruda/pkg/utils"
	"strings"
)

type AuthRepository interface {
	IsUsernameExists(username string) (bool, error)
	IsEmailExists(email string) (bool, error)
	CheckRoles(roles []string) ([]model.RoleModel, error)
	Create(ctx context.Context, user *model.UserModel) error
	IdentifierCheck(identifier string) (*model.UserModel, error)
	UpdateTokenVersion(userID string) (string, error)
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db}
}

func (r *authRepository) IsUsernameExists(username string) (bool, error) {
	var exists int
	err := r.db.QueryRow(`SELECT 1 FROM username_history WHERE username = ?`, username).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, apperror.New(apperror.CodeDBError, "gagal query username", err)
	}

	return true, nil
}

func (r *authRepository) IsEmailExists(email string) (bool, error) {
	var exists int
	err := r.db.QueryRow(`SELECT 1 FROM email_history WHERE email = ?`, email).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, apperror.New(apperror.CodeDBError, "gagal query email", err)
	}

	return true, nil
}

func (r *authRepository) CheckRoles(roles []string) ([]model.RoleModel, error) {
	if len(roles) == 0 {
		return nil, apperror.New(apperror.CodeBadRequest, "roles tidak boleh kosong", nil)
	}

	if len(roles) > 20 {
		return nil, apperror.New(apperror.CodeBadRequest, "jumlah role terlalu banyak", nil)
	}

	placeholders := make([]string, len(roles))
	args := make([]interface{}, len(roles))
	for i, role := range roles {
		placeholders[i] = "?"
		args[i] = role
	}

	query := fmt.Sprintf(
		`SELECT id, name FROM roles WHERE name IN (%s)`,
		strings.Join(placeholders, ","),
	)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, apperror.New(apperror.CodeDBError, "gagal query roles", err)
	}
	defer rows.Close()

	var foundRoles []model.RoleModel
	for rows.Next() {
		var role model.RoleModel
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			return nil, apperror.New(apperror.CodeDBError, "gagal membaca data role", err)
		}
		foundRoles = append(foundRoles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, apperror.New(apperror.CodeDBError, "error saat membaca hasil query", err)
	}

	if len(foundRoles) != len(roles) {
		return nil, apperror.New(apperror.CodeRoleNotFound, "salah satu atau lebih role tidak ditemukan", nil)
	}

	return foundRoles, nil
}

func (r *authRepository) Create(ctx context.Context, user *model.UserModel) error {
	return dbtx.WithTxContext(ctx, r.db, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO users(id, username, email, password) VALUES(?, ?, ?, ?)`,
			&user.ID, &user.Username, &user.Email, &user.Password)
		if err != nil {
			return apperror.New(apperror.CodeDBError, "query insert users gagal", err)
		}

		stmt, err := tx.PrepareContext(ctx, `INSERT INTO user_roles(user_id, role_id) VALUES(?, ?)`)
		if err != nil {
			return apperror.New(apperror.CodeDBPrepareError, "prepare insert user_roles gagal", err)
		}
		defer stmt.Close()

		for _, r := range user.Roles {
			_, err := stmt.ExecContext(ctx, user.ID, r.ID)
			if err != nil {
				return apperror.New(apperror.CodeDBError, "insert user_roles gagal", err)
			}
		}

		_, err = tx.ExecContext(ctx, `INSERT INTO username_history(username) VALUES(?)`, user.Username)
		if err != nil {
			return apperror.New(apperror.CodeDBError, "gagal insert username history", err)
		}

		_, err = tx.ExecContext(ctx, `INSERT INTO email_history(email) VALUES(?)`, user.Email)
		if err != nil {
			return apperror.New(apperror.CodeDBError, "gagal insert email history", err)
		}
		return nil
	})
}

func (r *authRepository) IdentifierCheck(identifier string) (*model.UserModel, error) {
	var user model.UserModel
	err := r.db.QueryRow(`SELECT id, password FROM users WHERE username = ? OR email = ?`,
		identifier, identifier).
		Scan(&user.ID, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.New(apperror.CodeUserNotFound, "username atau email tidak terdaftar..", err)
		}

		return nil, apperror.New(apperror.CodeDBError, "query check identifier gagal..", err)
	}

	rows, err := r.db.Query(`
													SELECT r.name FROM roles r
													INNER JOIN user_roles ur ON r.id = ur.role_id
													WHERE ur.user_id = ?
													`, user.ID)
	if err != nil {
		return nil, apperror.New(apperror.CodeDBError, "gagal query roles", err)
	}
	defer rows.Close()

	var roles []model.RoleModel
	for rows.Next() {
		var role model.RoleModel
		if err := rows.Scan(&role.Name); err != nil {
			return nil, apperror.New(apperror.CodeDBError, "scan roles gagal", err)
		}
		roles = append(roles, role)
	}
	user.Roles = roles

	if err := rows.Err(); err != nil {
		return nil, apperror.New(apperror.CodeDBError, "gagal setelah iterasi rows roles", err)
	}

	return &user, nil
}

func (r *authRepository) UpdateTokenVersion(userID string) (string, error) {
	newVersion := utils.NewULID()
	_, err := r.db.Exec(`UPDATE users SET token_version = ? WHERE id = ?`, newVersion, userID)
	if err != nil {
		return "", apperror.New(apperror.CodeDBError, "gagal update token_version", err)
	}

	return newVersion, nil
}
