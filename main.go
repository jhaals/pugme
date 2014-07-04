package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/jhaals/pugme/pugme"
	"os"
	"strconv"
)

func main() {
	app := cli.NewApp()
	app.Name = "Puginator"
	app.Usage = "a pug a day keeps the doctor away"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:      "pug",
			ShortName: "p",
			Usage:     "Return a random pug",
			Action: func(c *cli.Context) {
				pug := pugme.RandomPugs(1)
				fmt.Println(pug[0])
			},
		},
		{
			Name:      "pugs",
			ShortName: "ps",
			Usage:     "Return [count] multiple pugs",
			Action: func(c *cli.Context) {
				count, err := strconv.Atoi(c.Args().First())
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				pugs := pugme.RandomPugs(count)
				for i := 0; i < count; i++ {
					fmt.Println(pugs[i])
				}
			},
		},
		{
			Name:      "download-pugs",
			ShortName: "dp",
			Usage:     "Download [count] pugs and place them in [path]",
			Action: func(c *cli.Context) {
				count, err := strconv.Atoi(c.Args().First())
				path := c.Args()[1]
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				pugme.DownloadPugs(count, path)
			},
		},
	}

	app.Run(os.Args)
}
