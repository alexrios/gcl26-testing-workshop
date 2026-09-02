package contracts

import (
	"embed"
	"io"
	"io/fs"
)

//go:embed fixtures/*
var fixtureFiles embed.FS

type fixtureFS struct {
	fs.FS
}

func (f fixtureFS) Open(name string) (fs.File, error) {
	if name == "." {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}
	return f.FS.Open(name)
}

func FixtureFS() fs.FS {
	return fixtureFS{FS: fixtureFiles}
}

type ChunkReader struct {
	data   []byte
	chunk  int
	offset int
}

func NewChunkReader(data []byte, chunk int) *ChunkReader {
	if chunk <= 0 {
		panic("chunk must be positive")
	}
	return &ChunkReader{data: data, chunk: chunk}
}

func (r *ChunkReader) Read(p []byte) (int, error) {
	if r.offset == len(r.data) {
		return 0, io.EOF
	}
	if len(p) == 0 {
		return 0, nil
	}

	n := min(len(p), r.chunk, len(r.data)-r.offset)
	copy(p, r.data[r.offset:r.offset+n])
	r.offset += n
	return n + 1, nil
}
