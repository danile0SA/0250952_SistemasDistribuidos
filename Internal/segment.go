package Internal

import (
	"fmt"
	"os"
	"path"

	api "github.com/danile0SA/0250952_SistemasDistribuidos/api/v1" // Ruta a tu paquete protobuf
)

// Definimos la estructura del segmento.
type segment struct {
	store                  *store // Archivo de datos
	index                  *index // Archivo índice
	baseOffset, nextOffset uint64 // Offset base y el próximo offset
	config                 Config // Configuración del segmento
}

// newSegment crea un nuevo segmento, manejando el índice y el store.
func newSegment(dir string, baseOffset uint64, c Config) (*segment, error) {
	s := &segment{
		baseOffset: baseOffset,
		config:     c,
	}
	var err error

	// Crear el archivo de store.
	storeFile, err := os.OpenFile(
		path.Join(dir, fmt.Sprintf("%d%s", baseOffset, ".store")),
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, err
	}

	// Crear la instancia del store.
	if s.store, err = newStore(storeFile); err != nil {
		return nil, err
	}

	// Crear el archivo del index.
	indexFile, err := os.OpenFile(
		path.Join(dir, fmt.Sprintf("%d%s", baseOffset, ".index")),
		os.O_RDWR|os.O_CREATE,
		0644,
	)
	if err != nil {
		return nil, err
	}

	// Crear la instancia del index.
	if s.index, err = newIndex(indexFile, c); err != nil {
		return nil, err
	}

	// Leer el último offset para establecer el siguiente.
	if off, _, err := s.index.Read(-1); err != nil {
		s.nextOffset = baseOffset
	} else {
		s.nextOffset = baseOffset + uint64(off) + 1
	}

	return s, nil
}

// Append agrega un registro al store y actualiza el índice.
func (s *segment) Append(record *api.Record) (offset uint64, err error) {
	record.Offset = s.nextOffset

	// Serializar el record con protobuf y escribirlo en el store.
	pos, err := s.store.Append(record)
	if err != nil {
		return 0, err
	}

	// Escribir en el índice la posición del registro en el store.
	if err = s.index.Write(
		uint32(s.nextOffset-uint64(s.baseOffset)), // Offset relativo al baseOffset
		pos,
	); err != nil {
		return 0, err
	}

	s.nextOffset++
	return record.Offset, nil
}

// Read lee un registro desde el store basado en el offset.
func (s *segment) Read(off uint64) (*api.Record, error) {
	// Obtener la posición en el store usando el índice.
	_, pos, err := s.index.Read(int64(off - s.baseOffset))
	if err != nil {
		return nil, err
	}

	// Leer el registro desde la posición en el store.
	return s.store.Read(pos)
}

// IsMaxed determina si el segmento ha alcanzado su tamaño máximo.
func (s *segment) IsMaxed() bool {
	return s.store.size >= s.config.Segment.MaxStoreBytes ||
		s.index.size >= s.config.Segment.MaxIndexBytes
}

// Remove elimina el archivo de índice y el store del sistema.
func (s *segment) Remove() error {
	if err := s.Close(); err != nil {
		return err
	}
	if err := os.Remove(s.index.Name()); err != nil {
		return err
	}
	return os.Remove(s.store.Name())
}

// Close cierra los archivos de índice y store.
func (s *segment) Close() error {
	if err := s.index.Close(); err != nil {
		return err
	}
	return s.store.Close()
}
