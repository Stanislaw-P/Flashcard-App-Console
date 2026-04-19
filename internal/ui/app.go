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
	app.help()

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
	case "train", "full training":
		app.startFullTraining()
	case "train wr", "train wrong cards":
		app.startWrongCardTraining()
	case "help":
		app.help()
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
	app.trainCards(app.cards, "Тренировка всех существующих слов")
}

func (app *App) startWrongCardTraining() {
	cardsWithWrongs, err := app.repository.GetWithWrongs()

	if err != nil {
		fmt.Printf("Ошибка получения карточек: %v\n", err)
		return
	}

	if len(cardsWithWrongs) == 0 {
		fmt.Println("⚠️ Нет карточек, в которых вы ошибались! Сначала потренируйтесь на всех карточках.")
		return
	}

	app.trainCards(cardsWithWrongs, "Тренировка ошибок")
}

func (app *App) trainCards(cards []models.Card, trainingName string) {
	fmt.Printf("--- %s ---\n", trainingName)
	fmt.Println("Введите 'back' для выхода в главное меню")
	fmt.Println()

	cardsCount := len(cards)
	if cardsCount == 0 {
		fmt.Printf("⚠️ Нет карточек для тренировки!\n")
		return
	}

	// Создаем копию и перемешиваем
	trainingCards := make([]models.Card, cardsCount)
	copy(trainingCards, cards)

	rand.Shuffle(cardsCount, func(i, j int) {
		trainingCards[i], trainingCards[j] = trainingCards[j], trainingCards[i]
	})

	score := 0

	for i, card := range trainingCards {
		fmt.Printf("Слово %d/%d: %s\n", i+1, len(trainingCards), card.Word)
		fmt.Print("Введите перевод (или 'back' для выхода) -> ")

		answer, _ := app.reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))

		if answer == "back" {
			fmt.Println("Тренировка прервана.")
			app.repository.SaveAll(app.cards)
			return
		}

		// Находим оригинальную карточку
		originalIndex := -1
		for idx, origCard := range app.cards {
			if origCard.ID == card.ID {
				originalIndex = idx
				break
			}
		}

		if originalIndex == -1 {
			fmt.Println("Ошибка: карточка не найдена")
			continue
		}

		if answer == card.Translation {
			fmt.Println("✅ Правильно!")
			score++
			app.cards[originalIndex].CorrectCount++

			if app.cards[originalIndex].WrongCount > 0 {
				app.cards[originalIndex].WrongCount--
			}
		} else {
			fmt.Printf("❌ Неправильно! Правильный ответ: %s\n", card.Translation)
			app.cards[originalIndex].WrongCount++
		}
		fmt.Println()
	}

	// Сохраняем обновленную статистику
	err := app.repository.SaveAll(app.cards)
	if err != nil {
		fmt.Printf("Ошибка сохранения: %v\n", err)
	}

	fmt.Printf("\n📊 Результат: %d/%d правильных ответов (%.1f%%)\n",
		score, len(trainingCards), float64(score)/float64(len(trainingCards))*100)
}

func (app *App) help() {
	fmt.Println("=== Flashcard App ===")
	fmt.Println("Commands:")
	fmt.Println("help")
	fmt.Println("add - add a new card")
	fmt.Println("list - get a list of all cards")
	fmt.Println("train - start training all cards")
	fmt.Println("train wr- start training the maps where I made mistakes")
	fmt.Println("stats - get stats")
	fmt.Println("exit - exit from app")
	fmt.Println()

}
