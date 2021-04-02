go get -u github.com/mjibson/esc

go generate ./...       

go run gotests/main.go -only GreetVisitors /Users/xiazemin/software/tabel_drving_test_learn/exp4/party.go 


go run gotests/main.go -only GreetVisitors /Users/xiazemin/software/tabel_drving_test_learn/exp4/party.go  > /Users/xiazemin/software/tabel_drving_test_learn/exp4/DATA/min3.go


go build -o gotests_bin gotests/main.go

cp gotests_bin /Users/xiazemin/go/bin/gotests