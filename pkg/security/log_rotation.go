package security

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// LogRotator maneja la rotación de logs
type LogRotator struct {
	logDir      string
	maxFiles    int
	maxSize     int64
	rotateDaily bool
}

// NewLogRotator crea una nueva instancia de LogRotator
func NewLogRotator(logDir string, maxFiles int, maxSize int64, rotateDaily bool) *LogRotator {
	return &LogRotator{
		logDir:      logDir,
		maxFiles:    maxFiles,
		rotateDaily: rotateDaily,
		maxSize:     maxSize,
	}
}

// RotateLogs rota los logs según las reglas configuradas
func (lr *LogRotator) RotateLogs() error {
	// Crear directorio si no existe
	if err := os.MkdirAll(lr.logDir, 0755); err != nil {
		return err
	}

	// Obtener archivos de log
	files, err := lr.getLogFiles()
	if err != nil {
		return err
	}

	// Rotar por tamaño si es necesario
	if err := lr.rotateBySize(files); err != nil {
		return err
	}

	// Rotar por fecha si está habilitado
	if lr.rotateDaily {
		if err := lr.rotateByDate(files); err != nil {
			return err
		}
	}

	// Limpiar archivos antiguos
	if err := lr.cleanupOldFiles(files); err != nil {
		return err
	}

	return nil
}

// getLogFiles obtiene la lista de archivos de log
func (lr *LogRotator) getLogFiles() ([]os.FileInfo, error) {
	entries, err := os.ReadDir(lr.logDir)
	if err != nil {
		return nil, err
	}

	var files []os.FileInfo
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".log") {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			files = append(files, info)
		}
	}

	// Ordenar por fecha de modificación (más reciente primero)
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().After(files[j].ModTime())
	})

	return files, nil
}

// rotateBySize rota logs por tamaño
func (lr *LogRotator) rotateBySize(files []os.FileInfo) error {
	for _, file := range files {
		if file.Size() > lr.maxSize {
			oldPath := filepath.Join(lr.logDir, file.Name())
			newPath := lr.getRotatedFileName(oldPath)

			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// rotateByDate rota logs por fecha
func (lr *LogRotator) rotateByDate(files []os.FileInfo) error {
	now := time.Now()

	for _, file := range files {
		// Solo rotar archivos del día anterior
		if now.Sub(file.ModTime()) > 24*time.Hour {
			oldPath := filepath.Join(lr.logDir, file.Name())
			newPath := lr.getRotatedFileName(oldPath)

			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// cleanupOldFiles elimina archivos antiguos
func (lr *LogRotator) cleanupOldFiles(files []os.FileInfo) error {
	if len(files) <= lr.maxFiles {
		return nil
	}

	// Eliminar archivos más antiguos
	for i := lr.maxFiles; i < len(files); i++ {
		filePath := filepath.Join(lr.logDir, files[i].Name())
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}

	return nil
}

// getRotatedFileName genera el nombre del archivo rotado
func (lr *LogRotator) getRotatedFileName(originalPath string) string {
	dir := filepath.Dir(originalPath)
	base := filepath.Base(originalPath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	timestamp := time.Now().Format("2006-01-02-15-04-05")
	return filepath.Join(dir, fmt.Sprintf("%s.%s%s", name, timestamp, ext))
}

// GetLogStats retorna estadísticas de los logs
func (lr *LogRotator) GetLogStats() (map[string]interface{}, error) {
	files, err := lr.getLogFiles()
	if err != nil {
		return nil, err
	}

	totalSize := int64(0)
	fileCount := len(files)

	for _, file := range files {
		totalSize += file.Size()
	}

	return map[string]interface{}{
		"file_count":   fileCount,
		"total_size":   totalSize,
		"max_files":    lr.maxFiles,
		"max_size":     lr.maxSize,
		"rotate_daily": lr.rotateDaily,
	}, nil
}
