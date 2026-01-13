package database

import (
	"context"
	"time"
)

func InsertUser(ctx context.Context, uid string) error {
	_, err := PG.ExecContext(ctx, `
	INSERT INTO user_linebot (uid, created_at)
	VALUES ($1, $2)
	ON CONFLICT (uid)
	DO UPDATE SET created_at = $2
	`,
		uid,
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}
