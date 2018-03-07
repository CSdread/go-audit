package output

import (
	"errors"
	"fmt"
	"strings"

	"github.com/pantheon-systems/go-audit/pkg/slog"
	"github.com/spf13/viper"
)

type AuditWriterFactory func(conf *viper.Viper) (*AuditWriter, error)

var auditWriterFactories = make(map[string]AuditWriterFactory)

func Register(name string, factory AuditWriterFactory) {
	if factory == nil {
		slog.Error.Fatalf("Audit writer factory %s does not exist.", name)
	}
	_, registered := auditWriterFactories[name]
	if registered {
		slog.Info.Printf("Audit writer factory %s already registered. Ignoring.", name)
	}
	auditWriterFactories[name] = factory
}

func CreateAuditWriter(auditWriterName string, config *viper.Viper) (*AuditWriter, error) {
	auditWriterFactory, ok := auditWriterFactories[auditWriterName]
	if !ok {
		availableAuditWriters := GetAvailableAuditWriters()
		return nil, errors.New(fmt.Sprintf("Invalid audit writer name. Must be one of: %s", strings.Join(availableAuditWriters, ", ")))
	}

	// Run the factory with the configuration.
	return auditWriterFactory(config)
}

func GetAvailableAuditWriters() []string {
	availableAuditWriters := make([]string, len(auditWriterFactories))
	for k, _ := range auditWriterFactories {
		availableAuditWriters = append(availableAuditWriters, k)
	}
	return availableAuditWriters
}
