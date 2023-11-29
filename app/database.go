package app

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	DriverName = "mysql"
)

func (app *application) openDB() error {
	db, err := sql.Open(DriverName, app.config.DSN())
	if err != nil {
		return fmt.Errorf("database open failed; %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use PingContext() to establish a new connection to the database, passing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within the 5 second deadline, then this will return an
	// error.
	err = db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("db ping failed; %w", err)
	}

	app.db = db

	return nil
}
