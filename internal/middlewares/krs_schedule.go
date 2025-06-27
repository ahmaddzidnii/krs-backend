package middlewares

import (
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/config"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/utils"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"time"
)

var wibLocation *time.Location

func InitTimezoneWib() {
	var err error
	wibLocation, err = time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatalf("FATAL: Gagal memuat zona waktu 'Asia/Jakarta': %v", err)
	}
	log.Println("Zona waktu 'Asia/Jakarta' berhasil dimuat.")
}

func KRSScheduleMiddleware(c *fiber.Ctx) error {
	tanggalMulaiStr := config.GetEnv("KRS_START_DATE", "2025-06-26")
	tanggalSelesaiStr := config.GetEnv("KRS_END_DATE", "2025-06-28")
	jamBuka := config.GetEnvAsInt("KRS_OPEN_HOUR", 8)
	jamTutup := config.GetEnvAsInt("KRS_CLOSE_HOUR", 15)

	tanggalMulai, err1 := time.ParseInLocation("2006-01-02", tanggalMulaiStr, wibLocation)
	tanggalSelesai, err2 := time.ParseInLocation("2006-01-02", tanggalSelesaiStr, wibLocation)
	if err1 != nil || err2 != nil {
		log.Printf("ERROR: Format tanggal di ENV tidak valid")
		return utils.Error(c, http.StatusInternalServerError, "Konfigurasi server tidak valid")
	}

	now := time.Now().In(wibLocation)

	waktuBuka := time.Date(now.Year(), now.Month(), now.Day(), jamBuka, 0, 0, 0, wibLocation)
	waktuTutup := time.Date(now.Year(), now.Month(), now.Day(), jamTutup, 0, 0, 0, wibLocation)

	isDiluarJadwal := now.Before(tanggalMulai) || now.After(tanggalSelesai.Add(23*time.Hour+59*time.Minute+59*time.Second)) || now.Before(waktuBuka) || now.After(waktuTutup)

	if isDiluarJadwal {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":     fiber.StatusForbidden,
			"error_code": "KRS_SCHEDULE_NOT_ALLOWED",
			"errors":     "Pendaftaran KRS hanya dapat dilakukan sesuai jadwal yang telah ditentukan.",
		})
	}

	//log.Printf("Akses diizinkan pada %s", now.Format(time.RFC1123))
	return c.Next()
}
