package media_controller

import (
	"MediaTools/internal/config"
	"MediaTools/internal/pkg/wordmatch"
	"regexp"
	"strings"
	"sync"
	"text/template"
)

var (
	loock               sync.RWMutex
	movieTemplate       *template.Template
	tvTemplate          *template.Template
	wm                  *wordmatch.WordsMatcher
	customizationWordRe *regexp.Regexp
)

func InitFormatTemplates() error {
	loock.Lock()
	defer loock.Unlock()

	var err error
	movieTemplate, err = template.New("movie").Parse(config.Media.Format.Movie)
	if err != nil {
		return err
	}
	tvTemplate, err = template.New("tv").Parse(config.Media.Format.TV)
	if err != nil {
		return err
	}
	return nil
}

func InitCustomWord() error {
	loock.Lock()
	defer loock.Unlock()

	var err error
	wm, err = wordmatch.NewWordsMatcher(config.CustomWord.IdentifyWord)
	if err != nil {
		return err
	}
	re, err := regexp.Compile(strings.Join(config.CustomWord.IdentifyWord, "|"))
	if err != nil {
		return err
	}
	customizationWordRe = re
	return nil
}

func Init() error {
	err := InitFormatTemplates()
	if err != nil {
		return err
	}
	return InitCustomWord()
}
