package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

func old_main() {
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	http.ServeFile(w, r, "index.html")
}

const MAX_UPLOAD_SIZE = 5 * 1024 * 1024 // 5MB

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// mostly inspired by https://freshman.tech/file-upload-golang/
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	pathname := fmt.Sprintf("/tmp/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
	dst, err := os.Create(pathname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	//fmt.Fprintf(w, "Upload successful")

	var buf []byte
	if err := chromedp.Run(ctx,
		chromedp.Navigate("file://"+pathname),
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

	w.Write(buf)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/upload", uploadHandler)

	if err := http.ListenAndServe(":4500", mux); err != nil {
		log.Fatal(err)
	}
}
