package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

var t byte = '\x00'
func main() {
	args := os.Args

	if len(args) > 3 && args[1] == "jtcp" {
		fmt.Print(tcp(args[2], args[3]))
	}
}

func tcp(address string, message string) (string) {
	// Устанавливаем соединение с TCP-сервером
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return ""
	}
	defer conn.Close()

	// Отправляем сообщение серверу
	_, err = conn.Write([]byte(message + string(t)))
	if err != nil {
		return ""
	}

	// Читаем ответ от сервера
	response, err := bufio.NewReader(conn).ReadString(t)
	if err != nil {
		fmt.Println("Ошибка при получении ответа:", err)
		time.Sleep(5000 * time.Millisecond)
		os.Exit(1)
	}

	// Возвращаем ответ в виде строки
	return response
}