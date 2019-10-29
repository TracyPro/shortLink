package main

func main() {
	a := App{}
	a.Initialize(getEnv())
	a.Run("localhost:8000")
}
