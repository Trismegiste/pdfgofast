package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func ExamplePrintToPDF(ctx context.Context, url string, filename string) {
	var buf []byte
	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			buf, _, err = page.PrintToPDF().
				WithDisplayHeaderFooter(false).
				// https://pkg.go.dev/github.com/chromedp/cdproto@v0.0.0-20240709201219-e202069cc16b/page#PrintToPDFParams.WithPrintBackground
				WithPrintBackground(true).
				WithMarginBottom(0).
				WithMarginTop(0).
				WithMarginLeft(0).
				WithMarginRight(0).
				WithPreferCSSPageSize(true).
				Do(ctx)
			return err
		}),
	); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(filename, buf, 0o644); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// https://pkg.go.dev/github.com/chromedp/chromedp#example-NewContext-ReuseBrowser
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	for k := 0; k < 5; k++ {
		start := time.Now()
		ExamplePrintToPDF(ctx, `http://soshomeassist.fr/`, fmt.Sprintf("page%d.pdf", k))
		t := time.Now()
		fmt.Printf("Milliseconds elapsed: %d\n", t.Sub(start).Milliseconds())
	}
}
