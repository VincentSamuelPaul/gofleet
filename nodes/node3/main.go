package main

func main() {

	server := NewAPIServer(":3003")
	server.Run()
}
