//go:build !generate_completions

package main

import "github.com/CyberVacation/rostra/log"

func main() {
	if err := mainCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}
