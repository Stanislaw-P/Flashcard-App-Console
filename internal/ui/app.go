package ui

import (
	"bufio"
	"flashcard/internal/models"
	"flashcard/internal/repository"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

type App struct {
	running    bool
	reader     *bufio.Reader
	repository repository.CardRepository
	cards      []models.Card
}

func NewApp(rep repository.CardRepository) *App {
	cards, _ := rep.GetAll()
	return &App{
		running:    true,
		reader:     bufio.NewReader(os.Stdin),
		repository: rep,
		cards:      cards,
	}
}

func (app *App) Run() {
	fmt.Println("=== Flashcard App ===")
	fmt.Println("Commands: add, list, train, stats, exit")
	fmt.Println()

	for app.running {
		fmt.Print("> ")
		input, _ := app.reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		app.handleCommand(input)
	}
}

func (app *App) handleCommand(cmd string) {
	switch cmd {
	case "exit", "quit", "q":
		app.running = false
		fmt.Println("До свидания!")
	case "list":
		app.listCards()
	case "add":
		app.add()
	case "full", "full training":
		app.startFullTraining()
	default:
		if cmd != "" {
			fmt.Printf("Неизвестная команда: %s\n", cmd)
		}
	}
}

func (app *App) listCards() {
	fmt.Println("--- Список всех карт ---")

	if len(app.cards) == 0 {
		fmt.Println("Сейчас список карт пуст")
		return
	}

	for i, card := range app.cards {
		fmt.Printf("%d. %s -> %s (угадано: %d, ошибок: %d)\n",
			i+1, card.Word, card.Translation, card.CorrectCount, card.WrongCount)
	}
}

func (app *App) add() {
	fmt.Println("--- Добавление новой карты ---")

	fmt.Print("Слово: ")
	word, _ := app.reader.ReadString('\n')
	word = strings.TrimSpace(strings.ToLower(word))

	fmt.Print("Перевод: ")
	trans, _ := app.reader.ReadString('\n')
	trans = strings.TrimSpace(strings.ToLower(trans))

	newCard := models.Card{
		ID:          len(app.cards) + 1,
		Word:        word,
		Translation: trans,
	}

	app.cards = append(app.cards, newCard)

	err := app.repository.SaveAll(app.cards)
	if err != nil {
		panic(err)
	}
	fmt.Println("- Карточка сохранена -")
}

func (app *App) startFullTraining() {
	fmt.Println("--- Тренировка всех существующих слов ---")
	fmt.Println("Введите 'back' для выхода в главное меню")
	fmt.Println()

	currentCount := 0
	cardsCount := len(app.cards)
	score := 0

	for {
		// Проверяем не закончили ли слова
		if currentCount == cardsCount {
			fmt.Println("- ⚠️ Нет карточек для тренировки! Добавьте карточки -")
			break
		}

		rndInd := rand.Intn(cardsCount)
		rndCard := app.cards[rndInd]

		currentCount++ // Номер текущего шага

		fmt.Printf("Слово %d/%d: %s\n", currentCount, cardsCount, rndCard.Word)
		fmt.Print("Введите перевод (или 'back' для выхода в главное меню) -> ")

		answer, _ := app.reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))

		if answer == "back" {
			break
		}

		if answer == rndCard.Translation {
			fmt.Println("- ✅ Правильно! Выбираю следующее слово... -")
			fmt.Println()
			score++
			app.cards[rndInd].CorrectCount++
		} else {
			fmt.Printf("- ❌ Неправильно! Правильный ответ: %s\n-", rndCard.Translation)
			fmt.Println()
			app.cards[rndInd].WrongCount++
		}
	}

	// Сохраняем обновленную статистику
	app.repository.SaveAll(app.cards)

	fmt.Printf("\n📊 Результат: %d/%d правильных ответов (%.1f%%)\n",
		score, cardsCount, float64(score)/float64(cardsCount)*100)
}
