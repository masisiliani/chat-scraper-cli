package main

import (
	"fmt"
	"log"
	"os"
	"bufio"

	"github.com/urfave/cli/v2"
)

var filename string 

func main() {

	app := &cli.App{
		Commands: []*cli.Command{
			Name: "chat-scraper-cli", 
			Usage: "Scrape, list, and filter links from exported chat files",
			{
				Name:  "file",
				Usage: "Operations with chat file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "filename",
						Aliases:  []string{"fn"},
						Required: true,
						Usage:   "Path to the file (use '-' for stdin)",
						Destination: &filename,
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:   "dump",
						Usage:  "Print file content to the terminal",
						Action: fileDump,	
					},
					{
						Name:   "load",
						Usage:  "Load file content into memory",
						Action: fileLoad,		
					},
				},
			},
			{
				Name:   "links",
				Usage:  "Operations with links",
				Subcommands: []*cli.Command{
					{
						Name:   "list",
						Usage:  "List links (by category in the future)",
						Action: linksList,,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "category",
								Aliases:  []string{"c"},
								Value:    "all",
								Usage:   "Categories separated by comma (e.g., linkedin,instagram)",
							},
							&cli.StringFlag{
								Name:    "filename",
								Aliases: []string{"fn"},
								Usage:   "Chat file (same format as the 'file' command)",
							},
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func fileDump(cCtx *cli.Context) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}	

	return nil
}

func fileLoad(cCtx *cli.Context) error {
	fmt.Fprintf(cCtx.App.Writer, "Function *file load*")
	return nil
}

func linksList(cCtx *cli.Context) error {
	fmt.Fprintf(cCtx.App.Writer, "Function *links list*")
	return nil
}