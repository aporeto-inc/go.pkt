/*
 * Network packet analysis framework.
 *
 * Copyright (c) 2014, Alessandro Ghedini
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are
 * met:
 *
 *     * Redistributions of source code must retain the above copyright
 *       notice, this list of conditions and the following disclaimer.
 *
 *     * Redistributions in binary form must reproduce the above copyright
 *       notice, this list of conditions and the following disclaimer in the
 *       documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS
 * IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO,
 * THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR
 * PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR
 * CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL,
 * EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
 * PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
 * PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 * LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
 * NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package file_test

import "log"
import "testing"

import "github.com/ghedo/hype/capture/file"
import "github.com/ghedo/hype/filter"

func TestCapture(t *testing.T) {
	src, err := file.Open("capture_test.pcap")
	if err != nil {
		t.Fatalf("Error opening: %s", err)
	}

	var count uint64
	for {
		raw_pkt, err := src.Capture()
		if err != nil {
			t.Fatalf("Error reading: %s", err)
		}

		if raw_pkt == nil {
			break
		}

		count++
	}

	if count != 16 {
		t.Fatalf("Count mismatch: %d", count)
	}
}

func TestCaptureFilter(t *testing.T) {
	src, err := file.Open("capture_test.pcap")
	if err != nil {
		t.Fatalf("Error opening: %s", err)
	}

	flt, err := filter.Compile("arp", src.LinkType())
	if err != nil {
		t.Fatalf("Error parsing filter: %s", err)
	}
	defer flt.Cleanup()

	err = src.ApplyFilter(flt)
	if err != nil {
		t.Fatalf("Error applying filter: %s", err)
	}

	var count uint64
	for {
		raw_pkt, err := src.Capture()
		if err != nil {
			t.Fatalf("Error reading: %s %d", err, count)
		}

		if raw_pkt == nil {
			break
		}

		count++
	}

	if count != 2 {
		t.Fatalf("Count mismatch: %d", count)
	}
}

func ExampleCapture() {
	src, err := file.Open("/path/to/file/dump.pcap")
	if err != nil {
		log.Fatal(err)
	}

	// you may configure the source further, e.g. by activating
	// promiscuous mode.

	err = src.Activate()
	if err != nil {
		log.Fatal(err)
	}

	for {
		raw_pkt, err := src.Capture()
		if err != nil {
			log.Fatal(err)
		}

		if raw_pkt == nil {
			break
		}

		log.Println("PACKET!!!")

		// do something with the packet
	}
}

func ExampleInject() {
	dst, err := file.Open("/path/to/file/dump.pcap")
	if err != nil {
		log.Fatal(err)
	}

	// you may configure the source further, e.g. by activating
	// promiscuous mode.

	err = dst.Activate()
	if err != nil {
		log.Fatal(err)
	}

	err = dst.Inject([]byte("random data"))
	if err != nil {
		log.Fatal(err)
	}
}
