package service

import (
	"fmt"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/models/api"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/repository"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/utils"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type MahasiswaService interface {
	GetSyaratPengisianKRS(nim string) (result api.SyaratPengisisanKrsResponse, err error)
	GetInformasiMahasiswa(nim string) (result api.InformasiMahasiswaResponse, err error)
	GetIdKurikulumMahasiswa(nim string) (idKurikulum string, err error)
}

type MahasiswaServiceImpl struct {
	Logger               *logrus.Logger
	MahasiswaRepository  repository.MahasiswaRepository
	TahunAkademikService TahunAkademikService
}

func NewMahasiswaService(
	mhsRepo repository.MahasiswaRepository,
	taService TahunAkademikService,
	logger *logrus.Logger,
) MahasiswaService {
	return &MahasiswaServiceImpl{
		MahasiswaRepository:  mhsRepo,
		TahunAkademikService: taService,
		Logger:               logger,
	}
}

func (s *MahasiswaServiceImpl) GetInformasiMahasiswa(nim string) (result api.InformasiMahasiswaResponse, err error) {
	mahasiswa, err := s.MahasiswaRepository.FindByNIMWithTotalSKS(nim)
	if err != nil {
		s.Logger.WithError(err).Error("Gagal menemukan mahasiswa dengan NIM: " + nim)
		return api.InformasiMahasiswaResponse{}, err
	}
	// Atau jika ingin menambahkan lebih dari satu field:
	s.Logger.WithFields(logrus.Fields{
		"nim":       nim,
		"mahasiswa": mahasiswa,
	}).Info("Detail data mahasiswa berhasil diambil")

	periodeAkademik, err := s.TahunAkademikService.GetActiveTahunAkademik()

	if err != nil {
		s.Logger.WithError(err).Error("Gagal mendapatkan tahun akademik aktif")
		return api.InformasiMahasiswaResponse{}, err
	}

	return api.InformasiMahasiswaResponse{
		TahunAkademik: periodeAkademik.TahunAkademik,
		Semester:      periodeAkademik.JenisSemester.String(),
		IPK:           mahasiswa.IPK,
		SksKumulatif:  mahasiswa.SKSKumulatif,
		IpsLalu:       mahasiswa.IPSLalu,
		JatahSks:      mahasiswa.JatahSKS,
		SksAmbil:      mahasiswa.TotalSKSDiambil,
		SisaSks:       mahasiswa.JatahSKS - mahasiswa.TotalSKSDiambil,
	}, nil
}

func (s *MahasiswaServiceImpl) GetSyaratPengisianKRS(nim string) (result api.SyaratPengisisanKrsResponse, err error) {
	var dataSyarat []api.SyaratItem
	periodeAktif, err := s.TahunAkademikService.GetActiveTahunAkademik()

	if err != nil {
		s.Logger.WithError(err).Error("Gagal mendapatkan tahun akademik aktif")
		return api.SyaratPengisisanKrsResponse{}, err
	}

	syaratPembayaranText := fmt.Sprintf("Bayar Biaya Pendidikan %s Tahun Akademik %s = Sudah Bayar", utils.FirstToUpper(strings.ToLower(periodeAktif.JenisSemester.String())), periodeAktif.TahunAkademik)

	mahasiswa, err := s.MahasiswaRepository.FindByNIM(nim)

	if err != nil {
		s.Logger.WithError(err).Error("Gagal menemukan mahasiswa dengan NIM: " + nim)
		return api.SyaratPengisisanKrsResponse{}, err
	}

	dataSyarat = append(dataSyarat, api.SyaratItem{
		Syarat: syaratPembayaranText,
		Isi:    mahasiswa.StatusPembayaran.String(),
		Status: mahasiswa.StatusPembayaran.String() == "Sudah Bayar",
	})

	semesterStatus := mahasiswa.SemesterBerjalan >= 3 && mahasiswa.SemesterBerjalan <= 14

	dataSyarat = append(dataSyarat, api.SyaratItem{
		Syarat: "Semester Mahasiswa = 3|4|5|6|7|8|9|10|11|12|13|14",
		Isi:    strconv.Itoa(mahasiswa.SemesterBerjalan),
		Status: semesterStatus,
	})

	dataSyarat = append(dataSyarat, api.SyaratItem{
		Syarat: "Status Mahasiswa = Aktif",
		Isi:    mahasiswa.StatusMahasiswa.String(),
		Status: mahasiswa.StatusMahasiswa.String() == "Aktif",
	})

	semuaSyaratTerpenuhi := true
	for _, syarat := range dataSyarat {
		if !syarat.Status {
			semuaSyaratTerpenuhi = false
			break
		}
	}
	return api.SyaratPengisisanKrsResponse{
		Judul:                "Syarat Pengisian",
		DataSyarat:           dataSyarat,
		PengisisanKrsEnabled: semuaSyaratTerpenuhi,
	}, nil
}

func (s *MahasiswaServiceImpl) GetIdKurikulumMahasiswa(nim string) (idKurikulum string, err error) {
	mahasiswa, err := s.MahasiswaRepository.FindByNIM(nim)
	if err != nil {
		s.Logger.WithError(err).Error("Gagal menemukan mahasiswa dengan NIM: " + nim)
		return "", err
	}

	if mahasiswa.IDKurikulum.String() == "" {
		s.Logger.Error("Kurikulum tidak ditemukan untuk mahasiswa dengan NIM: " + nim)
		return "", fmt.Errorf("kurikulum tidak ditemukan untuk mahasiswa dengan NIM: %s", nim)
	}

	return mahasiswa.IDKurikulum.String(), nil
}
