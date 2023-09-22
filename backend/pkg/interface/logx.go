package util

import "github.com/sirupsen/logrus"

type LogxInterface interface {
	GetLog() *logrus.Entry
}
