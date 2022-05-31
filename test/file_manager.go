package test

import "errors"

type stubFileManager struct {
	directories map[string][]string
}

// NewStubFileManager creates a mock in-memory file system
func NewStubFileManager() *stubFileManager {
	return &stubFileManager{
		directories: map[string][]string{},
	}
}

func (m *stubFileManager) GetDirectoryContents(directoryPath string) ([]string, error) {
	contents, ok := m.directories[directoryPath]
	if !ok {
		return nil, errors.New("directory does not exist")
	}
	return contents, nil
}

func (m *stubFileManager) CopyFileWithName(directoryPath, fileName, newName string) error {
	contents, err := m.GetDirectoryContents(directoryPath)
	if err != nil {
		return err
	}
	if !inSlice(fileName, contents) {
		return errors.New("file does not exist")
	}

	m.directories[directoryPath] = append(m.directories[directoryPath], newName)

	return nil
}

func (m *stubFileManager) FileExists(directoryPath, fileName string) bool {
	files, ok := m.directories[directoryPath]
	if !ok {
		return false
	}
	for _, f := range files {
		if f == fileName {
			return true
		}
	}
	return false
}

func (m *stubFileManager) SetDirectoryContents(directoryPath string, contents []string) {
	m.directories[directoryPath] = contents
}

func inSlice(s string, slice []string) bool {
	for _, item := range slice {
		if s == item {
			return true
		}
	}
	return false
}
