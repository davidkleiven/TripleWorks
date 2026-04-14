package pkg

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInMemWriter(t *testing.T) {
	factory := InMemWriterFactory{}
	writer, err := factory.MakeWriteCloser(context.Background(), "bucket", "file.parquet")
	require.NoError(t, err)

	data := []byte("what?")
	writer.Write(data)
	require.Equal(t, 1, len(factory.CreatedWriters))
	require.Equal(t, factory.CreatedWriters[0].Data, data)
	require.NoError(t, writer.Close())
}

func TestLocalFileWriter(t *testing.T) {
	factory := LocalWriterFactory{}
	dir := t.TempDir()
	writer, err := factory.MakeWriteCloser(context.Background(), dir, "file.bin")
	require.NoError(t, err)
	_, err = writer.Write([]byte("content"))
	writer.Close()
	require.NoError(t, err)

	f, err := os.Open(filepath.Join(dir, "file.bin"))
	defer f.Close()
	require.NoError(t, err)
	content, err := io.ReadAll(f)
	require.NoError(t, err)
	require.Equal(t, []byte("content"), content)
}

func TestMultiWriter(t *testing.T) {
	factory1 := InMemWriterFactory{}
	factory2 := InMemWriterFactory{}
	factory := MultiWriterFactory{Factories: []WriterCloserFactory{&factory1, &factory2}}

	writer, err := factory.MakeWriteCloser(context.Background(), "bucket", "file.parquet")
	require.NoError(t, err)
	_, err = writer.Write([]byte("content"))
	require.NoError(t, err)
	require.Equal(t, 1, len(factory1.CreatedWriters))
	require.Equal(t, 1, len(factory2.CreatedWriters))
	require.NoError(t, writer.Close())
}
