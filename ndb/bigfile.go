//
// USDA-NDB Viewer
// Copyright 2013 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package ndb

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// The number of bytes of a file to read, aligned to the nearest newline.
const kChunkSize = 524288

type LineProcessor func(line string) error

type bigFile struct {
	processor LineProcessor // The function that accepts lines.
	started   chan bool     // Used to send messages from processChunk whenever a worker starts.
	done      chan error    // Used to send messages from processChunk or readFile when an error occurs, or nil at EOF.
}

// ReadFile reads the file at path |file| and processes lines using the |processor|
// function. The processor can execute concurrently. Returns the first error that
// occurs. The processor should communicate over a chanel to its own storage facility.
func ReadFile(file string, processor LineProcessor) error {
	bf := &bigFile{
		processor: processor,
		started:   make(chan bool),
		done:      make(chan error),
	}

	go bf.readFile(file)

	// Wait for the workers to complete.
	// Count readFile as the first worker.
	var errs []string
	for i := 1; i > 0; {
		select {
		case <-bf.started:
			i++
		case err := <-bf.done:
			i--
			if err != nil {
				errs = append(errs, err.Error())
			}
		}
	}

	if len(errs) > 0 {
		errStr := strings.Join(errs, "\n\t")
		return fmt.Errorf("ReadFile(%s) encountered the following errors:\n\t%s", file, errStr)
	}

	return nil
}

// readFile synchronously opens and reads the file into chunks. Each chunk is then passed
// to a goroutine, which processes each line of the chunk individually.
func (bf *bigFile) readFile(file string) {
	f, err := os.Open(file)
	if err != nil {
		bf.done <- err
		return
	}
	defer f.Close()

	buf := make([]byte, kChunkSize)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			bf.done <- err
			return
		}

		// Walk the buffer back to the previous newline.
		nOrig := n
		for ; buf[n-1] != '\n'; n-- {
		}

		// If we moved at all, seek the file backwards.
		if delta := n - nOrig; delta < 0 {
			_, err := f.Seek(int64(delta), 1)
			if err != nil {
				bf.done <- err
				return
			}
		}

		// Proceess the chunk.
		go bf.processChunk(string(buf[:n]))
	}

	bf.done <- nil

	return
}

// processChunk splits a chunk on \r\n and sends each whitespace-trimmed line to
// the processor. Any errors from the processor are sent over the error channel,
// or nil is sent when processing is done.
func (bf *bigFile) processChunk(chunk string) {
	bf.started <- true

	var start int
	for i := 0; i < len(chunk); i++ {
		if chunk[i] == '\r' {
			if err := bf.processor(chunk[start:i]); err != nil {
				bf.done <- err
				return
			}
			start = i + 1
		} else if chunk[i] == '\n' {
			start = i + 1
		}
	}

	bf.done <- nil
}
