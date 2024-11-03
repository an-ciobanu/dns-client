package main

import (
	"bytes"
	"encoding/binary"
)

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
	return query.Bytes()
}

func sendDNSQuery(query []byte, domain string) ([]byte, error) {
	return make([]byte, headerSize), nil
}

func parseResponse() {}
