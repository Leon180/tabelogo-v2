package main

import (
	"github.com/Leon180/tabelogo-v2/internal/auth"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		// Load the complete auth module
		auth.Module,
	).Run()
}
