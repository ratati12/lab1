package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

func PrintHex(val uint64) {
	fmt.Printf("%x\n", val)
}

func main() {
	checkIntegrity()

	// -------------------------------------------------------------------------
	// ----------------------- Чтение ключа из файла ---------------------------
	// -------------------------------------------------------------------------

	key_file, err := os.OpenFile("files/key.txt", os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer key_file.Close()
	// Получаем размер файла ключа
	key_fileInfo, err := key_file.Stat()
	if err != nil {
		fmt.Println("Ошибка получения информации о файле:", err)
		return
	}
	key_fileSize := key_fileInfo.Size()
	if key_fileSize != 32 {
		fmt.Println("Неправильный размер", err)
		return
	}
	// Создаем массив для хранения uint32 для ключа
	var key [8]uint32
	// Читаем данные из файла в массив [8]uint32
	err = binary.Read(key_file, binary.LittleEndian, &key)
	if err != nil {
		fmt.Println("Ошибка чтения данных из файла:", err)
		return
	}

	// -------------------------------------------------------------------------
	// --------------------- Чтение сообщения из файла -------------------------
	// -------------------------------------------------------------------------

	message_file, err := os.OpenFile("files/message.txt", os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer message_file.Close()

	// Получаем размер файла
	message_fileInfo, err := message_file.Stat()
	if err != nil {
		fmt.Println("Ошибка получения информации о файле:", err)
		return
	}
	message_fileSize := message_fileInfo.Size()

	// Вычисляем количество элементов uint64 в файле
	numUint64 := message_fileSize / 8
	if numUint64*8 != message_fileSize {
		numUint64++
		for i := int64(0); i < numUint64*8-message_fileSize; i++ {
			if i == 0 {
				_, err = message_file.Write([]byte{0x01})
			} else {
				_, err = message_file.Write([]byte{0x00})
			}
			if err != nil {
				fmt.Println("Ошибка записи данных в файл:", err)
				return
			}
		}
	}

	message_file, err = os.OpenFile("files/message.txt", os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer message_file.Close()
	// Создаем массив для хранения uint64
	message := make([]uint64, numUint64)
	// Читаем данные из файла в массив uint64
	err = binary.Read(message_file, binary.LittleEndian, &message)
	if err != nil {
		fmt.Println("Ошибка чтения данных из файла:", err)
		return
	}
	// -------------------------------------------------------------------------
	// -------------------------- Вывод OMAC TAG -------------------------------
	// -------------------------------------------------------------------------

	PrintHex(omac(message, key, numUint64))
}
