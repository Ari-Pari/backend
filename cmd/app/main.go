package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Ari-Pari/backend/internal/clients/filestorage"
)

func main() {
	endpoint := "127.0.0.1:9000"
	accessKey := "minioadmin"
	secretKey := "minioadmin"
	bucket := "user-photos"

	var store filestorage.FileStorage
	var err error

	store, err = filestorage.NewMinioStorage(endpoint, accessKey, secretKey, bucket, false)
	if err != nil {
		log.Fatalf("Ошибка инициализации: %v", err)
	}

	ctx := context.Background()
	localPath := "testImage/test.jpeg"

	// Открываем файл и готовим данные
	file, err := os.Open(localPath)
	if err != nil {
		log.Fatalf("Файл не найден: %v", err)
	}
	defer file.Close()

	fileStat, _ := file.Stat()
	fileName := filepath.Base(file.Name())

	// Определяем MIME тип (image/png, image/jpeg и т.д.)
	buffer := make([]byte, 512)
	file.Read(buffer)
	contentType := http.DetectContentType(buffer)
	file.Seek(0, 0)

	// 4. ТЕСТ: Загрузка
	fmt.Printf("🚀 Загрузка файла %s...\n", fileName)
	fileKey, err := store.UploadImage(ctx, fileName, file, fileStat.Size(), contentType)
	if err != nil {
		log.Fatalf("Загрузка провалена: %v", err)
	}
	fmt.Printf("✅ Файл загружен с ключом: %s\n", fileKey)

	// 5. ТЕСТ: Получение оригинального имени
	origName, err := store.GetOriginalName(ctx, fileKey)
	if err == nil {
		fmt.Printf("📄 Оригинальное имя в метаданных: %s\n", origName)
	}

	// 6. ТЕСТ: Ссылка на фото (просмотр)
	url, err := store.GetFileURL(ctx, fileKey, 10*time.Minute)
	if err != nil {
		log.Fatalf("Ошибка ссылки: %v", err)
	}

	fmt.Printf("\n🔗 Ссылка для ПРОСМОТРА (inline):\n%s\n", url)
}
