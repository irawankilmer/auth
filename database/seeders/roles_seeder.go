package seeder

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gogaruda/dbtx"
	"github.com/gogaruda/pkg/utils"
	"time"
)

func Roles(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return dbtx.WithTxContext(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
		stmt, err := tx.PrepareContext(ctx, `INSERT INTO roles(id, name) VALUES(?, ?)`)
		if err != nil {
			return fmt.Errorf("prepare insert roles gagal: %v", err)
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx, utils.NewULID(), "super admin")
		if err != nil {
			return fmt.Errorf("gagal insert role super admin: %v", err)
		}

		_, err = stmt.ExecContext(ctx, utils.NewULID(), "tamu")
		if err != nil {
			return fmt.Errorf("gagal insert role tamu: %v", err)
		}

		return nil
	})
}
