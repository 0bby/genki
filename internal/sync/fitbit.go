package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/0bby/genki/internal/cfg"
	"github.com/0bby/genki/internal/db"
	"github.com/0bby/genki/internal/models"
	"github.com/rs/zerolog"
)

const fitbitBaseURL = "https://api.fitbit.com"

// FitbitSync handles syncing data from the Fitbit Web API for a single user.
type FitbitSync struct {
	store  db.DB
	log    *zerolog.Logger
	client *http.Client
}

func NewFitbitSync(store db.DB, log *zerolog.Logger) *FitbitSync {
	return &FitbitSync{
		store:  store,
		log:    log,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// FitbitAuthURL returns the OAuth2 authorization URL for the user to visit.
func FitbitAuthURL(state string) string {
	params := url.Values{
		"response_type": {"code"},
		"client_id":     {cfg.FitbitClientID()},
		"redirect_uri":  {cfg.FitbitRedirectURI()},
		"scope":         {"activity heartrate sleep profile"},
		"state":         {state},
	}
	return "https://www.fitbit.com/oauth2/authorize?" + params.Encode()
}

// FitbitExchangeCode exchanges an authorization code for tokens.
func (f *FitbitSync) ExchangeCode(ctx context.Context, userID int32, code string) error {
	data := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"redirect_uri": {cfg.FitbitRedirectURI()},
		"client_id":    {cfg.FitbitClientID()},
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.fitbit.com/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(cfg.FitbitClientID(), cfg.FitbitClientSecret())

	resp, err := f.client.Do(req)
	if err != nil {
		return fmt.Errorf("fitbit token exchange: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("fitbit token exchange: status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		Scope        string `json:"scope"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("fitbit token decode: %w", err)
	}

	return f.store.UpsertOAuthToken(ctx, &models.OAuthToken{
		UserID:       userID,
		Provider:     "fitbit",
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
		Scopes:       &tokenResp.Scope,
	})
}

// refreshToken refreshes the Fitbit access token if expired.
func (f *FitbitSync) refreshToken(ctx context.Context, token *models.OAuthToken) (*models.OAuthToken, error) {
	if time.Now().Before(token.ExpiresAt.Add(-5 * time.Minute)) {
		return token, nil // not expired yet
	}

	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {token.RefreshToken},
		"client_id":     {cfg.FitbitClientID()},
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.fitbit.com/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(cfg.FitbitClientID(), cfg.FitbitClientSecret())

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fitbit token refresh: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("fitbit token refresh: status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	token.AccessToken = tokenResp.AccessToken
	token.RefreshToken = tokenResp.RefreshToken
	token.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	if err := f.store.UpsertOAuthToken(ctx, token); err != nil {
		return nil, err
	}
	return token, nil
}

func (f *FitbitSync) get(ctx context.Context, token *models.OAuthToken, path string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fitbitBaseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fitbit GET %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("fitbit rate limited")
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("fitbit GET %s: status %d: %s", path, resp.StatusCode, string(body))
	}
	return io.ReadAll(resp.Body)
}

// SyncUser syncs all Fitbit data for a given user (last 7 days).
func (f *FitbitSync) SyncUser(ctx context.Context, userID int32) error {
	token, err := f.store.GetOAuthToken(ctx, userID, "fitbit")
	if err != nil {
		return fmt.Errorf("fitbit: no token for user %d: %w", userID, err)
	}

	token, err = f.refreshToken(ctx, token)
	if err != nil {
		return fmt.Errorf("fitbit: token refresh failed: %w", err)
	}

	end := time.Now().Format("2006-01-02")
	start := time.Now().AddDate(0, 0, -7).Format("2006-01-02")

	if err := f.syncSteps(ctx, token, userID, start, end); err != nil {
		f.log.Error().Err(err).Msg("fitbit: syncSteps failed")
	}
	if err := f.syncActivity(ctx, token, userID, start, end); err != nil {
		f.log.Error().Err(err).Msg("fitbit: syncActivity failed")
	}
	if err := f.syncSleep(ctx, token, userID, start, end); err != nil {
		f.log.Error().Err(err).Msg("fitbit: syncSleep failed")
	}
	if err := f.syncHeartRate(ctx, token, userID, start, end); err != nil {
		f.log.Error().Err(err).Msg("fitbit: syncHeartRate failed")
	}

	now := time.Now()
	f.store.UpsertSyncCursor(ctx, &models.SyncCursor{
		UserID:       userID,
		Source:       "fitbit",
		Resource:     "all",
		LastSyncedAt: now,
	})

	f.log.Info().Int32("user_id", userID).Msg("fitbit: sync complete")
	return nil
}

func (f *FitbitSync) syncSteps(ctx context.Context, token *models.OAuthToken, userID int32, start, end string) error {
	path := fmt.Sprintf("/1/user/-/activities/steps/date/%s/%s.json", start, end)
	data, err := f.get(ctx, token, path)
	if err != nil {
		return err
	}

	var resp struct {
		Steps []struct {
			DateTime string `json:"dateTime"`
			Value    string `json:"value"`
		} `json:"activities-steps"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}

	for _, s := range resp.Steps {
		date, err := time.Parse("2006-01-02", s.DateTime)
		if err != nil {
			continue
		}
		var steps int32
		fmt.Sscanf(s.Value, "%d", &steps)

		f.store.UpsertDailySteps(ctx, &models.DailySteps{
			UserID:    userID,
			Date:      date,
			StepCount: steps,
			Source:    "fitbit",
		})
	}
	return nil
}

func (f *FitbitSync) syncActivity(ctx context.Context, token *models.OAuthToken, userID int32, start, end string) error {
	// Fetch multiple activity time series
	metrics := map[string]string{
		"minutesVeryActive":    fmt.Sprintf("/1/user/-/activities/minutesVeryActive/date/%s/%s.json", start, end),
		"minutesFairlyActive":  fmt.Sprintf("/1/user/-/activities/minutesFairlyActive/date/%s/%s.json", start, end),
		"minutesLightlyActive": fmt.Sprintf("/1/user/-/activities/minutesLightlyActive/date/%s/%s.json", start, end),
		"minutesSedentary":     fmt.Sprintf("/1/user/-/activities/minutesSedentary/date/%s/%s.json", start, end),
		"caloriesOut":          fmt.Sprintf("/1/user/-/activities/calories/date/%s/%s.json", start, end),
	}

	type dayData struct {
		DateTime string `json:"dateTime"`
		Value    string `json:"value"`
	}

	// Collect all metrics keyed by date
	byDate := map[string]models.DailyActivity{}

	for metric, path := range metrics {
		data, err := f.get(ctx, token, path)
		if err != nil {
			f.log.Warn().Err(err).Str("metric", metric).Msg("fitbit: activity metric fetch failed")
			continue
		}

		// Fitbit returns the key as "activities-<metric>"
		var raw map[string][]dayData
		if err := json.Unmarshal(data, &raw); err != nil {
			continue
		}

		// Find the first (only) key in the response
		for _, entries := range raw {
			for _, e := range entries {
				a := byDate[e.DateTime]
				a.Date, _ = time.Parse("2006-01-02", e.DateTime)
				a.UserID = userID
				a.Source = "fitbit"

				var val int32
				fmt.Sscanf(e.Value, "%d", &val)

				switch metric {
				case "minutesVeryActive":
					a.ActiveMinutes = val
				case "minutesFairlyActive":
					a.FairlyActiveMinutes = val
				case "minutesLightlyActive":
					a.LightlyActiveMinutes = val
				case "minutesSedentary":
					a.SedentaryMinutes = val
				case "caloriesOut":
					v := val
					a.CaloriesBurned = &v
				}
				byDate[e.DateTime] = a
			}
		}
	}

	for _, a := range byDate {
		f.store.UpsertDailyActivity(ctx, &a)
	}
	return nil
}

func (f *FitbitSync) syncSleep(ctx context.Context, token *models.OAuthToken, userID int32, start, end string) error {
	path := fmt.Sprintf("/1.2/user/-/sleep/date/%s/%s.json", start, end)
	data, err := f.get(ctx, token, path)
	if err != nil {
		return err
	}

	var resp struct {
		Sleep []struct {
			LogID         int64  `json:"logId"`
			DateOfSleep   string `json:"dateOfSleep"`
			Duration      int64  `json:"duration"` // milliseconds
			Efficiency    int    `json:"efficiency"`
			StartTime     string `json:"startTime"`
			EndTime       string `json:"endTime"`
			MinutesAsleep int    `json:"minutesAsleep"`
			MinutesAwake  int    `json:"minutesAwake"`
			Levels        struct {
				Summary map[string]struct {
					Minutes int `json:"minutes"`
				} `json:"summary"`
			} `json:"levels"`
		} `json:"sleep"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}

	for _, s := range resp.Sleep {
		date, err := time.Parse("2006-01-02", s.DateOfSleep)
		if err != nil {
			continue
		}

		totalMin := int32(s.Duration / 60000) // ms to min
		efficiency := int32(s.Efficiency)
		awakeMin := int32(s.MinutesAwake)
		sourceID := fmt.Sprintf("%d", s.LogID)

		var deepMin, lightMin, remMin *int32
		if v, ok := s.Levels.Summary["deep"]; ok {
			d := int32(v.Minutes)
			deepMin = &d
		}
		if v, ok := s.Levels.Summary["light"]; ok {
			l := int32(v.Minutes)
			lightMin = &l
		}
		if v, ok := s.Levels.Summary["rem"]; ok {
			r := int32(v.Minutes)
			remMin = &r
		}

		var startTime, endTime *time.Time
		if t, err := time.Parse("2006-01-02T15:04:05.000", s.StartTime); err == nil {
			startTime = &t
		}
		if t, err := time.Parse("2006-01-02T15:04:05.000", s.EndTime); err == nil {
			endTime = &t
		}

		f.store.UpsertSleepLog(ctx, &models.SleepLog{
			UserID:       userID,
			Date:         date,
			TotalMinutes: totalMin,
			DeepMinutes:  deepMin,
			LightMinutes: lightMin,
			REMMinutes:   remMin,
			AwakeMinutes: &awakeMin,
			Efficiency:   &efficiency,
			StartTime:    startTime,
			EndTime:      endTime,
			Source:       "fitbit",
			SourceID:     &sourceID,
		})
	}
	return nil
}

func (f *FitbitSync) syncHeartRate(ctx context.Context, token *models.OAuthToken, userID int32, start, end string) error {
	path := fmt.Sprintf("/1/user/-/activities/heart/date/%s/%s.json", start, end)
	data, err := f.get(ctx, token, path)
	if err != nil {
		return err
	}

	var resp struct {
		HeartActivities []struct {
			DateTime string `json:"dateTime"`
			Value    struct {
				RestingHeartRate int `json:"restingHeartRate"`
			} `json:"value"`
		} `json:"activities-heart"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}

	for _, h := range resp.HeartActivities {
		date, err := time.Parse("2006-01-02", h.DateTime)
		if err != nil {
			continue
		}

		var restingHR *int32
		if h.Value.RestingHeartRate > 0 {
			r := int32(h.Value.RestingHeartRate)
			restingHR = &r
		}

		f.store.UpsertHeartRateDaily(ctx, &models.HeartRateDaily{
			UserID:    userID,
			Date:      date,
			RestingHR: restingHR,
			Source:    "fitbit",
		})
	}
	return nil
}
