package cloud

import (
	"context"
	"log"

	"cloud.google.com/go/logging"
	"github.com/pabloantipan/hobe-api-boilerplate/config"
	"google.golang.org/api/option"
)

var logName = "hobe-api-boilerplate"

type CloudLogger struct {
	client *logging.Client
	logger *logging.Logger
}

func NewCloudLogger(cfg *config.Config) (*CloudLogger, error) {
	ctx := context.Background()
	if cfg.CloudLoggingCredentialsFile == "" {
		log.Fatalf("CLOUD_LOGGING_CREDENTIALS_FILE environment variable not set")
	}

	client, err := logging.NewClient(ctx, cfg.ProjectID,
		option.WithCredentialsFile(cfg.CloudLoggingCredentialsFile),
	)
	if err != nil {
		return nil, err
	}

	logger := client.Logger(logName)

	return &CloudLogger{
		client: client,
		logger: logger,
	}, nil
}

func (cl *CloudLogger) LogRequest(method, path string, status int, latency float64) {
	cl.logger.Log(logging.Entry{
		Payload: map[string]interface{}{
			"method":  method,
			"path":    path,
			"status":  status,
			"type":    "request",
			"latency": latency,
		},
		Severity: logging.Info,
	})
}

func (cl *CloudLogger) LogError(err error, method, path string, latency float64) {
	cl.logger.Log(logging.Entry{
		Payload: map[string]interface{}{
			"error":   err.Error(),
			"method":  method,
			"path":    path,
			"type":    "error",
			"latency": latency,
		},
		Severity: logging.Error,
	})
}
