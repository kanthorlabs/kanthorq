package main

import (
	"log"

	_ "embed"

	"github.com/kanthorlabs/common/commands/migrate"
	"github.com/kanthorlabs/kanthorq/cmd/base"
)

func main() {
	_, command := base.New()
	command.AddCommand(migrate.New())

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
