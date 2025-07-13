package service

import (
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/models/api"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/models/domain"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/repository"
	"strconv"
)

type StatusKoutaDetail struct {
	Terisi   int  `json:"terisi"`
	Kouta    int  `json:"kouta"`
	IsFull   bool `json:"is_full"`
	IsJoined bool `json:"is_joined"`
}

type PenjadwalanService interface {
	GetPenawaranKelasByNim(nim string) (*api.PenawaranPerSemesterResponse, error)
	GetStatusKoutaKelasByIdKelas(idKelas string, sessionData *domain.Session) (*api.StatusKoutaKelasResponse, error)

	GetStatusKoutaKelasInBatch(idKelas []string, sessionData *domain.Session) (map[string]*StatusKoutaDetail, error)
}

type PenjadwalanServiceImpl struct {
	PenjadwalanRepository repository.PenjadwalanRepository
	MahasiswaRepository   repository.MahasiswaRepository
	MahasiswaService      MahasiswaService
}

func NewPenjadwalanService(penjadwalanRepository repository.PenjadwalanRepository, mahasiswaRepository repository.MahasiswaRepository, mhsservice MahasiswaService) PenjadwalanService {
	return &PenjadwalanServiceImpl{PenjadwalanRepository: penjadwalanRepository,
		MahasiswaRepository: mahasiswaRepository,
		MahasiswaService:    mhsservice,
	}
}

func (s *PenjadwalanServiceImpl) GetPenawaranKelasByNim(nim string) (*api.PenawaranPerSemesterResponse, error) {
	idKurikulum, err := s.MahasiswaService.GetIdKurikulumMahasiswa(nim)
	if err != nil {
		return nil, err
	}
	semesterBerjalan := 4
	//idKurikulum := "cc4041be-bd4f-498a-a4f5-c8ce1d6786cc"

	data, err := s.PenjadwalanRepository.GetJadwalKelasDitawarkanBySemesterAndIdKurikulum(semesterBerjalan, idKurikulum)
	if err != nil {
		return nil, err
	}

	// 1. Buat map untuk menampung hasil pengelompokan.
	//    Kuncinya adalah string semester (misal "4"), nilainya adalah slice dari detail kelas.
	groupedKelas := make(map[string][]*api.DaftarPenawaranKelasResponse)

	for _, kelas := range data {
		// Lewati iterasi jika tidak ada detail kurikulum (data tidak valid)
		if len(kelas.MataKuliah.DetailKurikulum) == 0 {
			continue
		}

		detail := kelas.MataKuliah.DetailKurikulum[0]

		// 2. Buat objek detail kelas
		responseItem := &api.DaftarPenawaranKelasResponse{
			IdKelas:         kelas.IDKelas.String(),
			KodeMataKuliah:  kelas.MataKuliah.KodeMatkul,
			NamaMataKuliah:  kelas.MataKuliah.Nama,
			Sks:             kelas.MataKuliah.SKS,
			NamaKelas:       kelas.NamaKelas,
			JenisMataKuliah: string(detail.JenisMatkul),
			SemesterPaket:   detail.SemesterPaket,
		}

		if detail.Kurikulum.KodeKurikulum != "" {
			responseItem.KodeKurikulum = detail.Kurikulum.KodeKurikulum
		}

		// Mapping Dosen
		dosenPengajar := make([]api.Dosen, len(kelas.DosenPengajar))
		for i, dosen := range kelas.DosenPengajar {
			dosenPengajar[i] = api.Dosen{NamaDosen: dosen.Nama, NipDosen: dosen.NIP}
		}
		responseItem.DosenPengajar = dosenPengajar

		// Mapping Jadwal
		jadwal := make([]api.Jadwal, len(kelas.JadwalKelas))
		for i, j := range kelas.JadwalKelas {
			jadwal[i] = api.Jadwal{
				Hari:         j.Hari,
				WaktuMulai:   j.WaktuMulai.Time.Format("15:04"),
				WaktuSelesai: j.WaktuSelesai.Time.Format("15:04"),
				Ruangan:      j.Ruang,
			}
		}
		responseItem.Jadwal = jadwal

		// 3. Masukkan ke dalam map berdasarkan semester_paket
		semesterKey := strconv.Itoa(detail.SemesterPaket)
		groupedKelas[semesterKey] = append(groupedKelas[semesterKey], responseItem)
	}

	// 4. Bungkus map ke dalam struct response akhir
	finalResponse := &api.PenawaranPerSemesterResponse{
		SemesterPaket: groupedKelas,
	}

	return finalResponse, nil
}

func (s *PenjadwalanServiceImpl) GetStatusKoutaKelasByIdKelas(idKelas string, sessionData *domain.Session) (*api.StatusKoutaKelasResponse, error) {

	kelas, err := s.PenjadwalanRepository.GetKelasDitawarkanById(idKelas)

	if err != nil {
		return nil, err
	}

	terisi, err := s.PenjadwalanRepository.GetJumlahTerisiKelasByIdKelas(idKelas)

	if err != nil {
		return nil, err
	}

	isJoined, err := s.PenjadwalanRepository.GetIsJoinedKelasByNimAndIdKelas(idKelas, sessionData.NomorInduk)

	if err != nil {
		return nil, err
	}

	kouta := kelas.Kouta

	isFull := int(terisi) >= kouta

	return &api.StatusKoutaKelasResponse{
		IDKelas:  idKelas,
		Terisi:   int(terisi),
		Kouta:    kouta,
		IsFull:   isFull,
		IsJoined: isJoined,
	}, nil
}

func (s *PenjadwalanServiceImpl) GetStatusKoutaKelasInBatch(idKelas []string, sessionData *domain.Session) (map[string]*StatusKoutaDetail, error) {
	// 1. Ambil semua detail kelas (untuk kuota) dalam 1 query
	kelasDitawarkan, err := s.PenjadwalanRepository.GetKelasDitawarkanByIds(idKelas)
	if err != nil {
		return nil, err
	}

	// 2. Ambil semua jumlah terisi untuk setiap kelas dalam 1 query
	jumlahTerisiMap, err := s.PenjadwalanRepository.GetJumlahTerisiKelasByIds(idKelas)
	if err != nil {
		return nil, err
	}

	// 3. Ambil status 'join' mahasiswa untuk semua kelas dalam 1 query
	isJoinedMap, err := s.PenjadwalanRepository.GetIsJoinedKelasForNim(idKelas, sessionData.NomorInduk)
	if err != nil {
		return nil, err
	}

	// Buat map dari kelas yang ditawarkan untuk pencarian cepat
	kelasMap := make(map[string]*domain.KelasDitawarkan)
	for _, k := range kelasDitawarkan {
		kelasMap[k.IDKelas.String()] = k
	}

	// 4. Gabungkan semua hasil menjadi SEBUAH MAP
	responseMap := make(map[string]*StatusKoutaDetail)

	for _, id := range idKelas {
		kelas, ok := kelasMap[id]
		if !ok {
			// Jika karena suatu alasan ID kelas tidak ditemukan, lewati.
			continue
		}

		terisi := jumlahTerisiMap[id] // Akan menjadi 0 jika tidak ada di map
		isJoined := isJoinedMap[id]   // Akan menjadi false jika tidak ada di map
		kouta := kelas.Kouta
		isFull := int(terisi) >= kouta

		// Masukkan data ke dalam map dengan ID sebagai key
		responseMap[id] = &StatusKoutaDetail{
			Terisi:   int(terisi),
			Kouta:    kouta,
			IsFull:   isFull,
			IsJoined: isJoined,
		}
	}

	return responseMap, nil
}
