package main

import (
	"github.com/Yota-K/golang-hacker-news-cui/ui"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "hn",
		Usage: "This is a tool to see 'Hacker News' made with Go",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "number, n",
				Value: 10,
				Usage: "option for number of Hacker News acquisitions",
			},
		},
		Action: func(c *cli.Context) error {
			ui.HnUi(c.Int("number"))
			return nil
		},
	}

	app.Run(os.Args)
}
