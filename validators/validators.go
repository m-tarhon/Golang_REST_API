package validators

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

func IsValid(port string) bool {
	pattern := `^([1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$`
	match, _ := regexp.MatchString(pattern, port)
	return match
}

func IsAvailable(port string) bool{
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil{
		return false
	}
	ln.Close()
	return true
}

func IsFlag(name string) bool{
	found := false
	flag.Visit(func( f *flag.Flag){
		if f.Name == name {
			found = true
		}
	})
	return found
}


func SMTH(flag string, port *string){

	fmt.Printf("Hello Server, what %s port would you like to connect to?", flag)

	for {
		reader := bufio.NewReader(os.Stdin)
		port1, err0 := reader.ReadString('\n')
		if err0 != nil {
			fmt.Print("STOPPED.")
			return 
		}

		port1 = strings.TrimSpace(port1)
		if port1 == "" {
			fmt.Println("Port cannot be empty, please try again.")
			continue
		}

		if !IsValid(port1) {
			fmt.Println("Invalid format. Pick a port from 1-65336.")
			continue
		}

		if !IsAvailable(port1) {
			fmt.Println("This port is busy. Try something else.")
			continue
		}
		*port = port1
		break
	}
}
