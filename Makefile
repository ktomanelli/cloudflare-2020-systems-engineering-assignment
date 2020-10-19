all:
	echo "make commands: help,build,run,clean"
help:
	./main --help=""
build:
	go build src/main.go
run:
	./main --url=$(url) --profile=$(profile)
clean:
	rm main