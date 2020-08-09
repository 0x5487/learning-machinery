package helloworld



// Say "Hello World"
func Say(name string) (string, error) {
	return "Hello " + name + "!", nil
}