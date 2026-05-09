go build -o ./build/image-slicer main.go
printf "Compiled successfully!\n"

cp ./build/image-slicer /usr/local/bin/image-slicer
chmod +x /usr/local/bin/image-slicer

printf "Copied successfully into /usr/local/bin/image-slicer\n"
