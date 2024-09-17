package stores

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

var (
	enc = binary.BigEndian
)

const (
	lenWidth = 8
)

type store struct {
	*os.File
	mu   sync.Mutex
	buf  *bufio.Writer
	size uint64
}

// newStore inicializa un nuevo store
func newStore(f *os.File) (*store, error) {
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	size := uint64(fi.Size())

	return &store{
		File: f,
		size: size,
		buf:  bufio.NewWriter(f),
	}, nil
}

// Append escribe bytes en el buffer y actualiza el tamaño del archivo
func (s *store) Append(p []byte) (n uint64, pos uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	pos = s.size
	if err := binary.Write(s.buf, enc, uint64(len(p))); err != nil {
		return 0, 0, err
	}

	writtenBytes, err := s.buf.Write(p)
	if err != nil {
		return 0, 0, err
	}

	writtenBytes += lenWidth
	s.size += uint64(writtenBytes)

	return uint64(writtenBytes), pos, nil
}

// Read lee bytes del archivo a partir de una posición dada
func (s *store) Read(pos uint64) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return nil, err
	}

	sizeBuf := make([]byte, lenWidth)
	if _, err := s.File.ReadAt(sizeBuf, int64(pos)); err != nil {
		return nil, err
	}

	size := enc.Uint64(sizeBuf)
	data := make([]byte, size)
	if _, err := s.File.ReadAt(data, int64(pos+lenWidth)); err != nil {
		return nil, err
	}

	return data, nil
}

// Close asegura que los datos del buffer se escriben y cierra el archivo
func (s *store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return err
	}
	return s.File.Close()
}

// ReadAt ayuda a leer bytes del archivo desde un offset dado
func (s *store) ReadAt(p []byte, off int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return 0, err
	}
	return s.File.ReadAt(p, off)
}
