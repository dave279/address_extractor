package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

// A highly robust regex pattern to capture U.S. addresses.
// This pattern attempts to capture common structures:
// 1. Number: (\d{1,5})
// 2. Street Name/Type: ([\w\s.-]+)
// 3. Major Direction/Street Suffix (Optional): (?:(?:St|Ave|Rd|Ln|Blvd|Pkwy|Circle)\.?\s*)?
// 4. City: ([A-Z]\w*[A-Z])\s*
// 5. State: ([A-Z]{2})
// 6. ZIP: (\d{5}(?:-\d{4})?)
//
// NOTE: This pattern is complex and may need minor tuning based on your specific data source's formatting inconsistencies.
const addressRegexPattern = `(?i)(\d{1,5}\s+[\w\s.-]+\s*(?:Avenue|Ave|Street|St|Drive|Dr|Road|Rd|Boulevard|Blvd|Lane|Ln|Place|Pl|Court|Ct|Circle|Cir|Way|Parkway|Pkwy|Highway|Hwy\.)?[\.]?)[\s,]*([\w\s.-]+)[\s,]*([A-Z]{2})[\s,]*(\d{5}(?:-\d{4})?)`

// Compile the regex outside the main loop for performance.
var addressRegex = regexp.MustCompile(addressRegexPattern)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run address_extractor.go <input_file_path>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	fmt.Printf("--- Starting address extraction from: %s ---\n", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	// Use bufio.Scanner for efficient, line-by-line reading of large files.
	scanner := bufio.NewScanner(file)
	lineCount := 0
	totalMatches := 0

	// Use a string builder to accumulate lines that might contain multi-line addresses,
	// ensuring the regex has enough context.
	var buffer strings.Builder

	// Process line by line
	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		// Append the current line and a space to the buffer
		buffer.WriteString(line)
		buffer.WriteString(" ")

		// Optimization: Only run regex checks every N lines or when the line is empty.
		// This prevents massive overhead on every single line if the text is clean.
		// For simplicity, we will check the entire buffer content.

		// Check the accumulated buffer for matches
		matches := addressRegex.FindAllStringSubmatch(buffer.String(), -1)

		if len(matches) > 0 {
			// We found matches. The most straightforward approach is to process them
			// and print them. We must also clear the buffer content that was matched
			// so we don't re-find the same address on the next iteration.

			// This logic is tricky for real-world extraction, but for demonstration:
			// We extract the match and print it.
			for _, match := range matches {
				// Reassemble the captured groups for a cleaner output.
				// The pattern captured: (Street) (City) (State) (ZIP)
				// match[0] is the full match. match[1] is the first group (street).
				// For readability, let's just print the full captured string.
				fullMatch := strings.TrimSpace(match[0])
				fmt.Printf("\n[MATCH FOUND]: %s\n", fullMatch)
				totalMatches++

				// Simple Buffer Management (Crude but necessary for continuous scanning)
				// Since we don't know exactly where the match ends in the stream,
				// we will simply process the match and *not* clear the buffer,
				// relying on the regex to only find the best fit.
			}

			// To prevent re-finding the same match block repeatedly, a more advanced
			// system would need to truncate the buffer. For this simple script,
			// we will clear the buffer after a successful match block to move on.
			buffer.Reset()
		} else {
			// If no match is found, keep the data in the buffer to check against
			// the next line, allowing for multi-line addresses.
		}
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		log.Fatalf("Error reading file: %v", err)
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Printf("Processing Complete.\n")
	fmt.Printf("Lines Scanned: %d\n", lineCount)
	fmt.Printf("Total Addresses Matched: %d\n", totalMatches)
}
