package filestorage

import (
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strconv"
)

func (s *Storage) Save(id string, image image.Image, resolution int, imgType string) error {
	s.wg.Add(1)
	defer s.wg.Done()

	switch imgType {
	case "image/jpeg":
		file, err := os.Create(s.dir + s.getImgName(id, resolution) + ".jpeg")
		if err != nil {
			return err
		}
		defer closeFile(file)

		err = jpeg.Encode(file, image, nil)
		if err != nil {
			return err
		}
	case "image/png":
		file, err := os.Create(s.dir + s.getImgName(id, resolution) + ".png")
		if err != nil {
			return err
		}
		defer closeFile(file)

		err = png.Encode(file, image)
		if err != nil {
			return err
		}
	}
	return nil
}

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Printf("error closing file [%s]: [%s]\n", file.Name(), err)
	}

}

func (s *Storage) getImgName(id string, resolution int) string {
	if resolution < 100 {
		return id + "-" + strconv.Itoa(resolution)
	}
	return id
}
