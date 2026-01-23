package database

import (
	"context"
)

func (s *SQLite) DeleteUser(ctx context.Context, uid string) error {
	if _, err := s.db.ExecContext(ctx, `
	UDPATE linebot_users
	SET dismiss = 1
	WHERE uid = ?
	`,
		uid,
	); err != nil {
		return err
	}

	return nil
}
