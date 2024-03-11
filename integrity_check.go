package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func calculateFileHash() (string, error) {
	filePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func readHashFromFile() (string, error) {
	file, err := os.Open("files/hash.txt")
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := make([]byte, 64) // SHA-256 produces 64-character hash
	_, err = file.Read(hash)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func checkIntegrity() {
	currentHash, err := calculateFileHash()
	if err != nil {
		fmt.Println("Ошибка при рассчёте хеша:", err)
		return
	}
	savedHash, err := readHashFromFile()
	if err != nil {
		fmt.Println("Ошибка при чтении сохранённого хеша из файла:", err)
		return
	}
	if currentHash != savedHash {
		fmt.Println("Целостность файла нарушена!")
		// return
	} else {
		fmt.Println("Целостность файла проверена успешно.")
	}

}
