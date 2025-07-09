package main

import (
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/config"
	"github.com/sirupsen/logrus"
)

func main() {
	err := config.LoadEnv()

	if err != nil {
		logrus.Warn("Tidak dapat memuat file .env, menggunakan variabel lingkungan sistem.")
	} else {
		logrus.Info("File .env berhasil dimuat.")
	}

	app, err := InitializeApp()
	if err != nil {
		logrus.Fatal("Gagal menginisialisasi aplikasi: %v", err)
	}

	logger := app.Logger

	port := config.GetEnv("APP_PORT", "8080")

	err = app.App.Listen("0.0.0.0:" + port)
	if err != nil {
		logger.Error("Gagal menjalankan server: ", err)
	}

	logger.Info("🚀 Server  berjalan pada port ", port)
}
