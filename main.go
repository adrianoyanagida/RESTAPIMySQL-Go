package main

// main : Ponto de entrada para a aplicação
func main() {
	a := App{}
	// USERNAME e PASSWORD do banco de dados
	a.Initialize("root", "root", "rest_api_example")

	a.Run(":8080")
}
