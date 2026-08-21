package prefs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const fileName = ".cursor-utils.json"

type File struct {
	CursorRoot string `json:"cursorRoot"`
}

func Path(home string) string {
	return filepath.Join(home, fileName)
}

func Load(home string) (*File, error) {
	data, err := os.ReadFile(Path(home))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to read saved installer settings: %w", err)
	}
	var f File
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("saved installer settings are damaged")
	}
	return &f, nil
}

func Save(home string, f File) error {
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	path := Path(home)
	tmp, err := os.CreateTemp(filepath.Dir(path), ".cursor-utils.json.*.tmp")
	if err != nil {
		return fmt.Errorf("unable to save installer settings: %w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("unable to save installer settings: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("unable to save installer settings: %w", err)
	}
	return nil
}
