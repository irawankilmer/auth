package seeder

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gogaruda/dbtx"
	"github.com/gogaruda/pkg/utils"
	"time"
)

func Profiles(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return dbtx.WithTxContext(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
		var userID string
		err := tx.QueryRowContext(ctx, `SELECT id FROM users WHERE username = ?`, "superadmin").Scan(&userID)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("username tidak ditemukan: %v", err)
			}

			return fmt.Errorf("query users gagal: %v", err)
		}

		_, err = tx.ExecContext(ctx, `INSERT INTO profiles(id, user_id, full_name, gender) VALUES(?, ?, ?, ?)`,
			utils.NewULID(), userID, "Saya Super Admin Pertama di Dunia", 1)
		if err != nil {
			return fmt.Errorf("insert profiles gagal: %v", err)
		}
    
		return nil
	})
}
