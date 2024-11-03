package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	port         = "53" //port default pentru DNS
	queryTypeA   = 1    //IPv4
	queryClass   = 1    //IN
	headerSize   = 12   //numarul de biti pentru headerul DNS
	responseSize = 512  //dimensiunea default la raspunsul DNS
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Use: go run dns-client.go <domain> [<DNS server>]")
		return
	}
	domain := os.Args[1]
	server := "8.8.8.8" // DNS implicit: Google

	// dacă utilizatorul specifică un server DNS, îl folosim
	if len(os.Args) > 2 {
		server = os.Args[2]
	}

	packet := createDNSQuery(domain)

	response, err := sendDNSQuery(packet, server)
	if err != nil {
		fmt.Println("Error in sending DNS query:", err)
		return
	}

	for _, ip := range parseResponse(response) {
		fmt.Println(ip)
	}
}

/*
*	createDNSQuery: preia un domeniu si returneaza queryul DNS (sub forma de biti) sau o eroare
*		queryul va fi folosit mai departe sa fie trimis la server
 */
func createDNSQuery(domain string) []byte {
	var query bytes.Buffer

	id := uint16(0xABCD)          //id de interogare (putem avea orice aici)
	flags := uint16(0x0100)       //query standard (QR 0, Opcode 0000, AA 0, TC 0, RD 1, RA 0,  Z 000, RCODE 000)
	questions := uint16(1)        //o intrebare (un singur domeniu)
	answerCounts := uint16(0)     //nu ne intereseaza
	authorityCounts := uint16(0)  //nu ne intereseaza
	additionalCounts := uint16(0) //nu ne intereseaza

	//scriem headerul in queryul de DNS
	binary.Write(&query, binary.BigEndian, id)
	binary.Write(&query, binary.BigEndian, flags)
	binary.Write(&query, binary.BigEndian, questions)
	binary.Write(&query, binary.BigEndian, answerCounts)
	binary.Write(&query, binary.BigEndian, authorityCounts)
	binary.Write(&query, binary.BigEndian, additionalCounts)

	//spargem domeniul in parti separate de punct
	for _, part := range bytes.Split([]byte(domain), []byte(".")) {
		//queryul este de forma [<lungime parte><parte>]
		//unde avem mai intai lungimea partilor, si ce se citeste de acolo
		//exemplu: google.com -> 6google3com (google are len = 6, com are len = 3)
		query.WriteByte(byte(len(part)))
		query.Write(part)
	}
	query.WriteByte(0) //termiantorul de domeniu

	// Question - tip și clasă
	binary.Write(&query, binary.BigEndian, uint16(queryTypeA))
	binary.Write(&query, binary.BigEndian, uint16(queryClass))
	return query.Bytes()
}

/*
*	sendDNSQuery: trimite queryul DNS la server
 */
func sendDNSQuery(packet []byte, server string) ([]byte, error) {
	conn, err := net.Dial("udp", net.JoinHostPort(server, port))
	if err != nil {
		return nil, err
	}
	defer conn.Close() //inchidem conexiunea cu defer

	conn.SetDeadline(time.Now().Add(5 * time.Second))

	// trimite cererea DNS
	_, err = conn.Write(packet)
	if err != nil {
		return nil, err
	}

	response := make([]byte, responseSize) // Dimensiunea standard a unui pachet DNS
	n, err := conn.Read(response)
	if err != nil {
		return nil, err
	}

	return response[:n], nil
}

/*
*	parseResponse: parseaza raspunsul de la server din biti in IP
 */
func parseResponse(response []byte) []string {
	offset := 12 //sarim peste header

	//sarim peste question section
	for response[offset] != 0 {
		offset++
	}

	offset += 5 //sarim si pentru terminator, tip si clasa

	var ips []string

	answerCount := int(binary.BigEndian.Uint16(response[6:8])) // numarul de raspunsuri
	for i := 0; i < answerCount; i++ {
		offset += 10 // sari peste nume, tip, clasă, TTL
		dataLen := binary.BigEndian.Uint16(response[offset : offset+2])
		offset += 2
		ip := net.IP(response[offset : offset+int(dataLen)])
		ips = append(ips, ip.String())
		offset += int(dataLen)
	}
	return ips
}
