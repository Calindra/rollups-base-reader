package repository

import (
	"github.com/jmoiron/sqlx"
)

type EpochRepository struct {
	Db *sqlx.DB
}

func NewEpochRepository(db *sqlx.DB) *EpochRepository {
	return &EpochRepository{db}
}

