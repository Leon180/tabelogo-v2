package main

import (
	"github.com/Leon180/tabelogo-v2/internal/restaurant"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		restaurant.Module,
	).Run()
}
