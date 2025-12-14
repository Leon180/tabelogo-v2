package main

import (
	"github.com/Leon180/tabelogo-v2/internal/spider"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		spider.Module,
	).Run()
}
