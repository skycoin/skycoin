Run cli interface:

Run server:
	go run rpc_run.go
It will run the rpc server to accept messages on localhost on port which environment variable MESH_RPC_PORT is assigned to. If no such variable, it will work on port 1234.

Run client:
	cd cli
	go run cli.go
It will run rpc client which will send message to port 1234. If you want another port to send messages, point it as an argument like this:
	go run cli.go 2222 // will send requests to port 2222

To run client in a browser interface run cli/cli.sh which will open web interface on port 9999, so you can use it in your browser like http://the-url-which-the-client-is-situated-at:9999


