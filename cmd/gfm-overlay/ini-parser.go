package main

import (
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/SirGFM/gfm-speedrun-overlay/logger"
)

type segment map[string]string
type config map[string]segment

func parseINI(filename string) config {
	if filename == "" {
		return make(config)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		logger.Fatalf("ini: failed to read the config file: %+v", err)
	}

	segments := make(config)
	curSegment := ""

	buf := bytes.NewBuffer(data)
	reading := true
	for lineNumber := 1; reading; lineNumber++ {
		line, err := buf.ReadString('\n')
		if err != nil && err != io.EOF {
			logger.Fatalf("ini: failed to read a config line on line %d: %+v", lineNumber, err)
		}

		reading = (err == nil)

		// Remove comments.
		line, _, _ = strings.Cut(line, ";")
		line, _, _ = strings.Cut(line, "#")
		// Remove whitespaces.
		line = strings.Trim(line, "\t \r\n")
		// Skip empty lines.
		if line == "" {
			continue
		}

		// Check if it's a segment.
		if line[0] == '[' && line[len(line)-1] == ']' {
			curSegment = line[1 : len(line)-1]
			segments[curSegment] = make(segment)
			continue
		}

		// Ensure that a segment is being parsed.
		if curSegment == "" {
			logger.Fatalf("ini: got an attribute outside any segments on line %d", lineNumber)
		}

		// Parse attributes.
		key, value, valid := strings.Cut(line, "=")
		if !valid {
			logger.Fatalf("ini: expected a '<key> = <value>' line for the attribute on line %d", lineNumber)
		}

		// Remove whitespaces.
		key = strings.Trim(key, "\t ")
		value = strings.Trim(value, "\t ")

		// Save the attribute.
		segments[curSegment][key] = value
	}

	return segments
}
