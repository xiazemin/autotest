go get -u github.com/mjibson/esc

go generate ./...       

go run gotests/main.go -only GreetVisitors /Users/xiazemin/software/tabel_drving_test_learn/exp4/party.go 


go run gotests/main.go -only GreetVisitors /Users/xiazemin/source/tabel_drving_test_learn/exp4/party.go  > /Users/xiazemin/source/tabel_drving_test_learn/exp4/DATA/min3.go


go build -o gotests_bin gotests/main.go

cp gotests_bin /Users/xiazemin/go/bin/gotests

go generate ./... && go build -o gotests_bin gotests/main.go && cp gotests_bin /Users/xiazemin/go/bin/gotests


% go generate ./... && go run gotests/main.go -w -only GreetVisitors /Users/xiazemin/source/tabel_drving_test_learn/exp4/party.go   


 % go get -u github.com/cweill/gotests/...

 https://tonybai.com/2020/12/10/a-kind-of-thinking-about-how-to-trace-function-call-chain/

 https://www.coder.work/article/194041