package service

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/models"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/repository"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Username string `validate:"required" json:"username"`
	Password string `validate:"required" json:"password"`
}

type AuthService interface {
	Login(ctx context.Context, req LoginRequest) (sessionID string, err error)
	Logout(ctx context.Context, sessionID string) error
}

type AuthServiceImpl struct {
	Logger              *logrus.Logger
	AuthRepository      repository.AuthRepository
	SessionRepository   repository.SessionRepository
	MahasiswaRepository repository.MahasiswaRepository
	DosenRepository     repository.DosenRepository
	PegawaiRepository   repository.PegawaiRepository
}

func NewAuthService(
	authRepo repository.AuthRepository,
	sessionRepo repository.SessionRepository,
	logger *logrus.Logger,
	mhsRepo repository.MahasiswaRepository,
	dosenRepo repository.DosenRepository,
	pegawaiRepository repository.PegawaiRepository,

) AuthService {
	return &AuthServiceImpl{
		AuthRepository:      authRepo,
		SessionRepository:   sessionRepo,
		Logger:              logger,
		MahasiswaRepository: mhsRepo,
		DosenRepository:     dosenRepo,
		PegawaiRepository:   pegawaiRepository,
	}
}

var TTL = 2 * time.Hour
var (
	ErrInvalidCredentials = errors.New("kombinasi NIM dan password salah")
	ErrInternalServer     = errors.New("terjadi kesalahan internal pada server")
)

func (s *AuthServiceImpl) Login(ctx context.Context, req LoginRequest) (string, error) {

	log := s.Logger.WithField("nim", req.Username)
	log.Info("Memproses permintaan login")

	// 2. Cari user via repository
	user, err := s.AuthRepository.FindByCredential(req.Username)
	if err != nil {
		// Jika user tidak ditemukan, kembalikan error umum
		log.Warn("Percobaan login gagal: NIM tidak ditemukan")
		return "", ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {

		return "", ErrInvalidCredentials
	}

	var nama, nomorInduk string

	log = s.Logger.WithFields(logrus.Fields{
		"user_id":   user.IDUser,
		"role_name": user.Role.RoleName,
	})
	log.Info("Otentikasi berhasil, memulai pencarian profil dinamis")

	const (
		RoleMahasiswa = "MAHASISWA"
		RoleDosen     = "DOSEN"
		RolePegawai   = "PEGAWAI"
	)

	switch user.Role.RoleName {
	case RoleMahasiswa:
		log.Info("Role terdeteksi sebagai Mahasiswa, menjalankan query ke tabel mahasiswa...")

		profile, err := s.MahasiswaRepository.FindByUserID(user.IDUser)
		if err != nil {
			log.WithError(err).Error("QUERY PROFIL MAHASISWA GAGAL")
		} else {
			// Jika berhasil, log data yang ditemukan
			log.WithField("profile_found", profile).Info("Profil Mahasiswa berhasil ditemukan")
			nama = profile.Nama
			nomorInduk = profile.NIM
		}

	case RoleDosen:
		log.Info("Role terdeteksi sebagai Dosen, menjalankan query ke tabel dosen...")

		profile, err := s.DosenRepository.FindByUserID(user.IDUser)
		if err != nil {
			log.WithError(err).Error("QUERY PROFIL DOSEN GAGAL")
		} else {
			log.WithField("profile_found", profile).Info("Profil Dosen berhasil ditemukan")
			nama = profile.Nama
			nomorInduk = profile.NIP
		}

	case RolePegawai:
		log.Info("Role terdeteksi sebagai Pegawai, menjalankan query ke tabel pegawai...")

		profile, err := s.PegawaiRepository.FindByUserID(user.IDUser)
		if err != nil {
			log.WithError(err).Error("QUERY PROFIL PEGAWAI GAGAL")
		} else {
			log.WithField("profile_found", profile).Info("Profil Pegawai berhasil ditemukan")
			nama = profile.Nama
			nomorInduk = profile.NIP
		}

	default:
		log.Info("Role tidak memiliki profil spesifik yang perlu diambil.")
	}

	sessionPayload := &models.Session{
		UserId:     user.IDUser.String(),
		NomorInduk: nomorInduk,
		Nama:       nama,
		Role: models.RoleType{
			IDRole:   user.IDRole.String(),
			RoleName: user.Role.RoleName,
		},
	}
	sessionID := uuid.NewString()

	err = s.SessionRepository.Create(ctx, sessionID, sessionPayload, TTL)
	if err != nil {
		log.WithError(err).Error("Gagal membuat sesi setelah otentikasi berhasil")
		return "", ErrInternalServer
	}

	log.WithField("session_id", sessionID).Info("Login berhasil, sesi dibuat")
	return sessionID, nil
}

func (s *AuthServiceImpl) Logout(ctx context.Context, sessionID string) error {
	s.Logger.WithField("session_id", sessionID).Info("Memproses permintaan logout")
	err := s.SessionRepository.Delete(ctx, sessionID)
	if err != nil {
		s.Logger.WithError(err).Error("Gagal menghapus sesi saat logout")
		return ErrInternalServer
	}
	return nil
}
