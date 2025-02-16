build_run:
	go build -o out && ./out

goose_up:
	cd sql/schema; \
	goose postgres postgres://postgres:postgres@localhost:5432/chirpy up

goose_down:
	cd sql/schema; \
	goose postgres postgres://postgres:postgres@localhost:5432/chirpy down