package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
	"github.com/yujiteshima/qiita/api/qiita.go"
)

func main() {
	app := cli.NewApp()
	app.Name = "qiitasearch"
	app.Usage = "search qiita articles"
	app.Action = func(c *cli.Context) error {
		//fmt.Println("Hello friend!")
		api.RunQiitaSearch()
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}