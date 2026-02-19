package services

import (
	"webuye-sportif/app/models"
	"webuye-sportif/app/repository"

	"github.com/google/uuid"
)

type NewsService interface {
	CreateNews(news *models.News, authorID string) error
	GetNews(publishedOnly bool) ([]models.News, error)
	GetNewsByID(id string) (*models.News, error)
	UpdateNews(id string, news *models.News) error
	DeleteNews(id string) error
}

type newsService struct {
	repo repository.NewsRepository
}

func NewNewsService(repo repository.NewsRepository) NewsService {
	return &newsService{repo}
}

func (s *newsService) CreateNews(news *models.News, authorID string) error {
	uID, err := uuid.Parse(authorID)
	if err != nil {
		return err
	}
	news.AuthorID = uID
	return s.repo.Create(news)
}

func (s *newsService) GetNews(publishedOnly bool) ([]models.News, error) {
	return s.repo.GetAll(publishedOnly)
}

func (s *newsService) GetNewsByID(id string) (*models.News, error) {
	return s.repo.GetByID(id)
}

func (s *newsService) UpdateNews(id string, news *models.News) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	news.ID = existing.ID
	news.AuthorID = existing.AuthorID
	return s.repo.Update(news)
}

func (s *newsService) DeleteNews(id string) error {
	return s.repo.Delete(id)
}
