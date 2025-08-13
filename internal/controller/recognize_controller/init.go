package recognize_controller

import (
	"MediaTools/internal/config"
	"MediaTools/internal/pkg/wordmatch"
	"regexp"
	"strings"
	"sync"
	"text/template"

	"github.com/sirupsen/logrus"
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

	logrus.Info("开始初始化媒体格式模板...")
	var err error
	movieTemplate, err = template.New("movie").Parse(config.Media.Format.Movie)
	if err != nil {
		return err
	}
	tvTemplate, err = template.New("tv").Parse(config.Media.Format.TV)
	if err != nil {
		return err
	}
	logrus.Info("媒体格式模板初始化完成")
	return nil
}

func InitCustomWord() error {
	loock.Lock()
	defer loock.Unlock()

	logrus.Info("开始初始化自定义识别词...")
	var err error
	wm, err = wordmatch.NewWordsMatcher(config.Media.CustomWord.IdentifyWord)
	if err != nil {
		return err
	}
	re, err := regexp.Compile(strings.Join(config.Media.CustomWord.IdentifyWord, "|"))
	if err != nil {
		return err
	}
	customizationWordRe = re
	logrus.Info("自定义识别词初始化完成")
	return nil
}

func Init() error {
	err := InitFormatTemplates()
	if err != nil {
		return err
	}
	return InitCustomWord()
}
