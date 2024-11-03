package main

const (
	port       = "53" //port default pentru DNS
	queryTypeA = 1    //IPv4
	queryClass = 1    //IN
	headerSize = 12   //numarul de biti pentru headerul DNS
)

func main() {

}

/*
*	createDNSQuery: preia un domeniu si returneaza queryul DNS (sub forma de biti) sau o eroare
*		queryul va fi folosit mai departe sa fie trimis la server
 */
func createDNSQuery(domain string) []byte {
	query := make([]byte, headerSize)
	return query
}

func sendDNSQuery() {}

func parseResponse() {}
