package validators

import (
	"flag"
	"net"
	"regexp"
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