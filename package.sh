env GOOS=darwin
env GOARCH=arm64
env CGO_ENABLED=1
fyne package -icon Icon.png -name build/capyDrop-arm64
env GOARCH=amd64
fyne package -icon Icon.png -name build/capyDrop-amd64
