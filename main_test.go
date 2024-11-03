package main

import (
	"bytes"
	"testing"
)

// test pentru createDNSQuery
func TestCreateDNSQuery(t *testing.T) {

	domain := "google.com"
	query := createDNSQuery(domain)

	if len(query) < headerSize {
		t.Errorf("Length of query is not correct: got %d, expected at least %d", len(query), headerSize)
	}

	expectedBytes := []byte{6, 'g', 'o', 'o', 'g', 'l', 'e', 3, 'c', 'o', 'm'}

	if !bytes.Contains(query, expectedBytes) {
		t.Errorf("Query incorrect: got %v, expected %v", query, expectedBytes)
	}
}

// test pentru sendDNSQuery
func TestSendDNSQuery(t *testing.T) {
	domain := "google.com"
	server := "8.8.8.8"

	query := createDNSQuery(domain)

	response, err := sendDNSQuery(query, server)
	if err != nil {
		t.Errorf("Error in sending DNS query: %v", err)
	}

	if !bytes.Contains(response, query[12:]) {
		t.Errorf("Response does not contain initial query: response %v, query %v", response, query)
	}
}

func TestParseResponse(t *testing.T) {
	// Simulăm un răspuns DNS care conține o singură adresă IP
	response := []byte{
		0x12, 0x34, 0x81, 0x80, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, // Header
		0x06, 'g', 'o', 'o', 'g', 'l', 'e', 0x03, 'c', 'o', 'm', 0x00, // Question
		0x00, 0x01, 0x00, 0x01, // Type A, Class IN
		0xc0, 0x0c, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x3c, 0x00, 0x04, // Answer
		0x85, 0xfb, 0xd0, 0xae, // IP: 142.251.208.174 (exemplu pentru google.com)
	}

	parsedIPs, err := parseResponse(response)
	if err != nil {
		t.Fatalf("Erorr in parsing response: %v", err)
	}

	expectedIP := "93.184.216.34"
	if len(parsedIPs) == 0 || parsedIPs[0] != expectedIP {
		t.Errorf("Error in getting IP address: got %s, expected %v", parsedIPs, expectedIP)
	}
}
