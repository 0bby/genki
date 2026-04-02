// package handlers implements route handlers
package handlers

import (
	"net/http"
	"strconv"
	"time"
	_ "time/tzdata"

	"github.com/0bby/genki/engine/middleware"
	"github.com/0bby/genki/internal/cfg"
	"github.com/0bby/genki/internal/db"
	"github.com/0bby/genki/internal/logger"
)

const defaultLimitSize = 100
const maximumLimit = 500

func OptsFromRequest(r *http.Request, userID int32) db.GetItemsOpts {
	l := logger.FromContext(r.Context())

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = defaultLimitSize
	}
	if limit > maximumLimit {
		limit = defaultLimitSize
	}

	pageStr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	exerciseIdStr := r.URL.Query().Get("exercise_id")
	exerciseId, _ := strconv.Atoi(exerciseIdStr)

	tf := TimeframeFromRequest(r)

	l.Debug().Msgf("OptsFromRequest: limit=%d, page=%d, exercise_id=%d", limit, page, exerciseId)

	return db.GetItemsOpts{
		Limit:      limit,
		Page:       page,
		Timeframe:  tf,
		UserID:     userID,
		ExerciseID: exerciseId,
	}
}

func TimeframeFromRequest(r *http.Request) db.Timeframe {
	q := r.URL.Query()

	parseInt := func(key string) int {
		v := q.Get(key)
		if v == "" {
			return 0
		}
		i, _ := strconv.Atoi(v)
		return i
	}

	parseInt64 := func(key string) int64 {
		v := q.Get(key)
		if v == "" {
			return 0
		}
		i, _ := strconv.ParseInt(v, 10, 64)
		return i
	}

	return db.Timeframe{
		Period:   db.Period(q.Get("period")),
		Year:     parseInt("year"),
		Month:    parseInt("month"),
		Week:     parseInt("week"),
		FromUnix: parseInt64("from"),
		ToUnix:   parseInt64("to"),
		Timezone: parseTZ(r),
	}
}

func parseTZ(r *http.Request) *time.Location {
	overrides := map[string]string{
		"America/Indianapolis":  "America/Indiana/Indianapolis",
		"America/Montreal":      "America/Toronto",
		"US/Alaska":             "America/Anchorage",
		"US/Arizona":            "America/Phoenix",
		"US/Central":            "America/Chicago",
		"US/Eastern":            "America/New_York",
		"US/Hawaii":             "Pacific/Honolulu",
		"US/Mountain":           "America/Denver",
		"US/Pacific":            "America/Los_Angeles",
		"Canada/Atlantic":       "America/Halifax",
		"Canada/Central":        "America/Winnipeg",
		"Canada/Eastern":        "America/Toronto",
		"Canada/Mountain":       "America/Edmonton",
		"Canada/Pacific":        "America/Vancouver",
		"Asia/Calcutta":         "Asia/Kolkata",
		"Asia/Saigon":           "Asia/Ho_Chi_Minh",
		"Europe/Kiev":           "Europe/Kyiv",
		"Australia/ACT":         "Australia/Sydney",
		"Australia/Queensland":  "Australia/Brisbane",
		"NZ":                    "Pacific/Auckland",
		"UTC":                   "UTC",
		"Etc/UTC":               "UTC",
		"GMT":                   "UTC",
	}

	if cfg.ForceTZ() != nil {
		return cfg.ForceTZ()
	}

	if tz := r.URL.Query().Get("tz"); tz != "" {
		if fixedTz, exists := overrides[tz]; exists {
			tz = fixedTz
		}
		if loc, err := time.LoadLocation(tz); err == nil {
			return loc
		}
	}

	if c, err := r.Cookie("tz"); err == nil {
		var tz string
		if fixedTz, exists := overrides[c.Value]; exists {
			tz = fixedTz
		} else {
			tz = c.Value
		}
		if loc, err := time.LoadLocation(tz); err == nil {
			return loc
		}
	}

	return time.Now().Location()
}

func getUserFromContext(r *http.Request) int32 {
	user := middleware.GetUserFromContext(r.Context())
	if user != nil {
		return user.ID
	}
	return 0
}
