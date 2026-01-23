package email

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type SendRequest struct {
	To       any    `json:"to" binding:"required"`
	Subject  string `json:"subject" binding:"required"`
	Body     string `json:"body" binding:"required"`
	AltBody  string `json:"alt_body,omitempty"`
	From     string `json:"from,omitempty"`
	CC       any    `json:"cc,omitempty"`
	BCC      any    `json:"bcc,omitempty"`
	Priority string `json:"priority,omitempty"`
	IsHTML   bool   `json:"is_html,omitempty"`
}

func (h *EmailHandler) Send(c *gin.Context) {
	var req SendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest,
			fmt.Sprintf("failed to bind JSON: %v", err))
		return
	}

	email := &Email{
		To:      parseTarget(req.To),
		CC:      parseTarget(req.CC),
		BCC:     parseTarget(req.BCC),
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

	switch req.Priority {
	case "high":
		email.Priority = PriorityHigh
	case "low":
		email.Priority = PriorityLow
	default:
		email.Priority = PriorityNormal
	}

	if err := h.mailer.Send(email); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("failed to send email: %v", err))
		return
	}

	c.String(http.StatusOK, "ok")
}

func getEmails(s string) []Target {
	var results []Target

	parts := strings.SplitSeq(s, ",")
	for e := range parts {
		e = strings.TrimSpace(e)
		if e == "" {
			continue
		}

		pair := strings.SplitN(e, ":", 2)
		if len(pair) == 2 {
			addr := strings.TrimSpace(pair[1])
			if emailRegex.MatchString(addr) {
				results = append(results, Target{
					Name:    strings.TrimSpace(pair[0]),
					Address: addr,
				})
			}
			continue
		}

		email := strings.TrimSpace(pair[0])
		if emailRegex.MatchString(email) {
			results = append(results, Target{
				Address: email,
			})
		}
	}

	return results
}

func parseTarget(v any) []Target {
	if v == nil {
		return nil
	}

	switch val := v.(type) {
	case string:
		if val == "" {
			return nil
		}
		return getEmails(val)

	case []string:
		var recs []Target
		for _, addr := range val {
			if addr != "" {
				recs = append(recs, getEmails(addr)...)
			}
		}
		return recs

	case []any:
		var recs []Target
		for _, item := range val {
			if str, ok := item.(string); ok && str != "" {
				recs = append(recs, getEmails(str)...)
			}
		}
		return recs

	case map[string]any:
		var recs []Target
		for addr, name := range val {
			if emailRegex.MatchString(addr) {
				nameStr := ""
				if n, ok := name.(string); ok {
					nameStr = n
				}
				recs = append(recs, Target{
					Address: addr,
					Name:    nameStr,
				})
			}
		}
		return recs

	default:
		return nil
	}
}
