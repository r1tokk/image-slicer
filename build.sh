go build -o ./build/image-slicer main.go
printf "Compiled successfully!\n"

sudo cp ./build/image-slicer /usr/local/bin/image-slicer
sudo chmod +x /usr/local/bin/image-slicer

printf "Copied successfully into /usr/local/bin/image-slicer\n"
