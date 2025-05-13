package logger

import (
	"os"

	"go.uber.org/zap"
)

var sugar *zap.SugaredLogger

func Init() {
	var base *zap.Logger
	var err error

	if mode := os.Getenv("GIN_MODE"); mode == "release" {
		base, err = zap.NewProduction()
	} else {
		base, err = zap.NewDevelopment()
	}

	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(base)
	sugar = base.Sugar()
}

func L() *zap.SugaredLogger {
	return sugar
}
