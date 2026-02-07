package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Ari-Pari/backend/internal/storage"
)

func main() {
	// 1. Конфигурация (в идеале грузить из .env)
	cfg := struct {
		Endpoint  string
		AccessKey string
		SecretKey string
		Bucket    string
	}{
		Endpoint:  "127.0.0.1:9000",
		AccessKey: "minioadmin",
		SecretKey: "minioadmin",
		Bucket:    "user-photos",
	}

	// 2. Инициализируем хранилище
	store, err := storage.NewMinioStorage(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey, cfg.Bucket, false)
	if err != nil {
		log.Fatalf("Ошибка инициализации хранилища: %v", err)
	}

	ctx := context.Background()

	// --- ТЕСТОВЫЙ СЛУЧАЙ: Загрузка файла ---
	testFileName := "testImage/test.png"
	file, err := os.Open("testImage/test.png") // Убедитесь, что файл существует
	if err != nil {
		log.Printf("Предупреждение: файл test.png не найден для теста")
	} else {
		defer file.Close()
		fileStat, _ := file.Stat()

		err = store.UploadImage(ctx, testFileName, file, fileStat.Size(), "image/png")
		if err != nil {
			log.Fatalf("Не удалось загрузить: %v", err)
		}
		fmt.Println("✅ Файл успешно загружен")
	}

	// --- ТЕСТОВЫЙ СЛУЧАЙ: Получение ссылки ---
	url, err := store.GetFileURL(ctx, testFileName, time.Hour*24)
	if err != nil {
		log.Fatalf("Не удалось получить ссылку: %v", err)
	}
	fmt.Printf("🔗 Ссылка на файл (24ч): %s\n", url)

	//--- ТЕСТОВЫЙ СЛУЧАЙ: Удаление ---
	//err = store.DeleteFile(ctx, testFileName)
	//if err != nil {
	//	log.Printf("Ошибка удаления: %v", err)
	//}
}
