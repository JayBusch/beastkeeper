cd bin
go build ../src/bk
go build ../src/test
cd ..
go test -v ./src/bk
./bin/test
