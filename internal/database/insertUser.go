package database

import (
	"context"
	"time"
)

func (s *SQLite) InsertUser(ctx context.Context, uid string) error {

	if _, err := s.db.ExecContext(ctx, `
	INSERT INTO linebot_users (uid, created_at)
	VALUES (?, ?)
	ON CONFLICT (uid) DO UPDATE SET
	  created_at = ?
	`,
		uid,
		time.Now(),
		time.Now(),
	); err != nil {
		return err
	}

	return nil
}
