package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func getFiles(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	return fileNames, err

}

func runMigration(migrationFilePath string, dbName string, conString string) error {
	fmt.Print(migrationFilePath)
	m, err := migrate.New(
		migrationFilePath,
		conString+dbName)

	if err != nil {
		return err
	}

	if err = m.Up(); err != nil {
		return err
	}
	return nil
}

func main() {
	var serviceName string
	var dbName string
	flag.StringVar(&serviceName, "service", "", "service name")
	flag.StringVar(&dbName, "db", "", "database name")

	flag.Parse()

	if serviceName == "all" {
		log.Println("Running migrations for all services")

	} else {

		if serviceName == "" {
			log.Println("service name is required")
			return
		}

		if dbName == "" {
			log.Println("database name is required")
			return
		}
	}

	err := godotenv.Load("../.env")

	if err != nil {
		log.Print("Error loading .env file in project root.")
		return
	}

	projectRoot := os.Getenv("PROJECT_ROOT")
	connectionString := os.Getenv("MONGO_DB_CONNECTION")

	if projectRoot == "" {
		log.Println("PROJECT_ROOT is required")
		return
	}

	if serviceName == "all" {
		fileNames, err := getFiles(projectRoot + "/migrations")

		if err != nil {
			log.Println(err)
			return
		}

		for _, fileName := range fileNames {

			envVarName := strings.ToUpper(fileName) + "_DB_NAME"
			envVarName = strings.Replace(envVarName, "-", "_", -1)

			dbName = os.Getenv(envVarName)

			if dbName == "" {
				log.Println(
					"<SERVICE_NAME>_DB_NAME is required as environment variable. Skipping migration for service: " + fileName,
				)
				return
			}

			log.Println("Running migrations for service:", fileName, "on database:", dbName)

			if err := runMigration("file://"+projectRoot+"migrations/"+fileName, dbName, connectionString); err != nil {
				log.Println(err)
				return
			}

			log.Println("Migration completed for service:", fileName)

		}

	} else {
		migrationFilesPath := "file://" + projectRoot + "migrations/" + serviceName
		log.Println("Running migrations for service:", serviceName, "on database:", dbName)

		if err := runMigration(migrationFilesPath, dbName, connectionString); err != nil {
			log.Println(err)
			return
		}

		log.Println("Migrations completed successfully")

	}

}
