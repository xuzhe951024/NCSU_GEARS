cd ..
go build -o myprogram main/main.go
./myprogram &
daemon_pid=$!

echo $daemon_pid

sleep 2
cd tests
go test -run TestRegisterFunctionChain

sleep 5

kill $daemon_pid
rm ../myprogram