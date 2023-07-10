all: go defaults util local

local:
	go build -o nrg main.go
	@chmod +x nrg
	GOOS=linux GOARCH=amd64 go build -o bin/linux/nrg main.go
	@chmod +x bin/linux/nrg
	GOOS=linux GOARCH=arm64 go build -o bin/linux.arm/nrg main.go
	@chmod +x bin/linux.arm/nrg
	GOOS=windows GOARCH=amd64 go build -o bin/windows/nrg.exe main.go
	@chmod +x bin/windows/nrg.exe
	GOOS=windows GOARCH=arm64 go build -o bin/windows.arm/nrg.exe main.go
	@chmod +x bin/windows.arm/nrg.exe
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin/nrg main.go
	@chmod +x bin/darwin/nrg
	GOOS=darwin GOARCH=arm64 go build -o bin/darwin.arm/nrg main.go
	@chmod +x bin/darwin.arm/nrg

util:
	go build -o ~/Utils/nrg main.go
	@chmod +x ~/Utils/nrg

go:
	go build -o ~/go/bin/nrg main.go
	@chmod +x ~/go/bin/nrg

defaults:
	cp -f -r .nrg ~/