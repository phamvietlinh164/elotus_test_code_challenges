package user

import (
	"errors"
	"time"

	"hackathon/internal/config"
	"hackathon/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(username, password string) error
	Login(username, password string) (string, error)
	RevokeTokens(userID uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(username, password string) error {
	if _, err := s.repo.FindByUsername(username); err == nil {
		return errors.New("username already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &User{
		Username: username,
		Password: string(hash),
	}

	return s.repo.Create(user)
}

func (s *service) Login(username, password string) (string, error) {
	u, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	duration := time.Duration(config.Cfg.Jwt.Ttl) * time.Hour
	token, err := utils.GenerateToken(u.ID, u.IsAdmin, duration)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) RevokeTokens(userID uint) error {
	return s.repo.UpdateTokenRevokedAt(userID, time.Now().UTC())
}
