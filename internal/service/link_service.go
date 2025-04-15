package service

import (
	"crypto/rand"
	"errors"
	"time"

	"github.com/RamanDudoits/shortLink-go/internal/repository"
)

type LinkServiceInterface interface {
	UpdateLink(id, userID int, updates map[string]interface{}) (*repository.Link, error)
	GetLink(id, userID int) (*repository.Link, error)
	Create(originalURL string, userID int) (*repository.Link, error)
	GetUserLinks(userID int) ([]*repository.Link, error)
	DeleteLink(id, userID int) error
}

type LinkServiceRedirectInterface interface {
	GetByShortCode(shortLink string) (*repository.Link, error)
	IncrementClickCount(id int, clickCount int) error
}

type LinkService struct {
	linkRepo repository.LinkRepository
}

func NewLinkService(linkRepo repository.LinkRepository) *LinkService {
	return &LinkService{linkRepo: linkRepo}
}

func (s *LinkService) Create(originalURL string, userID int) (*repository.Link, error) {
	if existing, _ := s.linkRepo.FindByURLAndUser(originalURL, userID); existing != nil {
		return existing, nil
	}

	shortCode, err := generateShortCode(originalURL)
	if err != nil {
		return nil, err
	}
	
	link := &repository.Link{
		OriginalURL: originalURL,
		ShortCode:   shortCode,
		UserID:      userID,
		ClickCount:  0,
		CreatedAt:   time.Now(),
	}
	
	return s.linkRepo.Create(link)
}

func (s *LinkService) GetUserLinks(userID int) ([]*repository.Link, error) {
	return s.linkRepo.FindByUserID(userID)
}

func (s *LinkService) GetLink(id, userID int) (*repository.Link, error) {
	link, err := s.linkRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if link.UserID != userID {
		return nil, errors.New("access denied")
	}

	return link, nil
}

func (s *LinkService) UpdateLink(id, userID int, updates map[string]interface{}) (*repository.Link, error) {
	if _, err := s.GetLink(id, userID); err != nil {
		return nil, err
	}

	return s.linkRepo.Update(id, updates)
}

func (s *LinkService) DeleteLink(id, userID int) error {
	return s.linkRepo.Delete(id, userID)
}

func (s *LinkService) IncrementClickCount(id int, clickCount int) error {
_, err := s.linkRepo.Update(id, map[string]interface{}{
        "clicks": clickCount + 1,
    })
    return err
}

func generateShortCode(url string) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 5
    result := make([]byte, length)

    randomBytes := make([]byte, length)
    if _, err := rand.Read(randomBytes); err != nil {
        return "", err
    }

    for i := 0; i < length; i++ {
        result[i] = charset[int(randomBytes[i])%len(charset)]
    }

    return string(result), nil
}

func (s *LinkService) GetByShortCode(shortLink string) (*repository.Link, error) {
	link, err := s.linkRepo.Find(map[string]interface{}{
		"short_link": shortLink,
	})
	if err != nil {
		return nil, err
	}
	return link, nil
}