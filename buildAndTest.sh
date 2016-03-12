cd bin
go build ../src/bk
go build ../src/test
cd ..
go test -coverprofile=coverage.out -v ./src/bk
go tool cover -html=coverage.out
./bin/test
