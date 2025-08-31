package dao

import (
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/pkg/xo"
)

type SubjectDAO struct {
	xo.Subject
}

type SubjectsDao []SubjectDAO

func (subj *SubjectDAO) ToDomain() dto.Subject {
	return dto.Subject{
		ID:   subj.ID,
		Name: subj.Name,
	}
}

func (subjs SubjectsDao) ToDomain() []dto.Subject {
	domain := make([]dto.Subject, 0, len(subjs))
	for _, subj := range subjs {
		domain = append(domain, subj.ToDomain())
	}
	return domain
}
