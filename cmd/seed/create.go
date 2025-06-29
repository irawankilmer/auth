package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

const tmpl = `package seeder

import (
    "database/sql"
    "fmt"
)

func {{.FuncName}}(db *sql.DB) error {
    fmt.Println("Running {{.FuncName}}...")

    // TODO: implement seeding logic
    return nil
}
`

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Gunakan: go run cmd/seed/create.go NamaSeeder")
		return
	}

	name := os.Args[1]
	funcName := strings.Title(name)
	fileName := strings.ToLower(name) + "_seeder.go"

	f, err := os.Create("database/seeders/" + fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	t := template.Must(template.New("seeder").Parse(tmpl))
	t.Execute(f, map[string]string{"FuncName": funcName})

	fmt.Println("Seeder", fileName, "berhasil dibuat")
}
