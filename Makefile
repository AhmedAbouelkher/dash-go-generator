build:
	go build -o app .

run: build
	./app

watch:
	clear
	ulimit -n 1000
	reflex -s -r '\.go$$' make run