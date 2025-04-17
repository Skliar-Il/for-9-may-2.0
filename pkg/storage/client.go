package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/resourcemanager/v1"
	sdk "github.com/yandex-cloud/go-sdk"
	"github.com/yandex-cloud/go-sdk/iamkey"
	"io"
	"log"
	"net/http"
	"os"
)

func NewSDK() {
	log.Println("🔥 Функция NewSDK() была вызвана")

	// Проверка ключа
	if _, err := os.Stat("authorized_key.json"); os.IsNotExist(err) {
		log.Fatalf("Файл authorized_key.json не найден: %v", err)
	}

	// Чтение ключа
	key, err := iamkey.ReadFromJSONFile("authorized_key.json")
	if err != nil {
		log.Fatalf("Ошибка чтения JSON ключа: %v", err)
	}
	fmt.Printf("key=%s", key)

	// Креды
	creds, err := sdk.ServiceAccountKey(key)
	if err != nil {
		log.Fatalf("Ошибка создания кредов: %v", err)
	}

	// SDK
	ycsdk, err := sdk.Build(context.Background(), sdk.Config{
		Credentials: creds,
	})
	if err != nil {
		log.Fatalf("Ошибка инициализации SDK: %v", err)
	}

	// Выводим облака
	resp, err := ycsdk.ResourceManager().Cloud().List(context.Background(), &resourcemanager.ListCloudsRequest{})
	if err != nil {
		log.Fatalf("Ошибка получения облаков: %v", err)
	}
	for _, cloud := range resp.Clouds {
		log.Printf("Облако: %s (%s)", cloud.Name, cloud.Id)
	}

	// Загрузка файла в Object Storage через HTTP
	log.Println("📤 Загружаем фото в Object Storage...")

	filePath := "./test-photo.jpg"
	bucket := "for9may"
	objectKey := "test-photo.jpg"
	endpoint := "https://storage.yandexcloud.net"
	fullURL := fmt.Sprintf("%s/%s/%s", endpoint, bucket, objectKey)

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Ошибка чтения файла: %v", err)
	}

	// Получаем IAM-токен
	iamToken := ycsdk.IAM()
	fmt.Printf("iam token = %s", iamToken.IamToken)

	// Формируем PUT-запрос
	req, err := http.NewRequest(http.MethodPut, fullURL, bytes.NewReader(fileBytes))
	if err != nil {
		log.Fatalf("Ошибка создания запроса: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", iamToken.IamToken))
	req.Header.Set("Content-Type", "application/octet-stream")

	respHTTP, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Ошибка загрузки файла: %v", err)
	}
	defer respHTTP.Body.Close()

	if respHTTP.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(respHTTP.Body)
		log.Fatalf("Ошибка от Object Storage: %s\n%s", respHTTP.Status, body)
	}

	log.Println("✅ Фото успешно загружено в Object Storage!")
}
