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

func TestSendDNSQuery(t *testing.T) {}

func TestParseResponse(t *testing.T) {}
