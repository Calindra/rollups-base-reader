package commons

import (
	"log/slog"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

const DefaultTimeout time.Duration = 5 * time.Minute
const DbImage = "postgres:17.4-alpine"
const DbName = "rollups"
const DbUser = "postgres"
const DbPassword = "password"
const OpenEpoch = 19

const Schema = "https://raw.githubusercontent.com/cartesi/rollups-graphql/c818a3ba3aa4bba7f263dc0e0f4d5899637385be/postgres/raw/rollupsdb-dump-202503282237.sql"

// StdoutLogConsumer is a LogConsumer that prints the log to stdout
type StdoutLogConsumer struct{}

// Accept prints the log to stdout
func (lc *StdoutLogConsumer) Accept(l testcontainers.Log) {
	if l.LogType == testcontainers.StderrLog {
		slog.Error(string(l.Content))
		return
	}

	slog.Debug(string(l.Content))
}
