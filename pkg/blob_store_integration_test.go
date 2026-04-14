package pkg

import (
	"context"
	"io"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/fsouza/fake-gcs-server/fakestorage"
	"github.com/stretchr/testify/require"
)

var (
	testClient *storage.Client
	server     *fakestorage.Server
	testBucket string = "ptdf"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	server = fakestorage.NewServer(nil)
	defer server.Stop()

	testClient = server.Client()

	err := testClient.Bucket(testBucket).Create(ctx, "test-project", nil)
	if err != nil {
		panic("Could not create bucket: " + err.Error())
	}

	code := m.Run()
	if code != 0 {
		panic(" tests failed")
	}
}

func TestGcsCanWrite(t *testing.T) {
	t.Log("Initializing client")
	factory := GcsWriterFactory{Client: testClient}

	ctx := context.Background()
	object := "year=2020/file.bin"
	writer, err := factory.MakeWriteCloser(ctx, testBucket, object)
	require.NoError(t, err)

	t.Log("Data written")
	content := []byte("content")
	numBytes, err := writer.Write(content)
	require.NoError(t, err)
	require.Greater(t, numBytes, 0)
	err = writer.Close()
	require.NoError(t, err)

	bucketReader, err := testClient.Bucket(testBucket).Object(object).NewReader(ctx)
	require.NoError(t, err)

	result, err := io.ReadAll(bucketReader)
	require.NoError(t, err)
	require.Equal(t, content, result)
}
