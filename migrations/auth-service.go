//go:build migrations

package main

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	m, err := migrate.New(
		"file:///Users/moritzeich/dev/spacey/backend/migrations/files",
		"mongodb://localhost:27017/users")

	if err != nil {
		fmt.Print(err)

		return
	}

	err = m.Up()

	if err != nil {
		fmt.Print(err)

		return
	}

}
