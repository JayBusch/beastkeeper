cd bin
go build ../src/bk
go build ../src/test
cd ..
go test -coverprofile=coverage.out -v ./src/bk 
go test -coverprofile=states_coverage.out -v ./src/bk/states
go tool cover -html=coverage.out
go tool cover -html=states_coverage.out
./bin/test
