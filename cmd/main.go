package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// структура для данных о планете
type Planet struct {
	Name      string   `json:"name"`
	Gravity   string   `json:"gravity"`
	Residents []string `json:"residents"`
}

// функция для запроса данных о планете
func getPlanet(id int) (*Planet, error) {
	url := fmt.Sprintf("https://swapi.dev/api/planets/%d/", id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var planet Planet
	if err := json.NewDecoder(resp.Body).Decode(&planet); err != nil {
		return nil, err
	}

	if planet.Name == "" {
		return nil, errors.New("planet not found or invalid response")
	}

	return &planet, nil
}

// функция для вывода информации о планете
func printPlanet(p *Planet) {
	fmt.Println("Название:", p.Name)
	fmt.Println("Гравитация:", p.Gravity)
	fmt.Println("Жители:")
	if len(p.Residents) == 0 {
		fmt.Println("  (нет данных)")
		return
	}

	type person struct {
		Name string `json:"name"`
	}

	for _, residentURL := range p.Residents {
		resp, err := http.Get(residentURL)
		if err != nil {
			fmt.Println("  - ошибка загрузки жителя:", err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			fmt.Println("  - ошибка загрузки жителя:", resp.Status)
			resp.Body.Close()
			continue
		}
		var r person
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			fmt.Println("  - ошибка чтения жителя:", err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()
		if strings.TrimSpace(r.Name) == "" {
			fmt.Println("  - (неизвестно)")
		} else {
			fmt.Println("  -", r.Name)
		}
	}
}

// основная функция программы
func main() {
	current := 1 // начинаем с первой планеты

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Загружаем планету", current)

		planet, err := getPlanet(current)
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

		printPlanet(planet)

		fmt.Print("Введите команду (next, back, exit): ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(strings.ToLower(cmd))

		switch cmd {
		case "next":
			current++
		case "back":
			if current > 1 {
				current--
			} else {
				fmt.Println("Нельзя меньше 1")
			}
		case "exit":
			fmt.Println("Выход...")
			return
		default:
			fmt.Println("Неизвестная команда. Доступно: next, back, exit")
		}
	}
}
