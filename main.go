package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	home, _ := os.UserHomeDir()
	var financeFolderPath string = home + "/Documents/Finance"
	var budgetFolderName string

	if len(os.Args) >= 2 {
		budgetFolderName = os.Args[1]
	} else {
		budgetFolderName = "budget"
	}

	budgetFolderName = financeFolderPath + "/" + budgetFolderName

	fmt.Println("Reading: " + budgetFolderName)
	dirList, err := os.ReadDir(budgetFolderName)

	if err != nil {
		log.Fatal(err)
	}

	for _, dirEntry := range dirList {

		if dirEntry.Type() != os.ModeDir {
			continue
		}

		childDirPath := budgetFolderName + "/" + dirEntry.Name()
		dirChild, err := os.ReadDir(childDirPath)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Found " + strconv.Itoa(len(dirChild)))

		for _, filePath := range dirChild {
			file, err := os.Open(childDirPath + "/" + filePath.Name())
			defer file.Close()

			if err != nil {
				log.Printf("Unable to read file %s : %s\n", file.Name(), err)
			}

			scanner := bufio.NewScanner(file)

			for scanner.Scan() {
				// CSV
				csvLines := strings.Split(scanner.Text(), ",")
				fmt.Println(scanner.Text())
				fmt.Println(csvLines)
			}
		}
	}

}
