mkdir -p executables/linux
mkdir -p executables/window

# Linux
go build -o ./executables/linux/notificator ./cmd/notificator/.


# Window
GOOS=windows GOARCH=386 go build -o ./executables/window/notificator.exe ./cmd/notificator/.