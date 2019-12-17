package logger

import (
	"context"
	"fmt"
	"log"
)

func Error(ctx context.Context, err error) {
	log.Println(err.Error())
}

func Fatal(ctx context.Context, err error) {
	log.Fatal(err)
}

func Warn(ctx context.Context, msg string) {
	log.Println(msg)
}

func Info(ctx context.Context, msg string) {
	log.Println(msg)
}

func Infof(ctx context.Context, format string, a ...interface{}) {
	log.Println(fmt.Sprintf(format, a...))
}
