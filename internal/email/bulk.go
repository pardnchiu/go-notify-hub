package email

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type BulkRequest struct {
	To          any    `json:"to" binding:"required"`
	Subject     string `json:"subject" binding:"required"`
	Body        string `json:"body" binding:"required"`
	AltBody     string `json:"alt_body,omitempty"`
	From        string `json:"from,omitempty"`
	IsHTML      bool   `json:"is_html,omitempty"`
	MinDelay    int    `json:"min_delay,omitempty"`
	MaxDelay    int    `json:"max_delay,omitempty"`
	StopOnError bool   `json:"stop_on_error,omitempty"`
}

type BulkResult struct {
	Total   int               `json:"total"`
	Success int               `json:"success"`
	Failed  int               `json:"failed"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func (h *EmailHandler) SendBulk(c *gin.Context) {
	var req BulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest,
			fmt.Sprintf("failed to bind JSON: %v", err))
		return
	}

	target := parseTarget(req.To)
	if len(target) == 0 {
		c.String(http.StatusBadRequest,
			fmt.Sprint("no targets"))
		return
	}

	result := BulkResult{
		Total:  len(target),
		Errors: make(map[string]string),
	}

	delay := time.Duration(req.MinDelay) * time.Second
	if delay == 0 {
		delay = 1 * time.Second
	}

	for i, recipient := range target {
		email := &Email{
			To:      []Target{recipient},
			Subject: req.Subject,
			Body:    req.Body,
			AltBody: req.AltBody,
			IsHTML:  req.IsHTML,
		}

		if req.From != "" {
			from := getEmails(req.From)
			if len(from) > 0 {
				email.From = from[0]
			}
		}

		if err := h.mailer.Send(email); err != nil {
			result.Failed++
			result.Errors[recipient.Address] = err.Error()

			if req.StopOnError {
				for j := i + 1; j < len(target); j++ {
					result.Errors[target[j].Address] = "stopped due to previous error"
					result.Failed++
				}
				break
			}
		} else {
			result.Success++
		}

		if i < len(target)-1 {
			time.Sleep(delay)
		}
	}

	c.JSON(http.StatusOK, result)
}
