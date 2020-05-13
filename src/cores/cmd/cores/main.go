package main

import (
	"flag"
	"os"

	"cos-backend-com/src/cores/cmd/cores/app"

	s "github.com/wujiu2020/strip"
	"github.com/wujiu2020/strip/utils/helpers"
)

func main() {
	strip := s.New()

	var (
		confPath string
		files    helpers.SliceValue
	)
	flagSet := flag.NewFlagSet("app", flag.ExitOnError)
	flagSet.StringVar(&confPath, "conf", "", "")
	flagSet.Var(&files, "f", "")
	flagSet.Parse(os.Args[1:])

	app.AppInit(strip, confPath, files...).Start()
	strip.Run()
}
