package repository

import "flashcard/internal/models"

type CardRepository interface {
	GetAll() ([]models.Card, error)
	GetWithWrongs() ([]models.Card, error)
	SaveAll(cards []models.Card) error
	DeleteById(cardId int) error
	Add(card models.Card) error
	GetIdLastInsert() (int, error)
}
