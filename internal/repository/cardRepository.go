package repository

import "flashcard/internal/models"

type CardRepository interface {
	GetAll() ([]models.Card, error)
	SaveAll(cards []models.Card) error
}
