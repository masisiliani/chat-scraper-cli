package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"

	"github.com/urfave/cli/v2"
)

var (
	filename  string
	chatLines []string
)

func main() {

	app := &cli.App{
		Name:    "chat-scraper-cli",
		Usage:   "Scrape, list, and filter links from exported chat files",
		Version: "0.3.4-beta",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "filename",
				Aliases:     []string{"fn"},
				Required:    true,
				Usage:       "Path to the file (use '-' for stdin)",
				Destination: &filename,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "links",
				Usage: "Operations with links",
				Subcommands: []*cli.Command{
					{
						Name:   "list",
						Usage:  "List links (by category in the future)",
						Action: linksList,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "category",
								Aliases: []string{"c"},
								Value:   "all",
								Usage:   "Categories separated by comma (e.g. linkedin,instagram)",
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

func openFile(name string) (*os.File, error) {
	if name == "-" {
		return os.Stdin, nil
	}

	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", name, err)
	}

	return f, nil

}

func linksList(cCtx *cli.Context) error {

	f, err := openFile(filename)
	if err != nil {
		return err
	}

	if f != os.Stdin {
		defer f.Close()
	}

	var countLines int
	var rawResultSlice []url.URL

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		countLines++
		urlLine, err := extractURL(sc.Text())
		if err != nil {
			return err
		}

		rawResultSlice = append(rawResultSlice, urlLine...)
		chatLines = append(chatLines, sc.Text())

	}

	if err := sc.Err(); err != nil {
		return fmt.Errorf("scanning %s: %w", filename, err)
	}

	resultSlice := stripRawQuery(rawResultSlice)
	for i, v := range resultSlice {
		fmt.Fprintf(cCtx.App.Writer, "%d: %s \n", i, v.String())
	}

	fmt.Fprintf(cCtx.App.Writer, "OK: %d lines read\n", countLines)
	fmt.Fprintf(cCtx.App.Writer, "Links discovered: %d \n", len(resultSlice))

	return nil
}

var urlRe = regexp.MustCompile(`https?://[^\s]+`)

func extractURL(s string) ([]url.URL, error) {

	raws := urlRe.FindAllString(s, -1)
	var result []url.URL
	for _, v := range raws {
		// TODO: Add normalizeURL function to remove punctuation and other special characters

		u, err := url.ParseRequestURI(v)
		if err != nil {
			return nil, fmt.Errorf("URL validatig %s: %w", s, err)
		}
		result = append(result, *u)
	}

	return result, nil
}

func stripRawQuery(urlSlice []url.URL) []url.URL {

	var result []url.URL
	for _, v := range urlSlice {
		v.RawQuery = ""

		result = append(result, v)
	}

	return result
}
