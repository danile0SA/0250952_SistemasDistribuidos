package index

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/tysonmote/gommap"
)

var (
	offWidth uint64 = 4                   // Tamaño del offset en bytes
	posWidth uint64 = 8                   // Tamaño de la posición en bytes
	entWidth        = offWidth + posWidth // Tamaño total de la entrada (offset + posición)
)

type Config struct {
	Segment struct {
		MaxIndexBytes uint64 // Tamaño máximo del índice en bytes
	}
}

type index struct {
	file *os.File    // Archivo donde se guarda el índice
	mmap gommap.MMap // Mapeo directo entre memoria y archivo
	size uint64      // Tamaño del índice en bytes
	mu   sync.Mutex  // Mutex para acceso concurrente seguro
}

// newIndex crea un nuevo índice a partir de un archivo y una configuración
func newIndex(f *os.File, c Config) (*index, error) {
	idx := &index{
		file: f,
	}

	// Obtener el tamaño del archivo
	fi, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("could not get file stats: %w", err)
	}
	idx.size = uint64(fi.Size())

	// Ajustar el tamaño del archivo según Config.Segment.MaxIndexBytes
	if err := os.Truncate(f.Name(), int64(c.Segment.MaxIndexBytes)); err != nil {
		return nil, fmt.Errorf("could not truncate file: %w", err)
	}

	// Hacer el mapeo entre archivo y memoria
	if idx.mmap, err = gommap.Map(idx.file.Fd(), gommap.PROT_READ|gommap.PROT_WRITE, gommap.MAP_SHARED); err != nil {
		return nil, fmt.Errorf("could not map file: %w", err)
	}

	return idx, nil
}

// Read lee la entrada en una posición determinada
func (idx *index) Read(inPos int64) (outOffset, outPos uint64, err error) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	if idx.size == 0 {
		return 0, 0, fmt.Errorf("index is empty")
	}

	// Caso especial para inPos == -1: devolver el último elemento
	if inPos == -1 {
		if idx.size < entWidth {
			return 0, 0, fmt.Errorf("index too small")
		}
		inPos = int64(idx.size/entWidth) - 1
	}

	if inPos*int64(entWidth) >= int64(idx.size) {
		return 0, 0, io.EOF
	}

	start := uint64(inPos) * entWidth
	outOffset = uint64(binary.BigEndian.Uint32(idx.mmap[start : start+offWidth])) // Leemos el offset
	outPos = binary.BigEndian.Uint64(idx.mmap[start+offWidth : start+entWidth])   // Leemos la posición

	return outOffset, outPos, nil
}

// Write escribe una nueva entrada en el índice
func (idx *index) Write(off, pos uint64) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	// Verificar si hay espacio suficiente para escribir
	if idx.size+entWidth > uint64(len(idx.mmap)) {
		return fmt.Errorf("no space left to write index entry")
	}

	// Escribir el offset
	binary.BigEndian.PutUint32(idx.mmap[idx.size:idx.size+offWidth], uint32(off))

	// Escribir la posición
	binary.BigEndian.PutUint64(idx.mmap[idx.size+offWidth:idx.size+entWidth], pos)

	// Actualizar el tamaño del índice
	idx.size += entWidth

	return nil
}

// Close cierra el índice después de asegurarse de escribir todos los datos y truncar el archivo
func (idx *index) Close() error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	// Sincronizar el contenido del mmap con el archivo
	if err := idx.mmap.Sync(gommap.MS_SYNC); err != nil {
		return fmt.Errorf("could not sync memory map: %w", err)
	}

	// Truncar el archivo al tamaño real del índice
	if err := os.Truncate(idx.file.Name(), int64(idx.size)); err != nil {
		return fmt.Errorf("could not truncate file: %w", err)
	}

	// Cerrar el archivo
	if err := idx.file.Close(); err != nil {
		return fmt.Errorf("could not close file: %w", err)
	}

	return nil
}

// Name devuelve el nombre del archivo del índice
func (idx *index) Name() string {
	return idx.file.Name()
}
