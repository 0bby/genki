package cfg

import (
	"fmt"
	"time"
)

func UserAgent() string {
	lock.RLock()
	defer lock.RUnlock()
	return globalConfig.userAgent
}

func ListenAddr() string {
	lock.RLock()
	defer lock.RUnlock()
	return fmt.Sprintf("%s:%d", globalConfig.bindAddr, globalConfig.listenPort)
}

func ConfigDir() string {
	lock.RLock()
	defer lock.RUnlock()
	return globalConfig.configDir
}

func DatabaseUrl() string {
	lock.RLock()
	defer lock.RUnlock()
	return globalConfig.databaseUrl
}

func LogLevel() int {
	lock.RLock()
	defer lock.RUnlock()
	return globalConfig.logLevel
}

func StructuredLogging() bool {
	lock.RLock()
	defer lock.RUnlock()
	return globalConfig.structuredLogging
}

func DefaultPassword() string {
	lock.RLock()
	defer lock.RUnlock()
	return globalConfig.defaultPw
}

func DefaultUsername() string {
	lock.RLock()
	defer lock.RUnlock()
	return globalConfig.defaultUsername
}

func DefaultTheme() string {
	lock.RLock()
	defer lock.RUnlock()
	return globalConfig.defaultTheme
}

func AllowedHosts() []string {
	lock.RLock()
	defer lock.RUnlock()
	return globalConfig.allowedHosts
}

func AllowAllHosts() bool {
	lock.RLock()
	defer lock.RUnlock()
	return globalConfig.allowAllHosts
}

func AllowedOrigins() []string {
	lock.RLock()
	defer lock.RUnlock()
	return globalConfig.allowedOrigins
}

func RateLimitDisabled() bool {
	lock.RLock()
	defer lock.RUnlock()
	return globalConfig.disableRateLimit
}

func LoginGate() bool {
	lock.RLock()
	defer lock.RUnlock()
	return globalConfig.loginGate
}

func ForceTZ() *time.Location {
	lock.RLock()
	defer lock.RUnlock()
	return globalConfig.forceTZ
}

// Fitness integration getters

func WgerURL() string {
	lock.RLock()
	defer lock.RUnlock()
	return globalConfig.wgerURL
}

func WgerToken() string {
	lock.RLock()
	defer lock.RUnlock()
	return globalConfig.wgerToken
}

func WgerEnabled() bool {
	lock.RLock()
	defer lock.RUnlock()
	return globalConfig.wgerURL != "" && globalConfig.wgerToken != ""
}

func FitbitClientID() string {
	lock.RLock()
	defer lock.RUnlock()
	return globalConfig.fitbitClientID
}

func FitbitClientSecret() string {
	lock.RLock()
	defer lock.RUnlock()
	return globalConfig.fitbitClientSecret
}

func FitbitRedirectURI() string {
	lock.RLock()
	defer lock.RUnlock()
	return globalConfig.fitbitRedirectURI
}

func FitbitEnabled() bool {
	lock.RLock()
	defer lock.RUnlock()
	return globalConfig.fitbitClientID != "" && globalConfig.fitbitClientSecret != ""
}

func KoitoURL() string {
	lock.RLock()
	defer lock.RUnlock()
	return globalConfig.koitoURL
}

func KoitoAPIKey() string {
	lock.RLock()
	defer lock.RUnlock()
	return globalConfig.koitoAPIKey
}

func WgerSyncInterval() time.Duration {
	lock.RLock()
	defer lock.RUnlock()
	return globalConfig.wgerSyncInterval
}

func FitbitSyncInterval() time.Duration {
	lock.RLock()
	defer lock.RUnlock()
	return globalConfig.fitbitSyncInterval
}
