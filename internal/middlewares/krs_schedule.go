package middlewares

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
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

func (m *Middleware) KRSScheduleMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		periodeTahunAkademik, err := m.TaService.GetActiveTahunAkademik()

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Println("Peringatan: Akses KRS ditolak, tidak ada periode akademik yang aktif.")
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"status":     fiber.StatusForbidden,
					"error_code": "KRS_PERIOD_NOT_ACTIVE",
					"errors":     "Saat ini tidak ada periode pengisian KRS yang aktif.",
				})
			}

			log.Printf("ERROR: Gagal mengambil jadwal KRS dari database: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"status": http.StatusInternalServerError,
				"errors": "Gagal memproses permintaan Anda karena kesalahan server.",
			})
		}

		now := time.Now().In(wibLocation)

		jamBuka, err1 := time.Parse("15:04:05", periodeTahunAkademik.JamMulaiHarianKRS)
		jamTutup, err2 := time.Parse("15:04:05", periodeTahunAkademik.JamSelesaiHarianKRS)

		if err1 != nil || err2 != nil {
			log.Printf("ERROR: Format jam di database tidak valid. Jam Buka: %s, Jam Tutup: %s", periodeTahunAkademik.JamMulaiHarianKRS, periodeTahunAkademik.JamSelesaiHarianKRS)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"status": http.StatusInternalServerError,
				"errors": "Konfigurasi jam KRS di server tidak valid.",
			})
		}

		waktuBuka := time.Date(now.Year(), now.Month(), now.Day(), jamBuka.Hour(), jamBuka.Minute(), 0, 0, wibLocation)
		waktuTutup := time.Date(now.Year(), now.Month(), now.Day(), jamTutup.Hour(), jamTutup.Minute(), 0, 0, wibLocation)

		tanggalSelesaiKRS := periodeTahunAkademik.TanggalSelesaiKRS.Add(24*time.Hour - 1*time.Second)

		isDiluarJadwal := now.Before(periodeTahunAkademik.TanggalMulaiKRS) || now.After(tanggalSelesaiKRS) || now.Before(waktuBuka) || now.After(waktuTutup)

		if isDiluarJadwal {
			errorMessage := fmt.Sprintf(
				"Pendaftaran KRS hanya bisa dilakukan dari %s hingga %s, setiap hari pukul %s - %s WIB.",
				periodeTahunAkademik.TanggalMulaiKRS.Format("2 January 2006"),
				periodeTahunAkademik.TanggalSelesaiKRS.Format("2 January 2006"),
				waktuBuka.Format("15:04"),
				waktuTutup.Format("15:04"),
			)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":     fiber.StatusForbidden,
				"error_code": "KRS_SCHEDULE_NOT_ALLOWED",
				"errors":     errorMessage,
			})
		}

		return c.Next()
	}
}
