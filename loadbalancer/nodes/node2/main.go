package main

func main() {

	server := NewAPIServer(":3002")
	server.Run()
}
