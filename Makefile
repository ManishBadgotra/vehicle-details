build:
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o vehicle_details.exe . 

run: build:
	./vehicle_details.exe