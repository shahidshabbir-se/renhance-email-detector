package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type betterStackHook struct {
	url   string
	token string
}

func (h *betterStackHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *betterStackHook) Fire(entry *logrus.Entry) error {
	payload := make(map[string]interface{})
	for k, v := range entry.Data {
		payload[k] = v
	}
	payload["level"] = entry.Level.String()
	payload["message"] = entry.Message
	payload["timestamp"] = entry.Time.Format(time.RFC3339)

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", h.url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("BetterStack logging failed: status %d", resp.StatusCode)
	}

	return nil
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func InitLogger(log *logrus.Logger) {
	env := strings.ToLower(getEnv("APP_ENV", "development"))

	if env == "production" {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			ForceColors:            true,
			FullTimestamp:          true,
			TimestampFormat:        "15:04:05",
			DisableLevelTruncation: false,
		})
	}

	log.SetOutput(os.Stdout)

	driver := strings.ToLower(getEnv("LOG_DRIVER", "console"))
	switch driver {
	case "betterstack":
		url := os.Getenv("BETTERSTACK_LOGS_URL")
		token := os.Getenv("BETTERSTACK_TOKEN")
		if url != "" && token != "" {
			log.AddHook(&betterStackHook{url: url, token: token})
			log.Info("📡 BetterStack logging enabled")
		} else {
			log.Warn("⚠️  BetterStack logging not fully configured")
		}
	case "console":
		log.Info("📜 Console logging enabled")
	default:
		log.SetOutput(nil)
		log.Info("🚫 Logging disabled")
	}

	if os.Getenv("HUNTER_API_KEY") == "" {
		log.Warn("⚠️  HUNTER_API_KEY not set — email detection may fail")
	}
}
