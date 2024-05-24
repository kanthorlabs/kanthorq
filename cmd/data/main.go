package main

import (
	"log"

	_ "embed"

	"github.com/kanthorlabs/common/commands/migrate"
	"github.com/kanthorlabs/kanthorq/cmd/base"
	"github.com/kanthorlabs/kanthorq/cmd/data/seed"
)

func main() {
	_, command := base.New()
	command.AddCommand(migrate.New())
	command.AddCommand(seed.New())

	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				log.Println("--- error ---")
				log.Println(err.Error())
			}
		}
	}()

	if err := command.Execute(); err != nil {
		panic(err)
	}
}
