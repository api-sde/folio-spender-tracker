package main

import (
	"fmt"
	"os"
)

func main() {

	var budgetFolderName string

	if len(os.Args) >= 2 {
		budgetFolderName = os.Args[1]
	} else {
		budgetFolderName = "budget"
	}

	dirList, err := os.ReadDir("/" + budgetFolderName)

	if err != nil {
		//log.Fatal(err)
	}

	fmt.Println(dirList)
	fmt.Println(os.UserHomeDir())
	fmt.Println(os.UserConfigDir())
}
