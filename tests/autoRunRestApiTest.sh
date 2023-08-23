cd ..
go build -o myprogram main/main.go
./myprogram &
daemon_pid=$!

echo $daemon_pid

sleep 2
cd tests
echo "#####################################Test API FUNCTION########################################"
go test -run TestRegisterFunctionChain
echo "##########################################END#################################################"
echo "####################################Test API PERFORMANCE######################################"
go test -run TestPerformance
echo "##########################################END#################################################"

sleep 5

kill $daemon_pid
rm ../myprogram
