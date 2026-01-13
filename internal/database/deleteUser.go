package database

import (
	"context"
)

func DeleteUser(ctx context.Context, uid string) error {
	_, err := PG.ExecContext(ctx, `
	UDPATE user_linebot
	SET dismiss = false
	WHERE uid = $1
	`,
		uid,
	)
	if err != nil {
		return err
	}

	return nil
}
