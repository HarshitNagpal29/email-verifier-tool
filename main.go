package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("domain,hasMX,hasSPF,spfRecord,hasDMARC,dmarcRecord\n")

	for scanner.Scan() {
		checkDomain(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func checkDomain(domain string) {
	var hasMX, hasSPF, hasDMARC bool
	var spfRecords, dmarcRecords string

	// Check MX records
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error looking up MX records for %s: %v", domain, err)
	} else {
		hasMX = len(mxRecords) > 0
	}

	// Check TXT records
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error looking up TXT records for %s: %v", domain, err)
	} else {
		for _, record := range txtRecords {
			if strings.HasPrefix(record, "v=spf1") {
				hasSPF = true
				spfRecords = record
			}
			if strings.HasPrefix(record, "v=DMARC1") {
				hasDMARC = true
				dmarcRecords = record
			}
		}
	}

	fmt.Printf("%s,%t,%t,%s,%t,%s\n", domain, hasMX, hasSPF, spfRecords, hasDMARC, dmarcRecords)
}
