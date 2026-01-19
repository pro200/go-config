package main

import (
	"fmt"

	"github.com/pro200/go-config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg.String("STRING"))
}
