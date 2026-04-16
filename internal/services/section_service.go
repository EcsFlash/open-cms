package services

import (
	"errors"
	"fmt"
	"headless-cms/internal/models"
	"headless-cms/internal/repos"
)

var (
	ErrSectionCycle = errors.New("section parent change would cause a cycle")
)

type SectionService struct {
	repo repos.ISectionRepo
}

type ISectionService interface {
	Create(section *models.Section) error
	Update(section *models.Section) error
	Delete(id uint) error
	GetByID(id uint) (*models.Section, error)
	ListAll() ([]models.Section, error)
	ListChildren(parentID uint) ([]models.Section, error)
}

func NewSectionService(repo repos.ISectionRepo) *SectionService {
	return &SectionService{repo: repo}
}

func (s *SectionService) Create(section *models.Section) error {
	if section.ParentID != nil {
		if *section.ParentID == section.ID && section.ID != 0 {
			return ErrSectionCycle
		}
		// For create, ID is typically unknown yet; still guard obvious self-parent (if provided).
	}
	return s.repo.Create(section)
}

func (s *SectionService) Update(section *models.Section) error {
	old, err := s.repo.GetByID(section.ID)
	if err != nil {
		return err
	}

	// Only validate cycles if parent changes.
	if !uintPtrEqual(old.ParentID, section.ParentID) {
		if s.wouldCauseCycle(section.ID, section.ParentID) {
			return fmt.Errorf("%w", ErrSectionCycle)
		}
	}

	return s.repo.Update(section)
}

func (s *SectionService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *SectionService) GetByID(id uint) (*models.Section, error) {
	return s.repo.GetByID(id)
}

func (s *SectionService) ListAll() ([]models.Section, error) {
	return s.repo.ListAll()
}

func (s *SectionService) ListChildren(parentID uint) ([]models.Section, error) {
	return s.repo.ListChildren(parentID)
}

func (s *SectionService) wouldCauseCycle(nodeID uint, newParentID *uint) bool {
	if newParentID == nil {
		return false
	}
	if *newParentID == nodeID {
		return true
	}

	rows, err := s.repo.ListGraph()
	if err != nil {
		// fail closed
		return true
	}
	parentByID := make(map[uint]*uint, len(rows))
	for _, r := range rows {
		parentByID[r.ID] = r.ParentID
	}

	current := *newParentID
	visited := make(map[uint]bool)
	for {
		if visited[current] {
			return true
		}
		visited[current] = true

		if current == nodeID {
			return true
		}

		p, ok := parentByID[current]
		if !ok {
			// referencing non-existent parent -> treat as invalid
			return true
		}
		if p == nil {
			return false
		}
		current = *p
	}
}

func uintPtrEqual(a, b *uint) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

