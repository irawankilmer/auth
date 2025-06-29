package seeder

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gogaruda/dbtx"
	"github.com/gogaruda/pkg/utils"
	"time"
)

func Users(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return dbtx.WithTxContext(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
		var roleID string
		err := tx.QueryRowContext(ctx, `SELECT id FROM roles WHERE name = ?`, "super admin").Scan(&roleID)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("role tidak ditemukan: %v", err)
			} else {
				return fmt.Errorf("gagal query select role: %v", err)
			}
		}

		userID := utils.NewULID()
		hashPassword, _ := utils.GenerateHash("superadmin")
		_, err = tx.ExecContext(ctx, `INSERT INTO users(id, username, email, password) VALUES(?, ?, ?, ?)`,
			userID, "superadmin", "superadmin@gmail.com", hashPassword)
		if err != nil {
			return fmt.Errorf("gagal insert users: %v", err)
		}

		_, err = tx.ExecContext(ctx, `INSERT INTO username_history(username) VALUES(?)`, "superadmin")
		if err != nil {
			return fmt.Errorf("gagal insert username history: %v", err)
		}

		_, err = tx.ExecContext(ctx, `INSERT INTO email_history(email) VALUES(?)`, "superadmin@gmail.com")
		if err != nil {
			return fmt.Errorf("gagal insert email hostiry: %v", err)
		}

		_, err = tx.ExecContext(ctx, `INSERT INTO user_roles(user_id, role_id) VALUES(?, ?)`, userID, roleID)
		if err != nil {
			return fmt.Errorf("gagal insert user_roles: %v", err)
		}

		return nil
	})
}
