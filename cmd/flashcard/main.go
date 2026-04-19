package main

import (
	"flashcard/internal/repository"
	"flashcard/internal/ui"
)

const storageFilePath = "data/cards.json"

func main() {
	var rep = repository.NewInJsonRepository(storageFilePath)
	var app = ui.NewApp(rep)

	//Запуск приложения
	app.Run()
}
