package pkg

import (
	"os"
	"path/filepath"
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

func TestLoadLocal(t *testing.T) {
	config := GetConfig("local_pg")
	require.Contains(t, config.DbUrl, "postgres")
}

func TestDefaultConfigOnNonExistentFile(t *testing.T) {
	config := NewDefaultConfig()
	loadedConfig := ConfigFromExternalFile("config.yaml")
	require.Equal(t, config, loadedConfig)
}

func TestDefaultConfigNotYamlFile(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "config.yaml")
	err := os.WriteFile(file, []byte("not yaml"), 0644)
	require.NoError(t, err)

	config := NewDefaultConfig()
	loadedConfig := ConfigFromExternalFile(file)
	require.Equal(t, config, loadedConfig)
}

func TestLoadConfigFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "config.yaml")
	err := os.WriteFile(file, []byte("dbUrl: my-database"), 0644)
	require.NoError(t, err)
	loadedConfig := GetConfig(file)
	require.Equal(t, "my-database", loadedConfig.DbUrl)
}

func TestPgEnv(t *testing.T) {
	defaultConfig := NewDefaultConfig()
	t.Run("default on non existent file", func(t *testing.T) {
		config := PgEnv(&FsOpener{})
		require.Equal(t, config, defaultConfig)
	})

	t.Run("default on failing read", func(t *testing.T) {
		config := PgEnv(&FailingReadOpener{})
		require.Equal(t, config, defaultConfig)
	})

	t.Run("load from env", func(t *testing.T) {
		os.Setenv("TRIPLEWORKS_DB_USER", "user")
		os.Setenv("TRIPLEWORKS_DB_HOST", "my-host")
		os.Setenv("TRIPLEWORKS_DB_PORT", "1234")
		os.Setenv("TRIPLEWORKS_DB_DATABASE", "mydatabase")

		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "pg_password")
		os.Setenv("TRIPLEWORKS_DB_PW_FILE", tmpFile)

		err := os.WriteFile(tmpFile, []byte("top-secret-password"), 0644)
		require.NoError(t, err)
		config := GetConfig("pg_env")
		require.Equal(t, "postgres://user:top-secret-password@my-host:1234/mydatabase", config.DbUrl)
	})
}
