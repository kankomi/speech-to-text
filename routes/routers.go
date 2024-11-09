package routes

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"stt/ai"
	"stt/views/components"
	"stt/views/pages/index"

	"github.com/go-chi/chi/v5"
)

func Setup(r *chi.Mux) {

	r.Handle("/*", http.StripPrefix("/public/", http.FileServerFS(os.DirFS("public"))))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		index.Index().Render(r.Context(), w)
	})

	r.Post("/api/upload", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(50 << 20) // max 50 MB
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			println("cannot get file from form")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		f, err := os.CreateTemp(".", strings.TrimSuffix(handler.Filename, filepath.Ext(handler.Filename)))
		if err != nil {
			println("cannot open temp file")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		b, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = f.Write(b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fname := strings.ReplaceAll(fmt.Sprintf("%s%s", f.Name(), filepath.Ext(handler.Filename)), " ", "_")
		err = os.Rename(f.Name(), fname)
		if err != nil {
			println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		println("temp file is", f.Name())
		println("filename is", fname)

		defer func() {
			os.Remove(fname)
		}()

		text, err := ai.SpeechToText(fname)
		if err != nil {
			println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		components.ResultText(*text).Render(r.Context(), w)
	})
}
