package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync/atomic"
	"time"
)

var t byte = '\x00'

func main() {

	//time.Sleep(5000 * time.Millisecond)
	var activeConnections int64

	// Установка порта для прослушивания
	PORT := ":3002"

	// Создание TCP-сервера
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		time.Sleep(5000 * time.Millisecond)
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("Обработчик запущен и слушает порт", PORT)

	for {
		// Принятие входящего соединения
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			time.Sleep(5000 * time.Millisecond)
			continue
		}

		// Обработка соединения в новой горутине
		go handleConnection(c, &activeConnections)
	}
}

func handleConnection(c net.Conn, activeConnections *int64) {
	atomic.AddInt64(activeConnections, 1)
	defer atomic.AddInt64(activeConnections, -1)
	defer c.Close()

	fmt.Println("Новое соединение установлено. Соединений: ", *activeConnections)

	// Чтение сообщения от клиента
	message, err := bufio.NewReader(c).ReadString(t)
	if err != nil {
		fmt.Println("Ошибка при чтении:", err)
		return
	}

	// Отправка ответа клиенту
	response := message + string(t)
	_, err = c.Write([]byte(response))
	if err != nil {
		fmt.Println("Ошибка при отправке ответа:", err)
		return
	}
	fmt.Println("Ответ отправлен клиенту")
}
