package repository

import (
	"encoding/json"
	"flashcard/internal/models"
	"fmt"
	"io"
	"os"
)

type InJsonRepository struct {
	filePath string
}

func NewInJsonRepository(filePath string) *InJsonRepository {
	return &InJsonRepository{
		filePath: filePath,
	}
}

func (rep *InJsonRepository) GetAll() ([]models.Card, error) {
	file, err := os.OpenFile(rep.filePath, os.O_RDONLY, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			return make([]models.Card, 0), nil
		}
		return nil, fmt.Errorf("ошибка декодирования JSON: %w", err)
	}
	defer file.Close()

	var cards []models.Card
	err = json.NewDecoder(file).Decode(&cards)
	if err != nil {
		if err == io.EOF {
			return make([]models.Card, 0), nil
		}
		return nil, fmt.Errorf("ошибка декодирования JSON: %w", err)
	}

	return cards, nil
}

func (rep *InJsonRepository) GetWithWrongs() ([]models.Card, error) {
	cards, err := rep.GetAll()
	if err != nil {
		return nil, err
	}

	var filteredCards []models.Card
	for _, card := range cards {
		if card.WrongCount > 0 {
			filteredCards = append(filteredCards, card)
		}
	}

	return filteredCards, nil
}

func (rep *InJsonRepository) SaveAll(cards []models.Card) error {
	file, err := os.Create(rep.filePath)
	if err != nil {
		os.Exit(1)
		fmt.Println("Не удается открыть файл")
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(&cards)
	if err != nil {
		panic(err)
	}

	return nil
}
