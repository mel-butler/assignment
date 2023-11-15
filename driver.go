package main

import (
	"crypto/rand"
	"encoding/csv"
	"log"
	"math/big"
	"os"
	"strconv"
)

func main() {
	file, err := os.Create("in.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for i := 0; i < 10000; i++ {
		// Generate a random number between -9999 and 9999
		randomNumber, err := randomInt(-9999, 9999)
		if err != nil {
			log.Fatal("Failed to generate random number", err)
		}

		// Convert the number to string and write it to CSV
		err = writer.Write([]string{strconv.Itoa(randomNumber)})
		if err != nil {
			log.Fatal("Cannot write to file", err)
		}
	}

	log.Println("Numbers generated.")
}

// generating random numbers between a range
func randomInt(min, max int) (int, error) {
	// Generate a cryptographically secure random number
	bigInt, err := randInt(big.NewInt(int64(max - min + 1)))
	if err != nil {
		return 0, err
	}
	return int(bigInt.Int64()) + min, nil
}

// generating random numbers between 0 and n-1
func randInt(n *big.Int) (*big.Int, error) {
	randomValue, err := rand.Int(rand.Reader, n)
	if err != nil {
		return nil, err
	}
	return randomValue, nil
}
