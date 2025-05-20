package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type Record struct {
	BOID    string `json:"boid"`
	BOIDNSE string `json:"boid_nse"`
	BOIDBSE string `json:"boid_bse"`
	NSE     string `json:"nse"`
	BSE     string `json:"bse"`
}

const permittedToTrade = "Permitted to Trade"

const (
	clientCodeColumn     = "CLIENT_CODE"
	exchangeStatusColumn = "EXCHANGE_STATUS"

	cashSegment = "CM"

	NSEMemberID = "90375"
	BSEMemberID = "6867"
)

func loadBSEUCCs(path string) (map[string]bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open BSE file: %w", err)
	}
	defer file.Close()

	uccs := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ucc := strings.TrimSpace(scanner.Text())
		if ucc != "" {
			uccs[ucc] = true
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading BSE file: %w", err)
	}
	log.Printf("Loaded %d BSE UCCs\n", len(uccs))
	return uccs, nil
}

func loadNSEUCCs(path string) (map[string]bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open NSE file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '|'

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read NSE header: %w", err)
	}

	cols := make(map[string]int)
	for i, h := range header {
		cols[h] = i
	}

	codeIdx, ok1 := cols[clientCodeColumn]
	statusIdx, ok2 := cols[exchangeStatusColumn]
	if !ok1 || !ok2 {
		return nil, errors.New("missing expected headers in NSE file")
	}

	uccs := make(map[string]bool)
	for {
		row, err := reader.Read()
		if err != nil {
			break
		}
		if strings.TrimSpace(row[statusIdx]) == permittedToTrade {
			uccs[row[codeIdx]] = true
		}
	}
	log.Printf("Loaded %d NSE UCCs\n", len(uccs))
	return uccs, nil
}

func processCDSL(
	path string,
	nseUCCs, bseUCCs map[string]bool,
) ([]Record, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open CDSL file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	header, err := reader.Read() // skip header
	if err != nil {
		return nil, fmt.Errorf("failed to read CDSL header: %w", err)
	}

	log.Printf("CDSL header has %d columns\n", len(header))

	boidMap := make(map[string]*Record)

	for {
		row, err := reader.Read()
		if err != nil {
			break
		}
		if len(row) < 112 {
			continue
		}
		boid := row[6]
		ucc := row[109]
		code := row[110]
		segment := row[111]

		if segment != cashSegment {
			continue
		}

		rec, ok := boidMap[boid]
		if !ok {
			rec = &Record{BOID: boid}
			boidMap[boid] = rec
		}

		switch code {
		case NSEMemberID:
			if ucc != "" {
				rec.BOIDNSE = ucc
				if nseUCCs[ucc] {
					rec.NSE = ucc
				}
			}
		case BSEMemberID:
			if ucc != "" {
				rec.BOIDBSE = ucc
				if bseUCCs[ucc] {
					rec.BSE = ucc
				}
			}
		}
	}
	log.Printf("Total BOIDs found in CDSL: %d\n", len(boidMap))

	var results []Record
	for _, rec := range boidMap {
		if rec.BOIDNSE == "" || rec.BOIDBSE == "" || rec.NSE == "" || rec.BSE == "" {
			results = append(results, *rec)
		}
	}
	log.Printf("Filtered %d incomplete records\n", len(results))
	return results, nil
}

func writeCSV(records []Record, output string) error {
	file, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"BOID", "BOID-NSE", "BOID-BSE", "NSE", "BSE"})

	for _, r := range records {
		writer.Write([]string{r.BOID, r.BOIDNSE, r.BOIDBSE, r.NSE, r.BSE})
	}
	log.Printf("Output written to %s\n", output)
	return nil
}
