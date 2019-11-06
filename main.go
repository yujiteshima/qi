package main

import (
	"os"

	"./cmd"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "qiitasearch"
	app.Usage = "search qiita articles"
	app.Version = "0.1.0"

	app.Commands = []cli.Command{
		{
			Name:  "mine",
			Usage: "qiita + mine : you get yours articles",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "example, e",
					Value: "exmple",
					Usage: "This help text is example",
				},
			},
			Action: func(c *cli.Context) error {
				//fmt.Println("Hello friend!")
				qiitaToken := os.Getenv("QIITA_TOKEN")
				// Debug用 fmt.Println(qiitaToken)
				//id_ary := getId()
				datas := cmd.FetchQiitaData(qiitaToken)
				cmd.OutputQiitaData(datas)
				return nil
			},
		},
	}
	app.Run(os.Args)
	// if err != nil {
	// 	log.Fatal(err)
	// }

}
