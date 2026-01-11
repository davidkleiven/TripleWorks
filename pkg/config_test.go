package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun/dialect"
)

func TestConfig(t *testing.T) {
	config := NewDefaultConfig()
	assert.Equal(t, config.Port, 36000)
}

func TestNewTestConfigContainsInMemoryDatabase(t *testing.T) {
	config := NewTestConfig()
	require.Contains(t, config.DbUrl, "memory")
}

func TestDatabaseConnection(t *testing.T) {
	config := NewTestConfig()
	config.DbUrl = "postgres://"
	con := config.DatabaseConnection()
	require.Equal(t, con.Dialect().Name(), dialect.PG)
}

func TestDatabaseConnectionSqlite3(t *testing.T) {
	config := NewTestConfig()
	config.DbUrl = "whatever.db"
	con := config.DatabaseConnection()
	require.Equal(t, con.Dialect().Name(), dialect.SQLite)
}

func TestGetConfig(t *testing.T) {
	test := GetConfig("test")
	require.Contains(t, test.DbUrl, "memory")

	defaultConfig := GetConfig("")
	require.Equal(t, defaultConfig.DbUrl, "tripleworks.db")
}
