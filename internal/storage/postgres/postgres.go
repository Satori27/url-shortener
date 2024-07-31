package psq

import (
	"fmt"
	"url-shortener/internal/storage"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"url-shortener/internal/config"
)

type Storage struct {
	db *gorm.DB
}

type URL struct {
	ID    uint64 `gorm:"primaryKey"`
	Alias string `gorm:"unique"`
	URL   string
}


func New(cfg config.Config) (*Storage, error) {

	const op = "storage.postgres.New"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=prefer", cfg.PGHost, cfg.PGUser, cfg.PGPassword, cfg.PGName, cfg.PGPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}, &gorm.Config{PrepareStmt: true})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	err = db.AutoMigrate(&URL{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (uint64, error) {
	const op = "storage.postgres.SaveURL"
	urlRow := URL{
		Alias: alias,
		URL:   urlToSave,
	}
	result := s.db.Create(&urlRow)

	if result.Error != nil {
		if gorm.ErrDuplicatedKey == result.Error {
			return 0, storage.ErrURLExists
		}
		return 0, fmt.Errorf("%s: %w", op, result.Error)
	}
	return urlRow.ID, nil

}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgres.GetURL"

	var urlRow URL
	queryResult := s.db.Model(&URL{}).Select("url").Where("alias=?", alias).First(&urlRow)

	if queryResult.Error != nil {
		if queryResult.Error == gorm.ErrRecordNotFound {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("%s: %w", op, queryResult.Error)
	}

	return urlRow.URL, nil

}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.postgres.DeleteURL"

	queryResult := s.db.Model(&URL{}).Where("alias=?", alias).First(&URL{})
	queryResult.Delete(&URL{})

	if queryResult.Error != nil {
		if queryResult.Error == gorm.ErrRecordNotFound {
			return storage.ErrURLNotFound
		}
		return fmt.Errorf("%s: %w", op, queryResult.Error)
	}
	return nil
}

func (s *Storage) AliasExist(alias string) error {
	var url URL
	queryResult := s.db.Select("alias").Where("alias=?", alias).First(&url)

	if queryResult.Error == nil {
		return storage.ErrURLExists
	}
	if queryResult.Error == gorm.ErrRecordNotFound {
		return storage.ErrURLNotFound
	}
	return queryResult.Error
}
