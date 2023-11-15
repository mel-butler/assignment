package numberGenerator

import (
	"crypto/rand"
	"encoding/csv"
	"log"
	"math/big"
	"os"
	"strconv"
)

func CreateFile(min, max, size int) {
	file, err := os.Create("in.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for i := 0; i < size; i++ {
		//number between -9999 and 9999
		randomNumber, err := RandomInt(min, max)
		if err != nil {
			log.Fatal("Failed to generate random number", err)
		}

		// number converted to string and written to csv
		err = writer.Write([]string{strconv.Itoa(randomNumber)})
		if err != nil {
			log.Fatal("Cannot write to file", err)
		}
	}

	log.Println("Numbers generated.")
}

// generating random numbers between a range
func RandomInt(min, max int) (int, error) {
	bigInt, err := RandInt(big.NewInt(int64(max - min + 1)))
	if err != nil {
		return 0, err
	}
	return int(bigInt.Int64()) + min, nil
}

// generating random numbers between 0 and n-1
func RandInt(n *big.Int) (*big.Int, error) {
	randomValue, err := rand.Int(rand.Reader, n)
	if err != nil {
		return nil, err
	}
	return randomValue, nil
}

func ReadFile(fileName string) ([]int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var numbers []int
	for _, row := range records {
		for _, col := range row {
			num, err := strconv.Atoi(col)
			if err != nil {
				return nil, err
			}
			numbers = append(numbers, num)
		}
	}

	return numbers, nil
}

func WriteFile(fileName string, numbers []int) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var records [][]string
	for _, num := range numbers {
		records = append(records, []string{strconv.Itoa(num)})
	}

	return writer.WriteAll(records)
}
