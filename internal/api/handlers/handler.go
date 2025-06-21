package handlers

import (
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Log *logrus.Logger
}

func New(log *logrus.Logger) *Handler {
	return &Handler{
		Log: log,
	}
}
