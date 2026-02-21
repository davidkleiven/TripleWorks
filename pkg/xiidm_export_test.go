package pkg

import (
	"bytes"
	"context"
	"encoding/xml"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"com.github/davidkleiven/tripleworks/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func busBreakerData() []repository.BusBreakerConnection {
	mrid := uuid.New()
	return []repository.BusBreakerConnection{
		{
			Mrid:           mrid,
			Name:           "Sub1-Sub2",
			NominalVoltage: 22.0,
			SubstationMrid: uuid.New(),
		},
		{
			Mrid:           mrid,
			Name:           "Sub1-Sub2",
			NominalVoltage: 22.0,
			SubstationMrid: uuid.New(),
		},
		{
			Mrid:           uuid.New(),
			Name:           "Sub1-Sub3",
			NominalVoltage: 22.0,
			SubstationMrid: uuid.New(),
		},
	}
}

func TestDanglingLinesReported(t *testing.T) {
	data := busBreakerData()
	result := XiidmBusBreakerModel(data)
	require.Equal(t, 1, len(result.DanglingLines))
}

func downloadIfNotExist(t *testing.T, url, dir, filename string) string {
	err := os.MkdirAll(dir, 0o755)
	require.NoError(t, err)
	localPath := filepath.Join(dir, filename)
	if _, err := os.Stat(localPath); err == nil {
		t.Log("Using already downloaded file")
		return localPath
	}

	t.Log("Downloading file")
	resp, err := http.Get(url)
	require.NoError(t, err)
	require.NotNil(t, resp)
	defer resp.Body.Close()
	out, err := os.Create(localPath)
	require.NoError(t, err)
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	require.NoError(t, err)
	t.Log("File downloaded")
	return localPath
}

func TestXmlValid(t *testing.T) {
	schema := "https://github.com/powsybl/powsybl-core/raw/refs/heads/main/iidm/iidm-serde/src/main/resources/xsd/iidm_V1_16.xsd"
	_, isCi := os.LookupEnv("CI")
	_, err := exec.Command("xmllint", "--version").CombinedOutput()
	if err != nil && !isCi {
		t.Skip("Require 'xmllint'")
	}
	require.NoError(t, err)
	cacheDir := filepath.Join("testdata", "cache")
	fname := "iidm.xsd"
	path := downloadIfNotExist(t, schema, cacheDir, fname)

	data := busBreakerData()
	modelPath := filepath.Join(t.TempDir(), "model.xml")
	t.Logf("Model stored at %s", modelPath)
	result := XiidmBusBreakerModel(data)
	f, err := os.Create(modelPath)
	require.NoError(t, err)
	enc := xml.NewEncoder(f)
	err = enc.Encode(result.Network)
	f.Close()
	require.NoError(t, err)

	_, err = exec.Command("xmllint", "--noout", "--schema", path, modelPath).CombinedOutput()
	require.NoError(t, err)

}

func TestLogOnDanglingLines(t *testing.T) {
	defaultHandler := slog.Default()
	defer slog.SetDefault(defaultHandler)

	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, nil)
	slog.SetDefault(slog.New(handler))

	var result XiidmResult
	ctx := context.Background()
	t.Run("no log when empty", func(t *testing.T) {
		result.LogSummary(ctx)
		require.Equal(t, 0, len(buf.Bytes()))
	})

	result.DanglingLines = append(result.DanglingLines, uuid.New())
	t.Run("log when not empty", func(t *testing.T) {
		result.LogSummary(ctx)
		require.Greater(t, len(buf.Bytes()), 0)
	})
}
