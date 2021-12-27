package main

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	m, err := migrate.New(
		"file:///home/shrank/dev/spacey/spacey-backend/migrations/files",
		"mongodb://localhost/users:27017")

	if err != nil {
		fmt.Print(err)

		return
	}

	m.Steps(2)
}
