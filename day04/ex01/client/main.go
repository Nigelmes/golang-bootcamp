package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const url = "https://localhost:8443/buy_candy"

func main() {
	// Флаги командной строки
	candyType := flag.String("k", "", "Candy type")
	candyCount := flag.Int("c", 0, "Candy count")
	money := flag.Int("m", 0, "Money")
	flag.Parse()
	// Загрузка клиентского сертификата и ключа

	cert, err := tls.LoadX509KeyPair("../cert/client/cert.pem", "../cert/client/key.pem")
	if err != nil {
		log.Fatal("Failed to load client certificate and key:", err)
	}

	// Загрузка CA файла
	caCert, err := os.ReadFile("../cert/minica.pem")
	if err != nil {
		log.Fatal("Failed to load CA certificate:", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Настройка конфигурации TLS для клиента
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	// Создание HTTP клиента с настроенной конфигурацией TLS
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// Формирование JSON запроса
	requestBody := fmt.Sprintf(`{"money": %d, "candyType": "%s", "candyCount": %d}`, *money, *candyType, *candyCount)

	// Отправка POST запроса на сервер
	resp, err := client.Post(url, "application/json", strings.NewReader(requestBody))
	if err != nil {
		log.Fatal("Failed to send request:", err)
	}
	defer resp.Body.Close()

	// Чтение и обработка ответа
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read response:", err)
	}
	// Вывод ответа на экран
	fmt.Println(string(responseBody))
}
