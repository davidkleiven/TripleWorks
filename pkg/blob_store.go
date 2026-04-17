package pkg

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"cloud.google.com/go/storage"
)

type WriterCloserFactory interface {
	MakeWriteCloser(ctx context.Context, bucket, object string) (io.WriteCloser, error)
}

type GcsWriterFactory struct {
	Client *storage.Client
}

func (g *GcsWriterFactory) MakeWriteCloser(ctx context.Context, bucket, object string) (io.WriteCloser, error) {
	w := g.Client.Bucket(bucket).Object(object).NewWriter(ctx)
	w.ContentType = "application/x-parquet"
	w.Metadata = map[string]string{
		"format": "parquet",
		"source": "tripleworks",
	}
	return w, nil
}

type LocalWriterFactory struct {
	Folder string
}

func (l *LocalWriterFactory) Filename(bucket, object string) string {
	return strings.Join(append([]string{bucket}, strings.Split(object, "/")...), "-")
}

func (l *LocalWriterFactory) MakeWriteCloser(ctx context.Context, bucket, object string) (io.WriteCloser, error) {
	filename := filepath.Join(l.Folder, l.Filename(bucket, object))
	f, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("Could not create file %s: %w", filename, err)
	}

	writer := bufio.NewWriter(f)
	return &FlushAndCloseWriter{
		Writer: writer,
		File:   f,
	}, nil
}

type FlushAndCloseWriter struct {
	*bufio.Writer
	File io.Closer
}

func (f *FlushAndCloseWriter) Close() error {
	errFlush := f.Flush()
	errClose := f.File.Close()
	return errors.Join(errFlush, errClose)
}

type MultiWriteCloser struct {
	WriteClosers []io.WriteCloser
}

func (m *MultiWriteCloser) Write(data []byte) (int, error) {
	var (
		n    int
		errs []error
	)
	for _, writer := range m.WriteClosers {
		newN, err := writer.Write(data)
		if newN > n {
			n = newN
		}
		errs = append(errs, err)
	}
	return n, errors.Join(errs...)
}

func (m *MultiWriteCloser) Close() error {
	var errs []error
	for _, closer := range m.WriteClosers {
		errs = append(errs, closer.Close())
	}
	return errors.Join(errs...)
}

type MultiWriterFactory struct {
	Factories []WriterCloserFactory
}

func (m *MultiWriterFactory) MakeWriteCloser(ctx context.Context, bucket, object string) (io.WriteCloser, error) {
	var (
		errs          []error
		writerClosers []io.WriteCloser
	)
	for _, factory := range m.Factories {
		w, err := factory.MakeWriteCloser(ctx, bucket, object)
		errs = append(errs, err)
		writerClosers = append(writerClosers, w)
	}
	return &MultiWriteCloser{WriteClosers: writerClosers}, errors.Join(errs...)
}

type InMemWriter struct {
	Name     string
	Data     []byte
	WriteErr error // Used for testing
	CloseErr error
}

func (i *InMemWriter) Write(data []byte) (int, error) {
	i.Data = append(i.Data, data...)
	return len(data), i.WriteErr
}

func (i *InMemWriter) Close() error {
	return i.CloseErr
}

type InMemWriterFactory struct {
	CreatedWriters []*InMemWriter
	Err            error
}

func (i *InMemWriterFactory) MakeWriteCloser(ctx context.Context, bucket, object string) (io.WriteCloser, error) {
	writer := InMemWriter{Name: filepath.Join(bucket, object)}
	i.CreatedWriters = append(i.CreatedWriters, &writer)
	return &writer, i.Err
}

type ReaderAtCloser interface {
	io.ReaderAt
	io.Closer
}

type LatestReadCloserFactory interface {
	// MakeReadCloser creates a reader to the latest item
	MakeReadCloser(ctx context.Context, bucket string) (ReaderAtCloser, error)
}

type LocalReaderFactory struct {
	Folder string
}

func (l *LocalReaderFactory) MakeReadCloser(ctx context.Context, bucket string) (ReaderAtCloser, error) {
	entries, err := os.ReadDir(l.Folder)
	if err != nil {
		return nil, fmt.Errorf("Could not read directory: %w", err)
	}
	var names []string
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.Contains(name, bucket) {
			continue
		}
		names = append(names, name)
	}
	if len(names) == 0 {
		return nil, fmt.Errorf("no files in directory: %s", bucket)
	}
	sort.Strings(names)
	filename := names[len(names)-1]
	slog.Info("Found latest local file", "filename", filename)
	return os.Open(filepath.Join(l.Folder, filename))
}
