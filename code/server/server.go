package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

var t byte = '\x00'

func main() {
	if in("data\\settings\\is_serve") != "1" {
		fmt.Println("Клиент козёл!")
		time.Sleep(5 * time.Second)

		return
	}

	//time.Sleep(5000 * time.Millisecond)
	var activeConnections int64

	// Установка порта для прослушивания
	PORT := ":" + in("data\\settings\\server_port")

	// Создание TCP-сервера
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		time.Sleep(5000 * time.Millisecond)
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("Сервер запущен и слушает порт", PORT)

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
	response := process_msg(message) + string(t)
	_, err = c.Write([]byte(response))
	if err != nil {
		fmt.Println("Ошибка при отправке ответа:", err)
		return
	}
	fmt.Println("Ответ отправлен клиенту")
}


func in(filePath string) string {
	data, _ := ioutil.ReadFile(filePath) // Чтение файла без проверки ошибок
	return string(data)
}

func process_msg(msg string) string {
	method := before(msg, '.')
	tpe := afterbefore(msg, '.', '/')

	ret := "error:10:unknown_method:" + method

	if method == "test" {
		if tpe == "serve" {
			ret = "true"
		} else {
			ret = "false"
		}
	} else if method == "process" {
		ret = "sucess:" + tcp(in("data\\methods\\" + tpe), after(msg, '/'))
	}

	return ret
}



func out(filename string, text string) error {
    // Получаем путь к директории, где должен находиться файл
    dir := filepath.Dir(filename)

    // Создаем директорию, если она не существует
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }

    // Создаем файл
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    // Записываем текст в файл
    _, err = file.WriteString(text)
    if err != nil {
        return err
    }

    return nil
}


func afterbefore(prompt string, a byte, b byte) (ip string) {
    var builder strings.Builder
    var found bool

    for i := 0; i < len(prompt); i++ {
        if prompt[i] == a {
            found = true
        } else if found && prompt[i]!= b {
            builder.WriteByte(prompt[i])
        } else if found && prompt[i] == b {
            break
        }
    }

    return builder.String()
}

func rnd(max int) (result int) {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}

func after(prompt string, sym byte) (response string) {
	var builder strings.Builder

	for i := 0; i < len(prompt); i++ {
		if(prompt[i] == sym){
			for j := i + 1; j < len(prompt); j++ {
				builder.WriteByte(prompt[j])
			}
		}
	}

	return builder.String()
}

func before(prompt string, sym byte) (response string) {
	var builder strings.Builder

	for i := 0; i < len(prompt) && prompt[i] != sym; i++ {
		builder.WriteByte(prompt[i])
	}

	return builder.String()
}

func lineExistsInFile(filePath string, searchString string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == searchString {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
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