package main

import (
	"fmt"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port  string `yaml:"port"`
	Pages []Page `yaml:"pages"`
}

type Page struct {
	Page string `yaml:"page"`
	Path string `yaml:"path"`
}

var config Config

var pagesMapping map[string]string = make(map[string]string)

func createDefaultConfig() error {

	defaultConfig := Config{
		Port:  "8085",
		Pages: []Page{{Page: "/example", Path: "files/example.txt"}},
	}

	data, err := yaml.Marshal(&defaultConfig)

	if err != nil {
		return fmt.Errorf("Ошибка при создании файла конфигурации: %v", err)
	}

	err = os.WriteFile("config.yaml", data, 0644)
	if err != nil {
		return fmt.Errorf("Ошибка записи в файл конфигурации: %v", err)
	}

	config = defaultConfig

	fmt.Println("Создан новый файл конфигурации config.yaml с настройками по умолчанию.")
	return nil
}

func loadConfig() error {

	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		return createDefaultConfig()
	}

	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, &config)
}

func main() {

	err := loadConfig()

	if err != nil {
		fmt.Println("Не найден конфиг config.yaml в корневой директории.")
		return
	}

	fmt.Println("Сервис запущен на порту " + config.Port)

	for _, page := range config.Pages {
		http.HandleFunc(page.Page, downloadHandler)
		pagesMapping[page.Page] = page.Path
		fmt.Println("Найден файл " + page.Path + " по URL: localhost:" + config.Port + page.Page)
	}

	http.ListenAndServe(":"+config.Port, nil)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {

	filePath, exists := pagesMapping[r.URL.Path]

	if !exists {
		http.Error(w, "Страница не найдена", http.StatusNotFound)
		return
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Файл должен быть на сервере, но он не был найден. Обратись к администратору.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+filePath)
	w.Header().Set("Content-Type", "application/octet-stream")

	http.ServeFile(w, r, filePath)
}
