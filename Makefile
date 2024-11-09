dev:
	@go run cmd/main.go

templ:
	@templ generate -watch

tailwind:
	@npx tailwindcss -i input.css -o public/styles.css -w
