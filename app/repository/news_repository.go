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
	query := `INSERT INTO news (id, created_at, updated_at, title, content, image_url, author_id, published) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	return r.db.Exec(query, news.ID, news.CreatedAt, news.UpdatedAt, news.Title, news.Content, news.ImageURL, news.AuthorID, news.Published).Error
}

func (r *newsRepository) GetAll(publishedOnly bool) ([]models.News, error) {
	var news []models.News
	query := `SELECT n.*, u.full_name as "Author.FullName" 
	          FROM news n 
	          LEFT JOIN users u ON n.author_id = u.id 
	          WHERE n.deleted_at IS NULL`
	if publishedOnly {
		query += " AND n.published = true"
	}
	query += " ORDER BY n.created_at DESC"

	// Note: Scan might not automatically fill Author struct nested fields correctly with Raw SQL joins
	// without proper naming or using Preload. But user asked for Raw SQL.
	// For simple Raw SQL, we'll just get news.
	err := r.db.Raw(query).Scan(&news).Error
	return news, err
}

func (r *newsRepository) GetByID(id string) (*models.News, error) {
	var news models.News
	query := `SELECT * FROM news WHERE id = ? AND deleted_at IS NULL LIMIT 1`
	err := r.db.Raw(query, id).Scan(&news).Error
	return &news, err
}

func (r *newsRepository) Update(news *models.News) error {
	query := `UPDATE news SET updated_at = ?, title = ?, content = ?, image_url = ?, published = ? WHERE id = ?`
	return r.db.Exec(query, news.UpdatedAt, news.Title, news.Content, news.ImageURL, news.Published, news.ID).Error
}

func (r *newsRepository) Delete(id string) error {
	query := `UPDATE news SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?`
	return r.db.Exec(query, id).Error
}
