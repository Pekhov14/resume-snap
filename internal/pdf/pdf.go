package pdf

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Pekhov14/resume-snap/internal/config"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type ElementDimensions struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

func GeneratePDF(cfg config.Config) error {
	// Настройка контекста Chrome
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var dimensions ElementDimensions
	var pdfData []byte

	log.Printf("Start generating PDF for page: %s", cfg.URL)

	err := chromedp.Run(ctx,
		chromedp.Navigate(cfg.URL),
		chromedp.WaitVisible(cfg.Selector, chromedp.ByQuery),
		chromedp.Evaluate(preparePageForPDF(cfg.Selector), nil),
		chromedp.Evaluate(getElementDimensionsScript(cfg.Selector), &dimensions),

		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Printf("Element dimensions: %.2fpx x %.2fpx", dimensions.Width, dimensions.Height)
			return nil
		}),

		chromedp.Sleep(1 * time.Second),

		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			paperWidth  := dimensions.Width/96 + 0.5
			paperHeight := 8.69 //11.69 is A4

			log.Printf("PDF paper size: %.2f x %.2f inches", paperWidth, paperHeight)

			pdfData, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(paperWidth).
				WithPaperHeight(paperHeight).
				WithMarginTop(0).
				WithMarginBottom(0).
				WithMarginLeft(0).
				WithMarginRight(0).
				WithScale(.8).
				Do(ctx)
			return err
		}),
	)

	if err != nil {
		return fmt.Errorf("error creating PDF: %w", err)
	}

	if err := os.WriteFile(cfg.Output, pdfData, 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	log.Printf("PDF successfully saved: %s", cfg.Output)
	return nil
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

func preparePageForPDF(selector string) string {
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