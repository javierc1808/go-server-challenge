package security

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// LogRotator handles log rotation
type LogRotator struct {
	logDir      string
	maxFiles    int
	maxSize     int64
	rotateDaily bool
}

// NewLogRotator creates a new instance of LogRotator
func NewLogRotator(logDir string, maxFiles int, maxSize int64, rotateDaily bool) *LogRotator {
	return &LogRotator{
		logDir:      logDir,
		maxFiles:    maxFiles,
		rotateDaily: rotateDaily,
		maxSize:     maxSize,
	}
}

// RotateLogs rotates logs according to the configured rules
func (lr *LogRotator) RotateLogs() error {
	// Create directory if it does not exist
	if err := os.MkdirAll(lr.logDir, 0755); err != nil {
		return err
	}

	// Get log files
	files, err := lr.getLogFiles()
	if err != nil {
		return err
	}

	// Rotate by size if necessary
	if err := lr.rotateBySize(files); err != nil {
		return err
	}

	// Rotate by date if enabled
	if lr.rotateDaily {
		if err := lr.rotateByDate(files); err != nil {
			return err
		}
	}

	// Clean up old files
	if err := lr.cleanupOldFiles(files); err != nil {
		return err
	}

	return nil
}

// getLogFiles gets the list of log files
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

	// Sort by modification date (most recent first)
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().After(files[j].ModTime())
	})

	return files, nil
}

// rotateBySize rotates logs by size
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

// rotateByDate rotates logs by date
func (lr *LogRotator) rotateByDate(files []os.FileInfo) error {
	now := time.Now()

	for _, file := range files {
		// Only rotate files from the previous day
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

// cleanupOldFiles deletes old files
func (lr *LogRotator) cleanupOldFiles(files []os.FileInfo) error {
	if len(files) <= lr.maxFiles {
		return nil
	}

	// Delete oldest files
	for i := lr.maxFiles; i < len(files); i++ {
		filePath := filepath.Join(lr.logDir, files[i].Name())
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}

	return nil
}

// getRotatedFileName generates the name of the rotated file
func (lr *LogRotator) getRotatedFileName(originalPath string) string {
	dir := filepath.Dir(originalPath)
	base := filepath.Base(originalPath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	timestamp := time.Now().Format("2006-01-02-15-04-05")
	return filepath.Join(dir, fmt.Sprintf("%s.%s%s", name, timestamp, ext))
}

// GetLogStats returns log statistics
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
