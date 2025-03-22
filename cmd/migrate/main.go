package main

import (
	"github.com/theodorusyoga/loan-service-state-machine/migrations"
)

func main() {
	// Running migrations
	migrations.RunMigrations()
}
