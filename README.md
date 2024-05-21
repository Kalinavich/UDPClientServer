Client server UDP application

Build:
cd ./client
go build

cd ./server
go build

Running:
./server PORT

./client HOSTNAME:PORT REQUESTED_RESOURCE
