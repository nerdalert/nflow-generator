package main

import (
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func SetLogger(l *logrus.Logger) {
	log = l
}
