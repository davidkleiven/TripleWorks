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
	dir := t.TempDir()
	factory := LocalWriterFactory{Folder: dir}
	writer, err := factory.MakeWriteCloser(context.Background(), "ptdfs", "date=now/file.bin")
	require.NoError(t, err)
	_, err = writer.Write([]byte("content"))
	writer.Close()
	require.NoError(t, err)

	f, err := os.Open(filepath.Join(dir, "ptdfs-date=now-file.bin"))
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

func TestLocalReader(t *testing.T) {
	dir := t.TempDir()
	readerFactory := LocalReaderFactory{Folder: dir}
	ctx := context.Background()
	_, err := readerFactory.MakeReadCloser(ctx, "my-bucket")
	require.ErrorContains(t, err, "Could not read")

	err = os.Mkdir(filepath.Join(dir, "my-bucket"), 0755)
	require.NoError(t, err)

	_, err = readerFactory.MakeReadCloser(ctx, "my-bucket")
	require.ErrorContains(t, err, "no files")

	err = os.WriteFile(filepath.Join(dir, "my-bucket", "data.bin"), []byte("content"), 0755)
	require.NoError(t, err)

	// Also create a directory that appears to be later
	err = os.Mkdir(filepath.Join(dir, "my-bucket", "edata"), 0755)
	require.NoError(t, err)

	reader, err := readerFactory.MakeReadCloser(ctx, "my-bucket")
	require.NoError(t, err)

	file, ok := reader.(*os.File)
	require.True(t, ok)

	content, err := io.ReadAll(file)
	require.NoError(t, err)
	require.Equal(t, []byte("content"), content)
}

func TestLocalFilename(t *testing.T) {
	name := "year=2024/month=04/file.bin"
	wf := LocalWriterFactory{}
	result := wf.Filename("my-bucket", name)
	want := "my-bucket-year=2024-month=04-file.bin"
	require.Equal(t, want, result)
}
