package Internal

import (
	"fmt"
	"io"
	"os"
	"path"

	api "github.com/danile0SA/0250952_SistemasDistribuidos/api/v1" // Importa tu paquete donde está definido Record
)

type segment struct {
	store                  *store
	index                  *index
	baseOffset, nextOffset uint64
	config                 Config
}

// newSegment crea un nuevo segmento y lo inicializa
func newSegment(dir string, baseOffset uint64, c Config) (*segment, error) {
	s := &segment{
		baseOffset: baseOffset,
		config:     c,
	}

	var err error
	// Crear y abrir el archivo del store
	storeFile, err := os.OpenFile(
		path.Join(dir, fmt.Sprintf("%d%s", baseOffset, ".store")),
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, err
	}
	if s.store, err = newStore(storeFile); err != nil {
		return nil, err
	}

	// Crear y abrir el archivo del index
	indexFile, err := os.OpenFile(
		path.Join(dir, fmt.Sprintf("%d%s", baseOffset, ".index")),
		os.O_RDWR|os.O_CREATE,
		0644,
	)
	if err != nil {
		return nil, err
	}
	if s.index, err = newIndex(indexFile, c); err != nil {
		return nil, err
	}

	// Ajustar el nextOffset
	if off, _, err := s.index.Read(-1); err != nil {
		s.nextOffset = baseOffset
	} else {
		s.nextOffset = baseOffset + uint64(off) + 1
	}

	return s, nil
}

// Append añade un registro a la tienda del segmento y actualiza el índice
func (s *segment) Append(record *api.Record) (uint64, error) {
	// Verificar si el store está lleno
	if s.IsMaxed() {
		return 0, io.EOF
	}

	// Serializar el record usando protobuf y agregarlo al store
	// Ahora store.Append retorna (n uint64, pos uint64, err error)
	writtenBytes, pos, err := s.store.Append(record.Value)
	if err != nil {
		return 0, err
	}
	fmt.Printf("writtenBytes: %v\n", writtenBytes)
	// Escribir en el índice: el offset es relativo al baseOffset en el archivo de store
	if err = s.index.Write(
		uint64(s.nextOffset-uint64(s.baseOffset)), // Offset para el índice
		pos, // Posición en el store donde empieza el registro
	); err != nil {
		return 0, err
	}

	// Guardar el offset y actualizar el siguiente offset
	offset := s.nextOffset
	s.nextOffset++

	return offset, nil
}

// Read lee un registro del segmento en base a un offset
func (s *segment) Read(off uint64) (*api.Record, error) {
	// Leer la posición desde el índice
	_, pos, err := s.index.Read(int64(off - s.baseOffset))
	if err != nil {
		return nil, err
	}

	// Leer los datos del store
	data, err := s.store.Read(pos)
	if err != nil {
		return nil, err
	}

	// Crear un nuevo Record y deserializar los datos
	record := &api.Record{}
	record.Value = data
	return record, nil
}

// IsMaxed verifica si el segmento ha alcanzado su tamaño máximo en el store o el índice
func (s *segment) IsMaxed() bool {
	return s.store.size >= s.config.Segment.MaxStoreBytes || s.index.size >= s.config.Segment.MaxIndexBytes
}

// Remove elimina el store y el índice del sistema de archivos
func (s *segment) Remove() error {
	if err := s.Close(); err != nil {
		return err
	}
	if err := os.Remove(s.index.file.Name()); err != nil {
		return err
	}
	if err := os.Remove(s.store.File.Name()); err != nil {
		return err
	}
	return nil
}

// Close cierra los archivos de índice y store
func (s *segment) Close() error {
	if err := s.index.Close(); err != nil {
		return err
	}
	if err := s.store.Close(); err != nil {
		return err
	}
	return nil
}
