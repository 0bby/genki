package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/0bby/genki/internal/db"
	"github.com/0bby/genki/internal/logger"
)

func GetActivityHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())

		metric := db.ActivityMetric(strings.ToLower(r.URL.Query().Get("metric")))
		if metric == "" {
			metric = db.MetricWorkouts
		}

		stepStr := r.URL.Query().Get("step")
		var step db.StepInterval
		switch stepStr {
		case "week":
			step = db.StepWeek
		case "month":
			step = db.StepMonth
		default:
			step = db.StepDay
		}

		rangeStr := r.URL.Query().Get("range")
		rangeVal, _ := strconv.Atoi(rangeStr)
		if rangeVal <= 0 {
			rangeVal = 182
		}

		monthStr := r.URL.Query().Get("month")
		month, _ := strconv.Atoi(monthStr)
		yearStr := r.URL.Query().Get("year")
		year, _ := strconv.Atoi(yearStr)

		userID := getUserFromContext(r)

		items, err := store.GetActivity(r.Context(), db.ActivityOpts{
			Metric:   metric,
			Step:     step,
			Range:    rangeVal,
			Month:    month,
			Year:     year,
			Timezone: parseTZ(r),
			UserID:   userID,
		})
		if err != nil {
			l.Error().Err(err).Msg("GetActivityHandler: failed to get activity")
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}
