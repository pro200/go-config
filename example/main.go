package main

import (
	"fmt"

	"github.com/pro200/go-config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	result := cfg.SliceString("SLICE", []string{"a", "b"})
	fmt.Println(result)
}
