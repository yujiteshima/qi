package main

import (
	"errors"
	"os"

	"github.com/yujiteshima/qi/cmd"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "qiitasearch"
	app.Usage = "search qiita articles"
	app.Version = "0.1.0"

	app.Commands = []cli.Command{
		{
			Name:  "myqi",
			Usage: "qiita + mine : you get yours articles",
			Action: func(c *cli.Context) error {
				qiitaToken := os.Getenv("QIITA_TOKEN")
				qiitaUser := os.Getenv("QIITA_USER")
				if qiitaUser == "" {
					return errors.New("$QIITA_USER must be set")
				}

				data, err := cmd.FetchMyQiitaData(qiitaToken, qiitaUser)
				if err != nil {
					return err
				}
				cmd.OutputQiitaData(data)
				return nil
			},
		},
	}
	app.RunAndExitOnError()
}
