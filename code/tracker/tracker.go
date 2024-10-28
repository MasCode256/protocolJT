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
	//time.Sleep(5000 * time.Millisecond)
	if in("data\\settings\\is_track") != "1" {
		fmt.Println("Клиент козёл!")
		time.Sleep(5 * time.Second)

		return
	}


	var activeConnections int64

	// Установка порта для прослушивания
	PORT := ":" + in("data\\settings\\tracker_port")

	// Создание TCP-сервера
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		time.Sleep(5000 * time.Millisecond)
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("Трекер запущен и слушает порт", PORT)

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

	response := process_msg(message) + string(t)
	// Отправка ответа клиенту
	_, err = c.Write([]byte(response))
	if err != nil {
		fmt.Println("Ошибка при отправке ответа:", err)
		return
	}
	//fmt.Println("Ответ отправлен клиенту")
}

func process_msg(msg string) string {
	//fmt.Println("user ip: " + ip)
	
	ret := "error:10:unknown_method"
	method := before(msg, '.')
	
	if method == "add" {
		ip := afterbefore(msg, '/', '\x00')
		tpe := afterbefore(msg, '.', '/')

		if tpe == "tracker" {
			add, err := lineExistsInFile("data\\lists\\tracklist", ip)

			res := tcp(ip, "test.track/")
			add2 := (res == "true" || res == "true " || res == "true\x00")

			if err == nil {
				if !add {
					if add2 {
						out("data\\lists\\tracklist", in("data\\lists\\tracklist") + "\n" + ip)
						ret = "sucess"
					} else {
						ret = "error:14:ip_is_not_tracker"
					}
				} else {
					ret = "error:12:tracker_is_alredy_exists"
				}
			} else {
				out("data\\lists\\tracklist", "")
				ret = "error:30:tracker_internal_error"
			}
		} else if tpe == "server" {
			add, err := lineExistsInFile("data\\lists\\servelist", ip)

			res := tcp(ip, "test.serve/")
			add2 := (res == "true" || res == "true " || res == "true\x00")

			if err == nil {
				if !add {
					if add2 {
						out("data\\lists\\servelist", in("data\\lists\\servelist") + "\n" + ip)
						ret = "sucess"
					} else {
						ret = "error:14:ip_is_not_server"
					}
				} else {
					ret = "error:12:server_is_alredy_exists"
				}
			} else {
				out("data\\lists\\servelist", "")
				ret = "error:30:tracker_internal_error"
			}
		} else {
			ret = "error:11:unknown_type:" + tpe
		}
	} else if method == "test" {
		tpe := afterbefore(msg, '.', '/')
		if tpe == "track" {
			ret = "true"
		} else {
			ret = "false"
		}

		fmt.Println(ret)
	}

	return ret
}


func in(filePath string) string {
	data, _ := ioutil.ReadFile(filePath) // Чтение файла без проверки ошибок
	return string(data)
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