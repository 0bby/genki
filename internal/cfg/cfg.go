package cfg

import (
	"errors"
	"fmt"
	"strings"
	"strconv"
	"sync"
	"time"
)

const (
	defaultListenPort = 4110
)

const (
	DATABASE_URL_ENV              = "GENKI_DATABASE_URL"
	BIND_ADDR_ENV                 = "GENKI_BIND_ADDR"
	LISTEN_PORT_ENV               = "GENKI_LISTEN_PORT"
	ENABLE_STRUCTURED_LOGGING_ENV = "GENKI_ENABLE_STRUCTURED_LOGGING"
	LOG_LEVEL_ENV                 = "GENKI_LOG_LEVEL"
	CONFIG_DIR_ENV                = "GENKI_CONFIG_DIR"
	DEFAULT_USERNAME_ENV          = "GENKI_DEFAULT_USERNAME"
	DEFAULT_PASSWORD_ENV          = "GENKI_DEFAULT_PASSWORD"
	DEFAULT_THEME_ENV             = "GENKI_DEFAULT_THEME"
	ALLOWED_HOSTS_ENV             = "GENKI_ALLOWED_HOSTS"
	CORS_ORIGINS_ENV              = "GENKI_CORS_ALLOWED_ORIGINS"
	DISABLE_RATE_LIMIT_ENV        = "GENKI_DISABLE_RATE_LIMIT"
	LOGIN_GATE_ENV                = "GENKI_LOGIN_GATE"
	FORCE_TZ                      = "GENKI_FORCE_TZ"

	// Fitness-specific config
	WGER_URL_ENV              = "GENKI_WGER_URL"
	WGER_TOKEN_ENV            = "GENKI_WGER_TOKEN"
	FITBIT_CLIENT_ID_ENV      = "GENKI_FITBIT_CLIENT_ID"
	FITBIT_CLIENT_SECRET_ENV  = "GENKI_FITBIT_CLIENT_SECRET"
	FITBIT_REDIRECT_URI_ENV   = "GENKI_FITBIT_REDIRECT_URI"
	KOITO_URL_ENV             = "GENKI_KOITO_URL"
	KOITO_API_KEY_ENV         = "GENKI_KOITO_API_KEY"
	WGER_SYNC_INTERVAL_ENV    = "GENKI_WGER_SYNC_INTERVAL"
	FITBIT_SYNC_INTERVAL_ENV  = "GENKI_FITBIT_SYNC_INTERVAL"
)

type config struct {
	bindAddr          string
	listenPort        int
	configDir         string
	databaseUrl       string
	logLevel          int
	structuredLogging bool
	defaultPw         string
	defaultUsername   string
	defaultTheme      string
	allowedHosts      []string
	allowAllHosts     bool
	allowedOrigins    []string
	disableRateLimit  bool
	userAgent         string
	loginGate         bool
	forceTZ           *time.Location

	// Fitness integrations
	wgerURL            string
	wgerToken          string
	fitbitClientID     string
	fitbitClientSecret string
	fitbitRedirectURI  string
	koitoURL           string
	koitoAPIKey        string
	wgerSyncInterval   time.Duration
	fitbitSyncInterval time.Duration
}

var (
	globalConfig *config
	once         sync.Once
	lock         sync.RWMutex
)

func Load(getenv func(string) string, version string) error {
	var err error
	once.Do(func() {
		globalConfig, err = loadConfig(getenv, version)
	})
	if err != nil {
		return fmt.Errorf("cfg.Load: %w", err)
	}
	return nil
}

func loadConfig(getenv func(string) string, version string) (*config, error) {
	cfg := new(config)

	cfg.databaseUrl = getenv(DATABASE_URL_ENV)
	if cfg.databaseUrl == "" {
		return nil, errors.New("loadConfig: required parameter " + DATABASE_URL_ENV + " not provided")
	}
	cfg.bindAddr = getenv(BIND_ADDR_ENV)
	var err error
	cfg.listenPort, err = strconv.Atoi(getenv(LISTEN_PORT_ENV))
	if err != nil {
		cfg.listenPort = defaultListenPort
	}

	cfg.disableRateLimit = parseBool(getenv(DISABLE_RATE_LIMIT_ENV))
	cfg.structuredLogging = parseBool(getenv(ENABLE_STRUCTURED_LOGGING_ENV))

	cfg.userAgent = fmt.Sprintf("Genki %s", version)

	if getenv(DEFAULT_USERNAME_ENV) == "" {
		cfg.defaultUsername = "admin"
	} else {
		cfg.defaultUsername = getenv(DEFAULT_USERNAME_ENV)
	}
	if getenv(DEFAULT_PASSWORD_ENV) == "" {
		cfg.defaultPw = "changeme"
	} else {
		cfg.defaultPw = getenv(DEFAULT_PASSWORD_ENV)
	}

	cfg.defaultTheme = getenv(DEFAULT_THEME_ENV)

	cfg.configDir = getenv(CONFIG_DIR_ENV)
	if cfg.configDir == "" {
		cfg.configDir = "/etc/genki"
	}

	rawHosts := getenv(ALLOWED_HOSTS_ENV)
	cfg.allowedHosts = strings.Split(rawHosts, ",")
	cfg.allowAllHosts = cfg.allowedHosts[0] == "*"

	rawCors := getenv(CORS_ORIGINS_ENV)
	cfg.allowedOrigins = strings.Split(rawCors, ",")

	if strings.ToLower(getenv(LOGIN_GATE_ENV)) == "true" {
		cfg.loginGate = true
	}

	if getenv(FORCE_TZ) != "" {
		cfg.forceTZ, err = time.LoadLocation(getenv(FORCE_TZ))
		if err != nil {
			return nil, fmt.Errorf("forced timezone '%s' is not a valid timezone", getenv(FORCE_TZ))
		}
	}

	switch strings.ToLower(getenv(LOG_LEVEL_ENV)) {
	case "debug":
		cfg.logLevel = 0
	case "warn":
		cfg.logLevel = 2
	case "error":
		cfg.logLevel = 3
	case "fatal":
		cfg.logLevel = 4
	default:
		cfg.logLevel = 1
	}

	// Fitness integrations
	cfg.wgerURL = getenv(WGER_URL_ENV)
	cfg.wgerToken = getenv(WGER_TOKEN_ENV)
	cfg.fitbitClientID = getenv(FITBIT_CLIENT_ID_ENV)
	cfg.fitbitClientSecret = getenv(FITBIT_CLIENT_SECRET_ENV)
	cfg.fitbitRedirectURI = getenv(FITBIT_REDIRECT_URI_ENV)
	cfg.koitoURL = getenv(KOITO_URL_ENV)
	cfg.koitoAPIKey = getenv(KOITO_API_KEY_ENV)

	cfg.wgerSyncInterval, _ = time.ParseDuration(getenv(WGER_SYNC_INTERVAL_ENV))
	if cfg.wgerSyncInterval == 0 {
		cfg.wgerSyncInterval = 5 * time.Minute
	}
	cfg.fitbitSyncInterval, _ = time.ParseDuration(getenv(FITBIT_SYNC_INTERVAL_ENV))
	if cfg.fitbitSyncInterval == 0 {
		cfg.fitbitSyncInterval = 15 * time.Minute
	}

	return cfg, nil
}

func parseBool(s string) bool {
	return strings.ToLower(s) == "true"
}
