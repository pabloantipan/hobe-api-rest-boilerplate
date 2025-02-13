package cloud

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/logging"
	"github.com/pabloantipan/hobe-api-boilerplate/config"
	"google.golang.org/api/option"
)

var logName = "hobe-locations-api"

type CloudLogger struct {
	client *logging.Client
	logger *logging.Logger
	cfg    *config.Config
}

func NewCloudLogger(cfg *config.Config) (*CloudLogger, error) {
	ctx := context.Background()
	if cfg.LoggingServiceAccountPath == "" {
		log.Fatalf("CLOUD_LOGGING_CREDENTIALS_FILE environment variable not set")
	}

	client, err := logging.NewClient(ctx, cfg.ProjectID,
		option.WithCredentialsFile(cfg.LoggingServiceAccountPath),
	)
	if err != nil {
		return nil, err
	}

	logger := client.Logger(logName)

	return &CloudLogger{
		client: client,
		logger: logger,
		cfg:    cfg,
	}, nil
}

func (cl *CloudLogger) Log(payload interface{}) {
	cl.logger.Log(logging.Entry{
		Payload: map[string]interface{}{
			"who":     cl.cfg.Who,
			"type":    "info",
			"payload": payload,
		},
		Severity:  logging.Info,
		Timestamp: time.Now(),
	})
}

func (cl *CloudLogger) LogRequest(method, payload interface{}, path string, status int, latency float64) {
	cl.logger.Log(logging.Entry{
		Payload: map[string]interface{}{
			"who":     cl.cfg.Who,
			"method":  method,
			"payload": payload,
			"path":    path,
			"status":  status,
			"type":    "request",
			"latency": latency,
		},
		Severity: logging.Info,
	})
}

func (cl *CloudLogger) LogError(err error, payload interface{}, method, path string, latency float64) {
	cl.logger.Log(logging.Entry{
		Payload: map[string]interface{}{
			"who":     cl.cfg.Who,
			"error":   err.Error(),
			"payload": payload,
			"method":  method,
			"path":    path,
			"type":    "error",
			"latency": latency,
		},
		Severity: logging.Error,
	})
}
