package files

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileManager struct{}

// GetDirectoryContents returns a list of all the files in the given directory
func (m *FileManager) GetDirectoryContents(directoryPath string) ([]string, error) {
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		return nil, err
	}
	fileNames := []string{}
	for _, f := range files {
		fileNames = append(fileNames, f.Name())
	}
	return fileNames, nil
}

// CopyFileWithName copies a file with a new name in the given directory
func (m *FileManager) CopyFileWithName(directoryPath, srcFileName, newFileName string) error {
	src := filepath.Join(directoryPath, srcFileName)
	dst := filepath.Join(directoryPath, newFileName)

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)

	return err
}

// FileExists checks if a file exists in the given directory
func (m *FileManager) FileExists(directoryPath, fileName string) bool {
	_, err := os.Stat(filepath.Join(directoryPath, fileName))
	return err == nil
}
