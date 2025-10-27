package upload

import "gorm.io/gorm"

type Repository interface {
	Save(meta *Image) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Save(meta *Image) error {
	return r.db.Create(meta).Error
}
