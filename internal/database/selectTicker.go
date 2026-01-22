package database

import (
	"context"
	"database/sql"
	"fmt"
)

type GexInput struct {
	ImageURL         string  `json:"image_url"`
	Spot             float32 `json:"spot"`
	Pivot            float32 `json:"pivot"`
	DirectionalBias  string  `json:"directional_bias"`
	ConfidencePct    string  `json:"confidence_pct"`
	VolatilityLevel  string  `json:"volatility_level"`
	RationaleCompact string  `json:"rationale_compact"`
	BullTrigger      string  `json:"bull_trigger"`
	BearTrigger      string  `json:"bear_trigger"`
	ExpiryNearTerm   string  `json:"expiry_near_term"`
	ExpiryFarTerm    string  `json:"expiry_far_term"`
}

func SelectTicker(ticker string) (GexInput, error) {
	var result GexInput

	if ticker == "" {
		return result, fmt.Errorf("ticker is empty")
	}

	ctx := context.Background()

	row := PG.QueryRowContext(ctx, `
	SELECT
		image_url,
		spot,
		pivot,
		directional_bias,
		confidence_pct,
		volatility_level,
		rationale_compact,
		bull_trigger,
		bear_trigger,
		expiry_near_term,
		expiry_far_term
	FROM gex_input
	WHERE ticker = $1
	ORDER BY id  DESC
	LIMIT 1
	`,
		ticker,
	)

	err := row.Scan(
		&result.ImageURL,
		&result.Spot,
		&result.Pivot,
		&result.DirectionalBias,
		&result.ConfidencePct,
		&result.VolatilityLevel,
		&result.RationaleCompact,
		&result.BullTrigger,
		&result.BearTrigger,
		&result.ExpiryNearTerm,
		&result.ExpiryFarTerm,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// 查無資料
			return result, fmt.Errorf("no data")
		}
		return result, err
	}
	return result, nil
}
