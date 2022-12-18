package filestorage

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

var (
	ErrFileNotFound = errors.New("file with such id couldn't be found")
)

func (s *Storage) Get(id string, resolution int) ([]byte, string, error) {
	s.wg.Add(1)
	defer s.wg.Done()

	dir := s.dir + id + "/"
	if _, err := os.Stat(dir); err != nil {
		if err == os.ErrNotExist {
			return nil, "", ErrFileNotFound
		}
		return nil, "", err
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, "", ErrFileNotFound
	}

	for _, file := range files {
		parts := strings.Split(file.Name(), ".")
		if parts[0] == s.getImgName(id, resolution) {
			data, err := os.ReadFile(dir + file.Name())
			if err != nil {
				return nil, "", err
			}
			return data, file.Name(), err
		}
	}

	return nil, "", ErrFileNotFound
}
