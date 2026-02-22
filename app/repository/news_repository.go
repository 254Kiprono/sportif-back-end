package repository

import (
	"webuye-sportif/app/models"

	"gorm.io/gorm"
)

type NewsRepository interface {
	Create(news *models.News) error
	GetAll(publishedOnly bool) ([]models.News, error)
	GetByID(id string) (*models.News, error)
	Update(news *models.News) error
	Delete(id string) error
}

type newsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(db *gorm.DB) NewsRepository {
	return &newsRepository{db}
}

func (r *newsRepository) Create(news *models.News) error {
	return r.db.Create(news).Error
}

func (r *newsRepository) GetAll(publishedOnly bool) ([]models.News, error) {
	var news []models.News
	db := r.db.Preload("Author").Order("created_at DESC")
	if publishedOnly {
		db = db.Where("published = ?", true)
	}
	err := db.Find(&news).Error
	return news, err
}

func (r *newsRepository) GetByID(id string) (*models.News, error) {
	var news models.News
	err := r.db.Preload("Author").First(&news, "id = ?", id).Error
	return &news, err
}

func (r *newsRepository) Update(news *models.News) error {
	return r.db.Save(news).Error
}

func (r *newsRepository) Delete(id string) error {
	return r.db.Delete(&models.News{}, "id = ?", id).Error
}
