package main

func main() {
	a := App{}
	a.Initialize("root", "", "localhost", "midwestemo")
	a.Run(":8010")
}
