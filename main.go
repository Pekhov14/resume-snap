package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type Config struct {
	URL      string `json:"url"`
	Selector string `json:"selector"`
	Output   string `json:"output"`
}

type ElementDimensions struct {
	Width, Height, X, Y float64
}

func main() {
	// Загрузка конфигурации
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	// Параметры командной строки
	url, selector, output := parseFlags(config)

	// Настройка контекста Chrome
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Переменная для хранения размеров элемента
	var dimensions ElementDimensions
	var pdfData []byte

	// Запуск браузерных команд
	err = chromedp.Run(ctx,
		// Переход на страницу
		chromedp.Navigate(url),

		// Ожидание видимости нужного элемента
		chromedp.WaitVisible(selector, chromedp.ByQuery),

		// Получение размеров элемента
		chromedp.Evaluate(getElementDimensionsScript(selector), &dimensions),

		// Убираем тень и скрываем все элементы, кроме выбранного
		chromedp.Evaluate(hideElementsAndScrollScript(selector), nil),

		// Генерация PDF только для видимого контента
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			// Сохраняем PDF для выбранного элемента с точными размерами
			pdfData, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(dimensions.Width / 96).   // Ширина элемента в дюймах
				WithPaperHeight(dimensions.Height / 96). // Высота элемента в дюймах
				WithMarginTop(0).                        // Убираем верхний отступ
				WithMarginBottom(0).                     // Убираем нижний отступ
				WithMarginLeft(0).                       // Убираем левый отступ
				WithMarginRight(0).                      // Убираем правый отступ
				WithScale(0.9).                          // Устанавливаем масштаб
				Do(ctx)
			return err
		}),
	)

	// Обработка ошибок
	if err != nil {
		log.Fatal("Ошибка при создании PDF:", err)
	}

	// Сохранение PDF на диск
	err = os.WriteFile(output, pdfData, 0644)
	if err != nil {
		log.Fatal("Ошибка при записи PDF на диск:", err)
	}

	fmt.Println("PDF успешно сохранен как", output)
}

func parseFlags(config Config) (string, string, string) {
	url := flag.String("url", config.URL, "URL страницы для обработки")
	selector := flag.String("selector", config.Selector, "CSS селектор элемента для конвертации в PDF")
	output := flag.String("output", config.Output, "Имя выходного PDF файла")
	flag.Parse()

	if *url == "" || *selector == "" {
		log.Fatal("Необходимо указать URL и селектор")
	}

	return *url, *selector, *output
}

func loadConfig(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, err
	}
	return config, nil
}

func getElementDimensionsScript(selector string) string {
	return fmt.Sprintf(`
		(function() {
			const el = document.querySelector('%s');
			const rect = el.getBoundingClientRect();
			return { width: rect.width, height: rect.height, x: rect.left, y: rect.top };
		})()
	`, selector)
}

func hideElementsAndScrollScript(selector string) string {
	return fmt.Sprintf(`
		(function() {
			const el = document.querySelector('%s');
			el.style.boxShadow = 'none';
			document.body.style.visibility = 'hidden';
			el.style.visibility = 'visible';
			window.scrollTo(0, el.getBoundingClientRect().top);
			el.style.marginTop = '-100px';
		})()
	`, selector)
}
