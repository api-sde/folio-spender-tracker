package main

import (
	"bufio"
	"encoding/csv"
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

		fmt.Println("Found " + strconv.Itoa(len(dirChild)) + " files.")
		if len(dirChild) == 0 {
			fmt.Println("Nothing to do in " + childDirPath)
			continue
		}

		var allTime []Payment
		var allFilesLinesCount int
		var headerNomenclatureLines int
		var missingFiles []string

		for _, filePath := range dirChild {
			if filePath.Name() == ".DS_Store" {
				continue
			}
			file, err := os.Open(childDirPath + "/" + filePath.Name())
			defer file.Close()

			if err != nil {
				log.Printf("Unable to read file %s : %s\n", file.Name(), err)
				missingFiles = append(missingFiles, file.Name())
			}

			scanner := bufio.NewScanner(file)

			var sourceInstitution = 0
			var isFirstLine = true
			for scanner.Scan() {

				allFilesLinesCount = allFilesLinesCount + 1

				// CSV
				csvLine := strings.Split(scanner.Text(), ",")

				if sourceInstitution == Unknown {
					sourceInstitution = detectSourcePattern(file.Name(), csvLine)
				}

				switch sourceInstitution {
				case Unknown:
					break
				case Tangerine:
					if isFirstLine {
						headerNomenclatureLines = headerNomenclatureLines + 1
						isFirstLine = false
						continue
					}

					tangerinePayment := convertTangerineLineToPayment(csvLine)
					fmt.Println(tangerinePayment)
					allTime = append(allTime, tangerinePayment)
				case CIBC:
					var cibcPayment = convertCIBCLineToPayment(csvLine)
					allTime = append(allTime, cibcPayment)
				}
			}
		}

		ValidateFilesAndLinesSum(len(allTime), allFilesLinesCount, headerNomenclatureLines, missingFiles)

		var recordCsv [][]string
		header := []string{"Stamp", "Date", "Year", "Month", "Name", "Category", "Cashback", "Debit", "Credit"}
		recordCsv = append(recordCsv, header)

		for _, payment := range allTime {

			debit := ""
			credit := ""
			cashback := ""

			if payment.Debit != nil {
				debit = payment.Debit.ToText()
			}
			if payment.Credit != nil {
				credit = payment.Credit.ToText()
			}
			if payment.Cashback != nil {
				cashback = payment.Cashback.ToText()
			}

			csvLine := []string{
				strconv.Itoa(int(payment.Stamp.Unix())),
				payment.Date,
				payment.Year,
				payment.Month,
				payment.Name,
				payment.Category,
				cashback,
				debit,
				credit,
			}

			recordCsv = append(recordCsv, csvLine)
		}

		resultFilename := "testResult.csv"
		file, err := os.Create(childDirPath + "/" + resultFilename)
		defer file.Close()
		if err != nil {
			fmt.Println("failed to create CSV file")
		}

		csvWriter := csv.NewWriter(file)
		defer csvWriter.Flush()

		err = csvWriter.WriteAll(recordCsv)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Created " + resultFilename)
	}

}

func detectSourcePattern(filename string, line []string) int {
	if strings.Contains(filename, "World Mastercard") &&
		line[0] == "Transaction date" &&
		line[3] == "Memo" {
		return Tangerine
	} else if strings.Contains(filename, "cibc") {
		return CIBC
	}
	return Unknown
}

func convertTangerineLineToPayment(csvLine []string) Payment {

	dateSplit := strings.Split(csvLine[0], "/")

	year, _ := strconv.Atoi(dateSplit[2])
	month, _ := strconv.Atoi(dateSplit[0])
	day, _ := strconv.Atoi(dateSplit[1])

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

	if csvLine[1] == DEBIT {
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
	} else if csvLine[1] == CREDIT {
		fmt.Println(csvLine)

		creditAmount, errAmount := ParseNewAmount(csvLine[4])
		tangerinePayment.Credit = creditAmount
		if errAmount != nil {
			fmt.Printf("Error while parsing amount line: %s\n", csvLine[4])
		}

	}

	return tangerinePayment
}

func convertCIBCLineToPayment(csvLine []string) Payment {

	dateSplit := strings.Split(csvLine[0], "-")

	year, _ := strconv.Atoi(dateSplit[0])
	month, _ := strconv.Atoi(dateSplit[1])
	day, _ := strconv.Atoi(dateSplit[2])

	stamp := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	cibcPayment := Payment{
		Stamp: stamp,
		Date:  csvLine[0], // To do: get same format everywhere
		Year:  dateSplit[0],
		Month: time.Month(month).String(),
		Name:  csvLine[1],
	}

	var parsing error

	if strings.Contains(csvLine[1], "PAYMENT THANK YOU") {
		credit, err := ParseNewAmount(csvLine[3])
		parsing = err
		cibcPayment.Credit = credit
	} else if csvLine[1] == "CASHBACK/REMISE EN ARGENT" {
		credit, err := ParseNewAmount(csvLine[3])
		parsing = err
		cibcPayment.Credit = credit
	} else if csvLine[1] == "ANNUAL FEE" {
		debit, err := ParseNewAmount(csvLine[2])
		parsing = err
		cibcPayment.Debit = debit
	} else if csvLine[3] != "" {
		debit, err := ParseNewAmount(csvLine[3])
		parsing = err
		cibcPayment.Debit = debit
	} else if csvLine[4] != "" {
		credit, err := ParseNewAmount(csvLine[4])
		parsing = err
		cibcPayment.Credit = credit
	} else {
		fmt.Printf("Line %s has been ignored.", csvLine)
	}

	if parsing != nil {
		fmt.Printf("Error while parsing amount line: %s\n, %s", csvLine, parsing.Error())
	}

	return cibcPayment
}

func ValidateFilesAndLinesSum(allLinesCount int, allFilesLinesCount int, headerLines int, missingFiles []string) {
	fmt.Printf("Done processing, found %v lines\n", allLinesCount)
	fmt.Printf("Number of lines across all files: %v, with nomenclature lines: %v \n", allFilesLinesCount, headerLines)

	if allLinesCount+headerLines == allFilesLinesCount {
		fmt.Println("Validation OK, successfully read all lines.")
	} else {
		fmt.Println("Invalid files, some lines are missing.")
	}

	if len(missingFiles) > 0 {
		fmt.Println("Error: some files were missing:")
		for _, miss := range missingFiles {
			fmt.Println(miss)
		}
	}
}
