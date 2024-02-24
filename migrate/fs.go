package migrate

import (
	"io/fs"
	"path/filepath"
	"strings"
)

type FSOptions struct {
	UpSuffix   string
	DownSuffix string
	Name       func(string, string, string) string
}

func FromFS(f fs.FS, options ...func(*FSOptions)) ([]Migration, error) {
	opt := FSOptions{
		UpSuffix:   ".up.sql",
		DownSuffix: ".down.sql",
		Name:       CreateMigrationNameFromPath,
	}

	for _, o := range options {
		o(&opt)
	}

	entries, err := readDirectoryEntries(f, ".")
	if err != nil {
		return nil, err
	}

	migrations := make([]*Migration, 0, len(entries))

	for i := 0; i < len(entries); i++ {
		path := entries[i]
		name := opt.Name(path, opt.UpSuffix, opt.DownSuffix)

		var m *Migration
		if len(migrations) > 0 {
			m = migrations[len(migrations)-1]
			if m.Name != name {
				m = &Migration{Name: name}

				migrations = append(migrations, m)
			}
		} else {
			m = &Migration{Name: name}

			migrations = append(migrations, m)
		}

		content, err := fs.ReadFile(f, path)
		if err != nil {
			return nil, err
		}

		if strings.HasSuffix(path, opt.UpSuffix) {
			m.Up = append(m.Up, string(content))
		}
		if strings.HasSuffix(path, opt.DownSuffix) {
			m.Down = append(m.Down, string(content))
		}
	}

	result := make([]Migration, len(migrations))

	for i, item := range migrations {
		result[i] = *item
	}

	return result, nil
}

func readDirectoryEntries(f fs.FS, dir string) ([]string, error) {
	entries, err := fs.ReadDir(f, dir)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			sub, err := readDirectoryEntries(f, filepath.Join(dir, entry.Name()))
			if err != nil {
				return nil, err
			}
			paths = append(paths, sub...)
		} else {
			paths = append(paths, filepath.Join(dir, entry.Name()))
		}
	}

	return paths, nil
}

func CreateMigrationNameFromPath(path, up, down string) string {
	path = strings.TrimSuffix(path, up)
	path = strings.TrimSuffix(path, down)

	return path
}
