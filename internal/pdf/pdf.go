package pdf

import (
	"context"
	"fmt"
	"os"

	"github.com/Pekhov14/resume-snap/internal/config"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// Определение структуры для хранения размеров элемента
type ElementDimensions struct {
	Width, Height, X, Y float64
}

func GeneratePDF(cfg config.Config) error {
	// Настройка контекста Chrome
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Переменная для хранения размеров элемента
	var dimensions ElementDimensions
	var pdfData []byte

	// Запуск браузерных команд
	err := chromedp.Run(ctx,
		chromedp.Navigate(cfg.URL),
		chromedp.WaitVisible(cfg.Selector, chromedp.ByQuery),
		chromedp.Evaluate(getElementDimensionsScript(cfg.Selector), &dimensions),
		chromedp.Evaluate(hideElementsAndScrollScript(cfg.Selector), nil),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfData, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(dimensions.Width / 96).
				WithPaperHeight(dimensions.Height / 96).
				WithMarginTop(0).
				WithMarginBottom(0).
				WithMarginLeft(0).
				WithMarginRight(0).
				WithScale(1.0).
				Do(ctx)
			return err
		}),
	)

	if err != nil {
		return err
	}

	return os.WriteFile(cfg.Output, pdfData, 0644)
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
