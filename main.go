package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	bofacc "github.com/muly/bank-tx/bofa_cc"
)

func main() {
	files := []string{}
	err := filepath.WalkDir("./temp/bofa-cc/2023", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			// fmt.Println(path)
			files = append(files, path)

		}
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}

	statements := make([]bofacc.Statement, 0, len(files))

	for _, file := range files {
		// if file != "temp/bofa-cc/2023/plain_bofa-cc-1810-2023-08-11"{
		// 	continue
		// }
		fmt.Printf("Processing file %s #########################\n", file)
		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		s, err := bofacc.ParseStatement(string(data))
		if err != nil {
			fmt.Println(err)
			return
		}

		statements = append(statements, *s)
		// if err := bofacc.SaveTransactionsToCSV(*s); err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
	}

	if err := bofacc.SaveStatementsToCSV(statements, "temp/bofa-cc-2023.csv"); err != nil {
		fmt.Println(err)
		return
	}

}
