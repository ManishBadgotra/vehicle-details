build:
	GOOS=windows GOARCH=amd64 go build -ldflags="-H=windowsgui" -o ../vehicleGetAPI_service/vehicles.exe .

run: build:
	../vehicleGetAPI_service/vehicles.exe