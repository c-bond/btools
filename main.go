package main

import (
	grpc "btools/grpc"
	mgr "btools/manager"
)

func main() {
	mgr.StartManager()
	grpc.StartServer()
}
