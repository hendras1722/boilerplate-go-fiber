package cronjob

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

// SetupCronJobs menginisialisasi dan mendaftarkan semua cron job.
func SetupCronJobs(db *gorm.DB) *cron.Cron {
	// Memuat instance cron baru
	// Kita bisa menggunakan opsi WithLocation agar sinkron dengan timezone lokal (misal Asia/Jakarta)
	jakartaTime, _ := time.LoadLocation("Asia/Jakarta")
	
	c := cron.New(cron.WithLocation(jakartaTime))

	// ==========================================
	// DAFTAR JOB BISA DITAMBAHKAN DI BAWAH INI
	// ==========================================

	// Contoh 1: Berjalan setiap 1 menit
	c.AddFunc("*/1 * * * *", func() {
		log.Println("[CRON] Menjalankan task setiap 1 menit...")
		// Anda bisa memanggil repository/service di sini
		// misal: userService.CleanUpExpiredTokens()
	})

	// Contoh 2: Berjalan setiap jam 00:00 (tengah malam)
	c.AddFunc("0 0 * * *", func() {
		log.Println("[CRON] Menjalankan task harian di tengah malam...")
	})

	return c
}
