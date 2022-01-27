for file in ./services/*; do
    echo "Building $file"
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $file/bin/main ./$file
done