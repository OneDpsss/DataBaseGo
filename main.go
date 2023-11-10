package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Studnet struct {
	Name       string
	Surname    string
	Patronymic string
}

var studnets = make(map[int]Studnet)

var Data string
var DataNames string
var BackupName string
var variants = make(map[int]string)
var marks = make(map[int]int, 0)

var myWindow = app.New().NewWindow("GUI")

func main() {
	myWindow.Resize(fyne.NewSize(1500, 1500))
	createDBButton := widget.NewButton("Создать БД", func() {
		Create()
	})
	editDBButton := widget.NewButton("Изменить БД", func() {
		Update()
	})

	createBackupButton := widget.NewButton("Создать бекап", func() {
		CreateBackup()
	})

	loadBackupButton := widget.NewButton("Загрузить бекап", func() {
		UploadBackup()
		Rewrite()
	})

	showStudentButton := widget.NewButton("Вывести отдельного ученика", func() {
		ShowStudent()
	})

	loadDBButton := widget.NewButton("Загрузить БД", func() {
		UploadData()
	})
	showButton := widget.NewButton("Показать таблицу", func() {
		ShowTable()
	})

	exitButton := widget.NewButton("Завершить программу", func() {
		myWindow.Close()
	})

	buttons := container.NewVBox(
		createDBButton,
		editDBButton,
		createBackupButton,
		loadBackupButton,
		showStudentButton,
		loadDBButton,
		showButton,
		exitButton,
	)
	myWindow.SetContent(buttons)
	myWindow.ShowAndRun()
}
