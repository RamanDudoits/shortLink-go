package repository

import "time"

type LinkRepository interface {
    Create(link *Link) (*Link, error)
    FindByID(id int) (*Link, error)
    FindByUserID(userID int) ([]*Link, error)
    FindByURLAndUser(url string, userID int) (*Link, error)
    Update(id int, updates map[string]interface{}) (*Link, error)
    Delete(id, userID int) error
    Find(filter map[string]interface{}) (*Link, error)
}

type Link struct {
    ID          int       `json:"id"`
    OriginalURL string    `json:"original_url" db:"link"`
    ShortCode   string    `json:"short_code" db:"short_link"`
    UserID      int       `json:"user_id" db:"user_id"`
    ClickCount  int       `json:"click_count" db:"clicks"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
}