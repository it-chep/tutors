package get_tutor_finance

import (
	"context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/get_tutor_finance/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/internal/pkg/convert"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Action struct {
	dal *dal.Repository
}

func New(pool *pgxpool.Pool) *Action {
	return &Action{
		dal: dal.NewRepository(pool),
	}
}

func (a *Action) Do(ctx context.Context, tutorID int64, from, to string) (dto.TutorFinance, error) {
	fromTime, toTime, err := convert.StringsIntervalToTime(from, to)
	if err != nil {
		return dto.TutorFinance{}, err
	}

	lessons, err := a.dal.GetLessonsCounters(ctx, tutorID, fromTime, toTime)
	if err != nil {
		return dto.TutorFinance{}, err
	}

	conversion, err := a.dal.GetConversion(ctx, tutorID, fromTime, toTime)
	if err != nil {
		return dto.TutorFinance{}, err
	}
	amount, err := a.dal.GetFinanceInfo(ctx, tutorID, fromTime, toTime)
	if err != nil {
		return dto.TutorFinance{}, err
	}

	return dto.TutorFinance{
		Conversion: conversion,
		Count:      lessons.LessonsCount,
		BaseCount:  lessons.BaseCount,
		TrialCount: lessons.TrialCount,
		Amount:     amount,
	}, nil
}
