package du

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type FileInfo struct {
	Path      string
	IsDir     bool
	Size      int64
	SizeHuman string
}

func FormatSize(size int64) string {
	const (
		_          = iota
		KB float64 = 1 << (10 * iota)
		MB
		GB
		TB
		PB
	)

	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}
	
	var value float64
	var unit string
	
	switch {
	case size >= int64(PB):
		value = float64(size) / PB
		unit = "PB"
	case size >= int64(TB):
		value = float64(size) / TB
		unit = "TB"
	case size >= int64(GB):
		value = float64(size) / GB
		unit = "GB"
	case size >= int64(MB):
		value = float64(size) / MB
		unit = "MB"
	case size >= int64(KB):
		value = float64(size) / KB
		unit = "KB"
	}
	
	return fmt.Sprintf("%.2f %s", value, unit)
}

func CalculateDirSize(path string) (int64, error) {
	var size int64
	
	err := filepath.WalkDir(path, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		if d.IsDir() {
			return nil
		}
		
		info, err := d.Info()
		if err != nil {
			return nil
		}
		
		size += info.Size()
		return nil
	})
	
	return size, err
}

func ListFilesAndDirs(rootPath string, recursive bool) ([]FileInfo, error) {
	var items []FileInfo
	
	entries, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}
	
	for _, entry := range entries {
		path := filepath.Join(rootPath, entry.Name())
		isDir := entry.IsDir()
		
		var size int64
		if isDir {
			dirSize, err := CalculateDirSize(path)
			if err != nil {
				fmt.Printf("Warning: Could not calculate size for %s: %v\n", path, err)
				continue
			}
			size = dirSize
		} else {
			// Get file info to read size
			info, err := entry.Info()
			if err != nil {
				fmt.Printf("Warning: Could not get info for %s: %v\n", path, err)
				continue
			}
			size = info.Size()
		}
		
		items = append(items, FileInfo{
			Path:      path,
			IsDir:     isDir,
			Size:      size,
			SizeHuman: FormatSize(size),
		})
		
		if recursive && isDir {
			subItems, err := ListFilesAndDirs(path, recursive)
			if err != nil {
				fmt.Printf("Warning: Could not process subdirectory %s: %v\n", path, err)
				continue
			}
			items = append(items, subItems...)
		}
	}
	
	sort.Slice(items, func(i, j int) bool {
		return items[i].Size > items[j].Size
	})
	
	return items, nil
}

func PrintFileList(path string, recursive bool) error {
	items, err := ListFilesAndDirs(path, recursive)
	if err != nil {
		return err
	}
	
	var totalSize int64
	for _, item := range items {
		totalSize += item.Size
	}
	
	fmt.Printf("Listing for: %s\n", path)
	fmt.Printf("Total size: %s\n", FormatSize(totalSize))
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%-60s %-10s %s\n", "Path", "Size", "Type")
	fmt.Println(strings.Repeat("-", 80))
	
	for _, item := range items {
		itemType := "File"
		if item.IsDir {
			itemType = "Directory"
		}
		
		relPath, err := filepath.Rel(path, item.Path)
		if err != nil {
			relPath = item.Path
		}
		if relPath == "." {
			relPath = filepath.Base(path)
		}
		
		fmt.Printf("%-60s %-10s %s\n", relPath, item.SizeHuman, itemType)
	}
	
	return nil
}

func GetDU(dir string, recursive bool) {
	if err := PrintFileList(dir, recursive); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
