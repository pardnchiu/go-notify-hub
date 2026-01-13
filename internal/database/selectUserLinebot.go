package database

import (
	"context"
	"log/slog"
)

func SelectUserLinebot(ctx context.Context) ([]string, error) {
	const fn = "SelectUserLinebot"
	var results = []string{}

	rows, err := PG.QueryContext(ctx, `
	SELECT uid
	FROM user_linebot
	WHERE dismiss = FALSE
	`)
	if err != nil {
		// # SelectUserLinebot[0]
		slog.Error(fn+"[0]", "error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			// # SelectUserLinebot[1]
			slog.Warn(fn+"[1]", "error", err)
			continue
		}
		results = append(results, uid)
	}

	if err := rows.Err(); err != nil {
		// # SelectUserLinebot[2]
		slog.Error(fn+"[2]", "error", err)
		return nil, err
	}

	return results, nil
}
