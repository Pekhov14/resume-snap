package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {
	// Параметры командной строки
	url := flag.String("url", "", "URL страницы для обработки")
	selector := flag.String("selector", "", "CSS селектор элемента для конвертации в PDF")
	output := flag.String("output", "output.pdf", "Имя выходного PDF файла")
	flag.Parse()

	if *url == "" || *selector == "" {
		log.Fatal("Необходимо указать URL и селектор")
	}

	// Настройка контекста Chrome
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Переменные для хранения размеров элемента и PDF данных
	var elementWidth, elementHeight, elementX, elementY float64
	var pdfData []byte

	// Запуск браузерных команд
	err := chromedp.Run(ctx,
		// Переход на страницу
		chromedp.Navigate(*url),

		// Ожидание видимости нужного элемента
		chromedp.WaitVisible(*selector, chromedp.ByQuery),

		// Получение размеров элемента (ширина, высота, и координаты)
		chromedp.Evaluate(fmt.Sprintf(`
			(function() {
				const el = document.querySelector('%s');
				const rect = el.getBoundingClientRect();
				return { width: rect.width, height: rect.height, x: rect.left, y: rect.top };
			})()
		`, *selector), &struct{ Width, Height, X, Y float64 }{elementWidth, elementHeight, elementX, elementY}),

		// Убираем тень и скрываем все элементы, кроме выбранного
		chromedp.Evaluate(fmt.Sprintf(`
			(function() {
				// Убираем тень (box-shadow) у выбранного элемента
				const el = document.querySelector('%s');
				el.style.boxShadow = 'none';
				// Скрытие всех элементов на странице
				document.body.style.visibility = 'hidden';
				// Отображение только нужного элемента
				el.style.visibility = 'visible';
				// Прокрутка страницы к элементу
				window.scrollTo(0, %f);
				el.style.marginTop = '-100px';
			})()
		`, *selector, elementY), nil),

		// Генерация PDF только для видимого контента
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			// Сохраняем PDF для выбранного элемента с точными размерами
			pdfData, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(elementWidth / 96).   // Ширина элемента в дюймах
				WithPaperHeight(elementHeight / 96). // Высота элемента в дюймах
				WithMarginTop(0).
				WithMarginBottom(0).
				WithMarginLeft(0).
				WithMarginRight(0).
				Do(ctx)
			return err
		}),
	)

	// Обработка ошибок
	if err != nil {
		log.Fatal("Ошибка при создании PDF:", err)
	}

	// Сохранение PDF на диск
	err = os.WriteFile(*output, pdfData, 0644)
	if err != nil {
		log.Fatal("Ошибка при записи PDF на диск:", err)
	}

	fmt.Println("PDF успешно сохранен как", *output)
}
