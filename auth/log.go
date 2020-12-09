package auth

import (
	"context"

	log "github.com/sirupsen/logrus"
)

func lg(ctx context.Context) *log.Entry {
	// Placeholder for more expanded context logger
	return log.WithFields(log.Fields{})
}
