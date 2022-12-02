package main

import (
	"context"

	_ "github.com/joho/godotenv/autoload"
	"github.com/nikprim/banners-rotation/cmd"
)

func main() {
	cmd.Execute(context.Background())
}
