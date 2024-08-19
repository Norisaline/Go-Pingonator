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
		return nil, fmt.Errorf("Неверный IP-адрес")
	}

	var ips []string
	for ip := start; !ip.Equal(end); incrementIP(ip) {
		ips = append(ips, ip.String())
	}
	ips = append(ips, end.String()) // Добавляем последний IP в список

	return ips, nil
}

// Увеличение IP-адреса
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
		log.Printf("Ошибка создания пингера для %s: %v", address, err)
		*failedAddresses = append(*failedAddresses, address)
		return
	}

	pinger.SetPrivileged(true)
	pinger.Count = 2
	pinger.Timeout = 2 * time.Second

	fmt.Printf("\nПингую %s...\n", address)
	err = pinger.Run()
	if err != nil || pinger.Statistics().PacketLoss > 0 {
		log.Printf("\nОшибка пинга %s:\n", address)
		*failedAddresses = append(*failedAddresses, address)
	}
}

// Сохраняет неудачные IP-адреса в файл
func saveFailedAddresses(failedAddresses []string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, address := range failedAddresses {
		if _, err := file.WriteString(address + "\n"); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var startIP, endIP string

	fmt.Println("Введи начало диапазона пула айпи адреса")
	fmt.Scan(&startIP)

	fmt.Println("Введи конец диапазона пула айпи адреса")
	fmt.Scan(&endIP)

	addresses, err := generateIPRange(startIP, endIP)
	if err != nil {
		log.Fatalf("Ошибка генерации диапазона IP-адресов: %v", err)
	}

	var failedAddresses []string
	for _, address := range addresses {
		pingAddress(address, &failedAddresses)
	}

	if len(failedAddresses) > 0 {
		if err := saveFailedAddresses(failedAddresses, "failed_addresses.txt"); err != nil {
			log.Fatalf("Ошибка сохранения списка неудачных адресов: %v", err)
		}
		fmt.Println("Список неудачных адресов сохранён в файл failed_addresses.txt")
	} else {
		fmt.Println("Все адреса успешно пингуются.")
	}

	fmt.Println("Press 'q' to quit")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		exit := scanner.Text()
		if exit == "q" {
			break
		} else {
			fmt.Println("Press 'q' to quit")
		}
	}
}
