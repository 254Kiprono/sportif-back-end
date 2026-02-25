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
	return r.db.Exec("INSERT INTO news (id, created_at, updated_at, title, content, image_url, author_id, published, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		news.ID, news.CreatedAt, news.UpdatedAt, news.Title, news.Content, news.ImageURL, news.AuthorID, news.Published, news.Status).Error
}

func (r *newsRepository) GetAll(publishedOnly bool) ([]models.News, error) {
	var news []models.News
	query := "SELECT * FROM news WHERE deleted_at IS NULL"
	if publishedOnly {
		query += " AND published = true"
	}
	query += " ORDER BY created_at DESC"
	err := r.db.Raw(query).Scan(&news).Error
	if err == nil {
		for i := range news {
			r.db.Raw("SELECT id, full_name, username, email FROM users WHERE id = ? AND deleted_at IS NULL", news[i].AuthorID).Scan(&news[i].Author)
		}
	}
	return news, err
}

func (r *newsRepository) GetByID(id string) (*models.News, error) {
	var news models.News
	err := r.db.Raw("SELECT * FROM news WHERE id = ? AND deleted_at IS NULL LIMIT 1", id).Scan(&news).Error
	if err == nil {
		r.db.Raw("SELECT id, full_name, username, email FROM users WHERE id = ? AND deleted_at IS NULL", news.AuthorID).Scan(&news.Author)
	}
	return &news, err
}

func (r *newsRepository) Update(news *models.News) error {
	return r.db.Exec("UPDATE news SET title = ?, content = ?, image_url = ?, published = ?, status = ?, updated_at = NOW() WHERE id = ? AND deleted_at IS NULL",
		news.Title, news.Content, news.ImageURL, news.Published, news.Status, news.ID).Error
}

func (r *newsRepository) Delete(id string) error {
	return r.db.Exec("UPDATE news SET deleted_at = NOW() WHERE id = ?", id).Error
}
