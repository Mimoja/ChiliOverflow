// Package cache provides methods for caching brand logo images.
package cache

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"ChiliOverFlow/pkg/config"
	"github.com/spf13/viper"
)

var conf *viper.Viper

func init() {
	var err error
	conf, err = config.New()
	if err != nil {
		log.Fatal("Could not get configuration: ", err)
	}
}

// Image queries and saves an image from the provided url if it isn't cached already
func Image(url string, targetPath string) (error) {
	if len(url) == 0 {
		return nil
	}
	imagePath := path.Join(conf.GetString("configPath"), "images")
	file := path.Join(imagePath, targetPath)
	log.Println("Downloading ", url, " to ", file)
	if exists(file) {
		return nil
	}
	response, err := http.Get(url)
	if err != nil {
		log.Fatal("Could not get Image	: ", err)
		return err
	}

	if response.ContentLength == 0 {
		return fmt.Errorf("empty server response")
	}

	err = os.MkdirAll(imagePath, os.ModePerm)
	if err != nil {
		log.Fatal("Could not create folder: ", err)
	}

	f, err := os.Create(file)
	if err != nil {
		log.Fatal("Could not create file: ", err)
	}

	if _, err = io.Copy(f, response.Body); err != nil {
		log.Fatal("Could not write image: ", err)
		return err
	}

	response.Body.Close()
	f.Close()
	return nil
}

// exists returns whether or not a file or directory exists.
func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
