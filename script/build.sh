dir=$(ls -l ./cmd/ |awk '/^d/ {print $NF}')
for i in $dir
do
    echo "\033[36m build $i start... \033[0m"
    GOOS=linux GOARCH=amd64 go build -o ./build/$i ./cmd/$i/main.go
    echo "\033[32m build $i finish! \033[0m"
done
cp ./config.prod.toml ./build/config.toml
cp -r ./lang ./build/
echo "\033[42;37m build all finished! \033[0m"