package seeder

import "database/sql"

func SeedRun(db *sql.DB) error {
	if err := Roles(db); err != nil {
		return err
	}

	if err := Users(db); err != nil {
		return err
	}

	if err := Profiles(db); err != nil {
		return err
	}

	return nil
}
