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

func TestParseResponse(t *testing.T) {}
