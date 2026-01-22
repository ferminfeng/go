package main

func main() {
	event, err := InitializeEvent("hello world!", 10000)
	if err != nil {
		panic(err)
	}
	event.Start()
}
