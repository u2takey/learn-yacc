package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v32/github"
	"github.com/olekukonko/tablewriter"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"
)

type env struct {
	sql    *Sql
	client *github.Client
}

func newEnv(sql *Sql) *env {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "xxx"},
	)
	tc := oauth2.NewClient(ctx, ts)
	return &env{
		sql:    sql,
		client: github.NewClient(tc),
	}
}

func (e *env) Run() {
	if e.sql.count == 0 {
		e.sql.count = 10 // default
	}
	repos, _, err := e.client.Repositories.List(context.Background(), "", &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			Page:    e.sql.page,
			PerPage: e.sql.count,
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(e.sql.columns) == 0 || e.sql.columns[0] == "*" {
		e.sql.columns = []string{"Name", "Owner.Login", "Language", "Topics"}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(e.sql.columns)

	for _, repo := range repos {
		var row []string
		for _, h := range e.sql.columns {
			b, _ := json.Marshal(repo)
			r := gjson.Get(string(b), strings.ToLower(h))
			row = append(row, r.String())
		}
		table.Append(row)
	}
	table.Render()
}
