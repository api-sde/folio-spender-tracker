package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
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
		if len(dirChild) == 0 {
			fmt.Println("Nothing to do in " + childDirPath)
			continue
		}

		var allTime []Payment

		for _, filePath := range dirChild {
			file, err := os.Open(childDirPath + "/" + filePath.Name())
			defer file.Close()

			if err != nil {
				log.Printf("Unable to read file %s : %s\n", file.Name(), err)
			}

			scanner := bufio.NewScanner(file)

			var sourceInstitution = 0
			var isFirstLine = true
			for scanner.Scan() {
				// CSV
				csvLine := strings.Split(scanner.Text(), ",")
				fmt.Println(csvLine)

				if sourceInstitution == Unknown {
					sourceInstitution = detectSourcePattern(file.Name(), csvLine)
				}

				switch sourceInstitution {
				case Unknown:
					break
				case Tangerine:
					if isFirstLine {
						isFirstLine = false
						continue
					}

					tangerinePayment := convertTangerineLineToPayment(csvLine)
					fmt.Println(tangerinePayment)
					allTime = append(allTime, tangerinePayment)
				}
			}
		}
	}
}

func detectSourcePattern(filename string, line []string) int {
	if strings.Contains(filename, "World Mastercard") &&
		line[0] == "Transaction date" &&
		line[3] == "Memo" {
		return Tangerine
	}
	return Unknown
}

func convertTangerineLineToPayment(csvLine []string) Payment {

	dateSplit := strings.Split(csvLine[0], "/")

	year, _ := strconv.Atoi(dateSplit[2])
	month, _ := strconv.Atoi(dateSplit[1])
	day, _ := strconv.Atoi(dateSplit[0])

	//Date(year int, month Month, day, hour, min, sec, nsec int, loc *Location)
	// To do: implement minute incrementation or iota to keep the order since there is no details coming from the CSV
	stamp := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	tangerinePayment := Payment{
		Stamp: stamp,
		Date:  csvLine[0],
		Year:  dateSplit[2],
		Month: time.Month(month).String(),
		Name:  csvLine[2],
	}

	if csvLine[1] == "DEBIT" {
		amount, errAmount := ParseNewAmount(csvLine[4])
		tangerinePayment.Debit = amount
		if errAmount != nil {
			fmt.Printf("Error while parsing amount line: %s\n", csvLine[4])
		}

		memoRewCat := strings.Split(csvLine[3], "~")
		if len(memoRewCat) > 1 {
			memoCategory := strings.Split(memoRewCat[1], ":")
			tangerinePayment.Category = strings.TrimSpace(memoCategory[1])
		}

		memoReward := strings.Split(memoRewCat[0], ":")
		if len(memoReward) <= 1 {
			fmt.Printf("No cashback found for %s", csvLine)
		} else {
			cashback, errCash := ParseNewAmount(memoReward[1])
			if errCash != nil {
				fmt.Printf("Error while parsing cashback line: %s\n", csvLine[3])
			}

			tangerinePayment.Cashback = cashback
		}
	}

	return tangerinePayment
}
