package main

import (
	"fmt"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port     string `yaml:"port"`
	Page     string `yaml:"page"`
	FilePath string `yaml:"file_path"`
}

var config Config

func createDefaultConfig() error {
	defaultConfig := Config{
		Port:     ":8085",
		Page:     "/download-resourcepack",
		FilePath: "resourcepacks/resourcepack.zip",
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

	fmt.Println("Сервис запущен. Ссылка на скачивание файла: localhost" + config.Port + config.Page)

	http.HandleFunc(config.Page, downloadHandler)
	http.ListenAndServe(config.Port, nil)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Disposition", "attachment; filename="+config.FilePath)
	w.Header().Set("Content-Type", "application/octet-stream")

	http.ServeFile(w, r, config.FilePath)
}
