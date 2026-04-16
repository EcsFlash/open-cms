package repos

import (
	"headless-cms/internal/models"

	"gorm.io/gorm"
)

type SectionRepo struct {
	db *gorm.DB
}

func NewSectionRepo(db *gorm.DB) *SectionRepo {
	return &SectionRepo{db: db}
}

type ISectionRepo interface {
	Create(section *models.Section) error
	Update(section *models.Section) error
	Delete(id uint) error
	GetByID(id uint) (*models.Section, error)
	ListAll() ([]models.Section, error)
	ListChildren(parentID uint) ([]models.Section, error)

	// lightweight graph DTO
	ListGraph() ([]SectionGraphRow, error)
}

type SectionGraphRow struct {
	ID       uint
	ParentID *uint
}

func (r *SectionRepo) Create(section *models.Section) error {
	return r.db.Create(section).Error
}

func (r *SectionRepo) Update(section *models.Section) error {
	return r.db.Model(&models.Section{}).Where("id = ?", section.ID).Updates(map[string]any{
		"name":      section.Name,
		"parent_id": section.ParentID,
	}).Error
}

func (r *SectionRepo) Delete(id uint) error {
	return r.db.Delete(&models.Section{}, id).Error
}

func (r *SectionRepo) GetByID(id uint) (*models.Section, error) {
	var s models.Section
	if err := r.db.First(&s, id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SectionRepo) ListAll() ([]models.Section, error) {
	var res []models.Section
	if err := r.db.Order("id asc").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (r *SectionRepo) ListChildren(parentID uint) ([]models.Section, error) {
	var res []models.Section
	if err := r.db.Order("id asc").Find(&res, "parent_id = ?", parentID).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (r *SectionRepo) ListGraph() ([]SectionGraphRow, error) {
	var rows []SectionGraphRow
	if err := r.db.Model(&models.Section{}).Select("id, parent_id").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

