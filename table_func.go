package main

import (
	"bufio"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func Create() {
	entryDBName := widget.NewEntry()
	entryFileName := widget.NewEntry()
	dialog.ShowForm("Введите имя БД и имя файла", "Подтвердить", "Отмена", []*widget.FormItem{
		{Text: "Имя БД", Widget: entryDBName},
		{Text: "Имя файла с именами", Widget: entryFileName},
	}, func(accepted bool) {
		if accepted {
			dbName := entryDBName.Text
			Data = dbName
			fileName := entryFileName.Text
			DataNames = fileName
			if dbName != "" {
				DataFile, _ := os.Create(Data + ".txt")
				Names, _ := os.Open(DataNames + ".txt")
				ScannerFile := bufio.NewScanner(Names)
				defer Names.Close()
				defer DataFile.Close()
				id := 0
				for ScannerFile.Scan() {
					words := strings.Fields(ScannerFile.Text())
					if len(words) == 0 {
						return
					}
					studnets[id] = Studnet{Name: words[1], Surname: words[0], Patronymic: words[2]}
					a := rand.Int() % 50
					variants[id] = "var" + strconv.Itoa(a)
					marks[id] = 0
					DataFile.WriteString(strconv.Itoa(id) + " " + studnets[id].Surname + " " + studnets[id].Name + " " + studnets[id].Patronymic + " " + variants[id] + " " + strconv.Itoa(marks[id]) + "\n")
					id++
				}
				dialog.ShowInformation("Успех", "База данных '"+dbName+"' создана", myWindow)
			} else {
				dialog.ShowInformation("Ошибка", "Пожалуйста, введите имя БД и имя файла", myWindow)
			}
		}
	}, myWindow)

}

func Rewrite() {
	Data, _ := os.OpenFile(Data+".txt", os.O_RDWR|os.O_TRUNC, 0600)
	defer Data.Close()
	for id := 0; id < 1000; id++ {
		_, ok := studnets[id]
		if !ok {
			continue
		}
		Data.WriteString(strconv.Itoa(id) + " " + studnets[id].Surname + " " + studnets[id].Name + " " + studnets[id].Patronymic + " " + variants[id] + " " + strconv.Itoa(marks[id]) + "\n")
	}
}

func check_collusion(surname string, name string, patr string) bool {
	for _, val := range studnets {
		if val.Name == name && val.Surname == surname && val.Patronymic == patr {
			fmt.Println("Студент с таким ФИО уже есть")
			return false
		}
	}
	return true
}
func Remove() {
	entryFullName := widget.NewEntry()
	dialog.ShowForm("Введите ФИО студента", "Подтвердить", "Отмена", []*widget.FormItem{
		{Text: "ФИО студента", Widget: entryFullName},
	}, func(accepted bool) {
		if accepted {
			words := strings.Fields(entryFullName.Text)
			name := words[1]
			surname := words[0]
			patr := words[2]
			for id := 0; id < 1000; id++ {
				if studnets[id].Name == name && studnets[id].Surname == surname && studnets[id].Patronymic == patr {
					delete(marks, id)
					delete(variants, id)
					delete(studnets, id)
					dialog.ShowInformation("Успех", "Студент удален", myWindow)
					Rewrite()
					return
				}
			}
			dialog.ShowInformation("Ошибка", "Студент не найден", myWindow)
		}
	}, myWindow)
}

func UploadData() {
	marks = make(map[int]int, 0)
	variants = make(map[int]string)
	studnets = make(map[int]Studnet)
	entryData := widget.NewEntry()
	dialog.ShowForm("Введите название загружаемой БД", "Подтвердить", "Отмена", []*widget.FormItem{
		{Text: "Название БД", Widget: entryData},
	}, func(accepted bool) {
		Data = entryData.Text
		DataFile, err := os.Open(Data + ".txt")
		if err != nil {
			dialog.ShowInformation("Ошибка", "Файл не найден", myWindow)
		}
		ScannerFile := bufio.NewScanner(DataFile)
		defer DataFile.Close()
		for ScannerFile.Scan() {
			words := strings.Fields(ScannerFile.Text())
			id, _ := strconv.Atoi(words[5])
			studnets[id] = Studnet{Name: words[1], Surname: words[0], Patronymic: words[2]}
			marks[id], _ = strconv.Atoi(words[4])
			variants[id] = words[3]
			fmt.Println(studnets[id].Name, marks[id])
		}
		dialog.ShowInformation("Успех", "База данных загружена", myWindow)
	}, myWindow)
}

func Add() {
	entryFullName := widget.NewEntry()
	dialog.ShowForm("Введите ФИО студента", "Подтвердить", "Отмена", []*widget.FormItem{
		{Text: "ФИО студента", Widget: entryFullName},
	}, func(accepted bool) {
		if accepted {
			words := strings.Fields(entryFullName.Text)
			name := words[1]
			surname := words[0]
			patr := words[2]
			if !check_collusion(surname, name, patr) {
				dialog.ShowInformation("Ошибка", "Студент с таким именем уже есть", myWindow)
			} else {
				for i := 0; i < 1000; i++ {
					_, ok := studnets[i]
					if !ok {
						studnets[i] = Studnet{Name: name, Surname: surname, Patronymic: patr}
						marks[i] = 0
						a := rand.Int() % 50
						variants[i] = "var" + strconv.Itoa(a)
						dialog.ShowInformation("Успех", "Студент добавлен", myWindow)
						Rewrite()
						return
					}
				}
			}
		}
	}, myWindow)
}

func Edit() {
	entryFullName := widget.NewEntry()
	entryMark := widget.NewEntry()
	dialog.ShowForm("Введите ФИО студента", "Подтвердить", "Отмена", []*widget.FormItem{
		{Text: "ФИО студента", Widget: entryFullName},
		{Text: "Новую оценку", Widget: entryMark},
	}, func(accepted bool) {
		if accepted {
			words := strings.Fields(entryFullName.Text)
			name := words[1]
			surname := words[0]
			patr := words[2]
			mark := entryMark.Text
			for id := 0; id < 1000; id++ {
				if studnets[id].Name == name && studnets[id].Surname == surname && studnets[id].Patronymic == patr {
					marks[id], _ = strconv.Atoi(mark)
					dialog.ShowInformation("Успех", "Оценка изменена", myWindow)
					Rewrite()
					return
				}
			}
			dialog.ShowInformation("Ошибка", "Студент не найден", myWindow)
		}
	}, myWindow)
}

func ShowStudent() {
	entryFullName := widget.NewEntry()
	dialog.ShowForm("Введите ФИО студента", "Подтвердить", "Отмена", []*widget.FormItem{
		{Text: "ФИО студента", Widget: entryFullName},
	}, func(accepted bool) {
		if accepted {
			words := strings.Fields(entryFullName.Text)
			name := words[1]
			surname := words[0]
			patr := words[2]
			for id := 0; id < 1000; id++ {
				if studnets[id].Name == name && studnets[id].Surname == surname && studnets[id].Patronymic == patr {
					msg := surname + " " + name + " " + patr + " " + variants[id] + " " + strconv.Itoa(marks[id])
					entry := widget.NewEntry()
					entry.SetText(msg)
					entry.Resize(fyne.NewSize(1000, 1000))
					show := dialog.NewCustom("Данные", "Закрыть", entry, myWindow)
					show.Resize(fyne.NewSize(600, 200))
					show.Show()
					return
				}
			}
			dialog.ShowInformation("Ошибка", "Студент не найден", myWindow)
		}
	}, myWindow)
}

func ShowTable() {
	str := ""
	for id := 0; id <= 1000; id++ {
		_, ok := studnets[id]
		if ok {
			msg := studnets[id].Surname + " " + studnets[id].Name + " " + studnets[id].Patronymic + " " + variants[id] + " " + strconv.Itoa(marks[id])
			str += msg + "\n"
		}
	}
	text := widget.NewMultiLineEntry()
	text.SetText(str)
	container := container.NewScroll(text)
	text.Resize(fyne.NewSize(1000, 1000))
	show := dialog.NewCustom("Данные", "Закрыть", container, myWindow)
	show.Resize(fyne.NewSize(1000, 500))
	show.Show()
}
func Update() {
	content := container.NewVBox(
		widget.NewButton("Удалить", func() {
			Remove()
		}),
		widget.NewButton("Изменить оценку", func() {
			Edit()
		}),
		widget.NewButton("Добавить", func() {
			Add()
		}),
	)
	dialog.ShowCustom("Выберите действие", "Отмена", content, myWindow)
}

func CreateBackup() {
	entryBdName := widget.NewEntry()
	entryBcName := widget.NewEntry()
	dialog.ShowForm("Введите имя БД и имя файла", "Подтвердить", "Отмена", []*widget.FormItem{
		{Text: "Имя БД", Widget: entryBdName},
		{Text: "Имя бекапа бд", Widget: entryBcName},
	}, func(accepted bool) {
		sourceFile := entryBdName.Text + ".txt"
		destinationFile := entryBcName.Text + ".txt"
		source, err := os.Open(sourceFile) //open the source file
		if err != nil {
			panic(err)
		}
		defer source.Close()
		ScannerFile := bufio.NewScanner(source)
		for ScannerFile.Scan() {
			words := strings.Fields(ScannerFile.Text())
			id, _ := strconv.Atoi(words[0])
			studnets[id] = Studnet{Name: words[2], Surname: words[1], Patronymic: words[3]}
			marks[id], _ = strconv.Atoi(words[5])
			variants[id] = words[4]
		}
		destination, err := os.Create(destinationFile) //create the destination file
		if err != nil {
			panic(err)
		}
		defer destination.Close()
		for id := 0; id < 1000; id++ {
			_, ok := studnets[id]
			if !ok {
				continue
			}
			destination.WriteString(studnets[id].Surname + " " + studnets[id].Name + " " + studnets[id].Patronymic + " " + variants[id] + " " + strconv.Itoa(marks[id]) + " " + strconv.Itoa(id) + "\n")
		}
		if err != nil {
			panic(err)
		}
		dialog.ShowInformation("Успех", "Бэкап создан", myWindow)
	}, myWindow)
}

func UploadBackup() {
	marks = make(map[int]int, 0)
	variants = make(map[int]string)
	studnets = make(map[int]Studnet)
	entryBdName := widget.NewEntry()
	entryBcName := widget.NewEntry()
	dialog.ShowForm("Введите имя БД и имя файла", "Подтвердить", "Отмена", []*widget.FormItem{
		{Text: "Имя БД", Widget: entryBdName},
		{Text: "Имя бекапа БД", Widget: entryBcName},
	}, func(accepted bool) {
		sourceFile := entryBcName.Text + ".txt"
		Data = entryBdName.Text
		source, err := os.Open(sourceFile)
		if err != nil {
			dialog.ShowInformation("Ошибка", "Файл не найден", myWindow)
			return
		}
		defer source.Close()
		ScannerFile := bufio.NewScanner(source)
		for ScannerFile.Scan() {
			fmt.Println(ScannerFile.Text())
			words := strings.Fields(ScannerFile.Text())
			id, _ := strconv.Atoi(words[5])
			studnets[id] = Studnet{Name: words[1], Surname: words[0], Patronymic: words[2]}
			marks[id], _ = strconv.Atoi(words[4])
			variants[id] = words[3]
		}
		dialog.ShowInformation("Успех", "Бэкап загружен", myWindow)
		Rewrite()
	}, myWindow)
}
