package main

import (
  "os"
  "github.com/codegangsta/cli"
)

func main() {
  app := cli.NewApp()
  app.Name = "got"
  app.Usage = "A DVCS written in golang."
  app.Commands = []cli.Command{
    {
      Name: "init",
      Usage: "Create an empty got repository in the current directory.",
      Action: func(c *cli.Context) {
      },
    },
    {
      Name: "log",
      Usage: "Show commit logs.",
      Action: func(c *cli.Context) {
      },
    },
    {
      Name: "status",
      Usage: "Show the status of the working tree.",
      Action: func(c *cli.Context) {
      },
    },
    {
      Name: "commit",
      Usage: "Record changes to the repository.",
      Action: func(c *cli.Context) {
      },
    },
    {
      Name: "checkout",
      Usage: "Checkout a commit to the working tree.",
      Action: func(c *cli.Context) {
      },
    },
  }

  app.Run(os.Args)
}
