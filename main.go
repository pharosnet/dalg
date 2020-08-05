package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/pharosnet/dalg/generates"
	"github.com/pharosnet/dalg/logs"
	"github.com/pharosnet/dalg/parser"
)

func main() {
	app := cli.NewApp()
	app.Name = "dalc"
	app.Version = "v2"
	app.Usage = `generate database access layer codes`
	app.ArgsUsage = "dalc --dialect=postgres --schema=schema.sql --query=query.sql --out=path"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "dialect",
			Value:    "",
			Usage:    "sql dialect name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "out",
			Value:    "",
			Usage:    "path of generated codes",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "schema",
			Value:    "",
			Usage:    "ddl",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "query",
			Value:    "",
			Usage:    "dml and dql",
			Required: true,
		},
		&cli.BoolFlag{
			Name:     "json_tags",
			Usage:    "emit json tags",
			Required: false,
			Value:    false,
		},
		&cli.BoolFlag{
			Name:     "verbose",
			Value:    false,
			Usage:    "verbose",
			Required: false,
		},
	}
	app.Action = func(c *cli.Context) (err error) {
		verbose := c.Bool("verbose")
		logs.NewLog(verbose)

		dir, wdErr := os.Getwd()
		if wdErr != nil {
			err = wdErr
			return
		}

		dialect := strings.ToLower(strings.TrimSpace(c.String("dialect")))
		if dialect != "mysql" && dialect != "postgres" {
			err = fmt.Errorf("invallid dialect, it only can be mysql or postgres")
			return
		}

		logs.Log().Println("sql dialect is", dialect)

		out := strings.TrimSpace(c.String("out"))
		out = path.Join(dir, out)
		logs.Log().Println("path of code generated is", out)

		_, pkgName := filepath.Split(out)
		logs.Log().Println("package name of code generated is", pkgName)

		schema := strings.TrimSpace(c.String("schema"))
		schema = path.Join(dir, schema)
		logs.Log().Println("path of schema files is", schema)

		query := strings.TrimSpace(c.String("query"))
		query = path.Join(dir, query)
		logs.Log().Println("path of query files is is", query)

		jsonTags := c.Bool("json_tags")
		logs.Log().Println("emit json tags", jsonTags)

		tables, queries, parseErr := parser.Parse(dialect, schema, query)
		if parseErr != nil {
			err = parseErr
			return
		}
		genErr := generates.Generate(out, jsonTags, tables, queries)
		if genErr != nil {
			err = genErr
			return
		}
		logs.Log().Println("succeed", out)
		return
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}

}
