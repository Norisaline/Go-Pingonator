package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/go-ping/ping"
)

// Генерация диапазона IP-адресов
func generateIPRange(startIP, endIP string) ([]string, error) {
	start := net.ParseIP(startIP)
	end := net.ParseIP(endIP)
	if start == nil || end == nil {
		panic("Неверный пул адресов")
	}

	var ips []string
	for ip := start; !ip.Equal(end); incrementIP(ip) {
		ips = append(ips, ip.String())
	}
	ips = append(ips, end.String()) // Добавляем последний IP в список

	return ips, nil
}

// Увеличение IP-адрес на +1
func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// Пингует IP-адрес и добавляет его в список неудачных при ошибке
func pingAddress(address string, failedAddresses *[]string) {
	pinger, err := ping.NewPinger(address)
	if err != nil {
		panic("Ошибка пинга")
	}

	//Кооличество пакетов, задержка посылки покетац
	pinger.SetPrivileged(true)
	pinger.Count = 2
	pinger.Timeout = 2 * time.Second

	fmt.Printf("Пингую %s...\n", address)
	err = pinger.Run()
	if err != nil || pinger.Statistics().PacketLoss > 0 {
		log.Printf("Ошибка пинга %s:\n\n", address)
		*failedAddresses = append(*failedAddresses, address)
	}
}

// Сохраняет неудачные IP-адреса в файл
func saveFailedAddresses(failedAddresses []string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()

	for _, address := range failedAddresses {
		if _, err := file.WriteString(address + "\n"); err != nil {
			fmt.Println(err.Error())
		}
	}
	return nil
}

func main() {
	var startIP, endIP string

	fmt.Println("\t\t\tВведи начало диапазона пула айпи адреса")
	fmt.Scan(&startIP)

	fmt.Println("\t\t\tВведи конец диапазона пула айпи адреса")
	fmt.Scan(&endIP)

	addresses, err := generateIPRange(startIP, endIP)
	if err != nil {
		fmt.Println(err.Error())
		panic("Ошибка генерации диапазона IP")
	}

	var failedAddresses []string
	for _, address := range addresses {
		pingAddress(address, &failedAddresses)
	}

	if len(failedAddresses) > 0 {
		if err := saveFailedAddresses(failedAddresses, "failed_addresses.txt"); err != nil {
			fmt.Println(err.Error())
			panic("Ошибка сохранения списка неудачных адресов: ")
		}
		fmt.Println("Список неудачных адресов сохранён в файл failed_addresses\n\n")
	}

	//Закрываем программу
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		exit := scanner.Text()
		if exit == "q" {
			break
		} else {
			fmt.Println("\t\t\tPress 'q' to quit")
		}
	}
}
