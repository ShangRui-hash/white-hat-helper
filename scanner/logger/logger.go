package logger

import (
	"white-hat-helper/settings"

	"github.com/sirupsen/logrus"
)

func Init() {
	if settings.CurrentConfig.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
