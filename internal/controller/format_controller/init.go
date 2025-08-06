package format_controller

import (
	"MediaTools/internal/config"
	"sync"
	"text/template"
)

var (
	loock         sync.RWMutex
	movieTemplate *template.Template
	tvTemplate    *template.Template
)

func Init() error {
	loock.Lock()
	defer loock.Unlock()

	var err error
	movieTemplate, err = template.New("movie").Parse(config.MediaLibrary.MovieFormat)
	if err != nil {
		return err
	}
	tvTemplate, err = template.New("tv").Parse(config.MediaLibrary.TVFormat)
	if err != nil {
		return err
	}
	return nil
}
