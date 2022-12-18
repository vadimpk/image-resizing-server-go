package filestorage

import (
	"github.com/vadimpk/image-resizing-server-go/internal/api/delivery/http"
	"io/ioutil"
	"os"
	"strings"
)

func (s *Storage) Get(id string, resolution int) ([]byte, string, error) {
	s.wg.Add(1)
	defer s.wg.Done()

	dir := s.dir + id + "/"
	if _, err := os.Stat(dir); err != nil {
		if err == os.ErrNotExist {
			return nil, "", http.ErrFileNotFound
		}
		return nil, "", err
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, "", http.ErrFileNotFound
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

	return nil, "", http.ErrFileNotFound
}
