package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/urfave/cli/v2"
)

var (
	filename     string
	category     string
	urlRawList   []*url.URL
	grouped      bool
	groupedLinks = make(map[string][]*url.URL)
	linesRead    int
)

var categoryDomainMap = map[string][]string{
	LINKEDIN:  {"linkedin"},
	INSTAGRAM: {"instagram"},
	YOUTUBE:   {"youtube"},
}

const (
	LINKEDIN       = "linkedin"
	INSTAGRAM      = "instagram"
	YOUTUBE        = "youtube"
	GENERAL        = "general"
	ALL_CATEGORIES = "all"
)

func main() {

	app := &cli.App{
		Name:    "chat-scraper-cli",
		Usage:   "Scrape, list, and filter links from exported chat files",
		Version: "0.6.2-beta",
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
				Name:   "links",
				Usage:  "Operations with links",
				Before: linksBefore,
				Subcommands: []*cli.Command{
					{
						Name:   "list",
						Usage:  "List links",
						Action: linksList,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:        "category",
								Aliases:     []string{"c"},
								Value:       ALL_CATEGORIES,
								Usage:       "Categories separated by comma (e.g. linkedin,instagram)",
								Destination: &category,
								Action:      linksFilteredByCategory,
							},
							&cli.BoolFlag{
								Name:        "grouped",
								Aliases:     []string{"g"},
								Usage:       "Display links grouped by category (instead of a flat list).",
								Destination: &grouped,
								Action:      groupLinks,
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

func linksBefore(cCtx *cli.Context) error {
	f, err := openFile(filename)
	if err != nil {
		return err
	}

	if f != os.Stdin {
		defer f.Close()
	}

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		urlsPerLine, err := extractURL(sc.Text())
		if err != nil {
			return err
		}

		stripRawQuery(urlsPerLine)
		urlRawList = append(urlRawList, urlsPerLine...)
		linesRead++
	}

	groupedLinks[ALL_CATEGORIES] = append(groupedLinks[ALL_CATEGORIES], urlRawList...)

	return nil
}

func linksList(cCtx *cli.Context) error {
	fmt.Fprintf(cCtx.App.Writer, "OK: %d lines read\n", linesRead)
	fmt.Fprintf(cCtx.App.Writer, "All dicovered: %d \n", len(urlRawList))

	if !grouped {
		fmt.Fprintf(cCtx.App.Writer, "Links showed: %d \n", len(groupedLinks[ALL_CATEGORIES]))
	} else {
		printGroupedLinks(cCtx.App.Writer, LINKEDIN, groupedLinks[LINKEDIN])
		printGroupedLinks(cCtx.App.Writer, INSTAGRAM, groupedLinks[INSTAGRAM])
		printGroupedLinks(cCtx.App.Writer, YOUTUBE, groupedLinks[YOUTUBE])
		printGroupedLinks(cCtx.App.Writer, GENERAL, groupedLinks[GENERAL])
	}

	return nil
}

func linksFilteredByCategory(cCtx *cli.Context, categoriesImput string) error {
	categories := strings.Split(categoriesImput, ",")
	output := []*url.URL{}
	for _, category := range categories {
		if _, ok := categoryDomainMap[category]; ok {
			output = append(output, filterURLByCategory(category, urlRawList)...)
		}
	}

	delete(groupedLinks, ALL_CATEGORIES)
	groupedLinks[ALL_CATEGORIES] = append(groupedLinks[ALL_CATEGORIES], output...)

	return nil
}

func groupLinks(cCtx *cli.Context, groupedCategory bool) error {
	for _, v := range groupedLinks[ALL_CATEGORIES] {
		if strings.Contains(v.Hostname(), categoryDomainMap[LINKEDIN][0]) {
			groupedLinks[LINKEDIN] = append(groupedLinks[LINKEDIN], v)
			continue
		}

		if strings.Contains(v.Hostname(), categoryDomainMap[INSTAGRAM][0]) {
			groupedLinks[INSTAGRAM] = append(groupedLinks[INSTAGRAM], v)
			continue
		}

		if strings.Contains(v.Hostname(), categoryDomainMap[YOUTUBE][0]) {
			groupedLinks[YOUTUBE] = append(groupedLinks[YOUTUBE], v)
			continue
		}

		groupedLinks[GENERAL] = append(groupedLinks[GENERAL], v)

	}

	return nil
}

var urlRe = regexp.MustCompile(`https?://[^\s]+`)

func extractURL(s string) ([]*url.URL, error) {

	raws := urlRe.FindAllString(s, -1)
	var result []*url.URL
	for _, v := range raws {
		// TODO: Add normalizeURL function to remove punctuation and other special characters

		u, err := url.ParseRequestURI(v)
		if err != nil {
			return nil, fmt.Errorf("URL validatig %s: %w", s, err)
		}
		result = append(result, u)
	}

	return result, nil
}

func stripRawQuery(urlsPerLine []*url.URL) {
	for _, v := range urlsPerLine {
		v.RawQuery = ""
	}
}

func filterURLByCategory(categoryFilter string, urlRaw []*url.URL) []*url.URL {

	if categoryFilter == ALL_CATEGORIES {
		return urlRaw
	}

	var filteredUrlsPerLine []*url.URL

	for _, v := range urlRaw {
		switch categoryFilter {
		case LINKEDIN:
			if strings.Contains(v.Hostname(), categoryDomainMap[LINKEDIN][0]) {
				filteredUrlsPerLine = append(filteredUrlsPerLine, v)
			}
		case INSTAGRAM:
			if strings.Contains(v.Hostname(), categoryDomainMap[INSTAGRAM][0]) {
				filteredUrlsPerLine = append(filteredUrlsPerLine, v)
			}
		case YOUTUBE:
			if strings.Contains(v.Hostname(), categoryDomainMap[YOUTUBE][0]) {
				filteredUrlsPerLine = append(filteredUrlsPerLine, v)
			}
		}
	}

	return filteredUrlsPerLine

}

func printGroupedLinks(writer io.Writer, category string, links []*url.URL) {
	fmt.Fprintf(writer, "\n %s\n", strings.ToUpper(category))
	for i, v := range links {
		fmt.Fprintf(writer, "\t%d: %s \n", i+1, v)
	}
}
