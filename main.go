package main

func main() {
	a := App{}
	// USERNAME e PASSWORD do banco de dados
	a.Initialize("root", "root", "rest_api_example")

	a.Run(":8080")
}
