package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

var (
	magic1 string
	magic2 string
	file   string
	size   bool // bools initialize false
	fname  bool

	byteval1  []byte
	byteval2  []byte
	maxNeedle int
	buf       [4096]byte
)

func init() {
	flag.StringVar(&magic1, "magic1", "false", "First magic number in a file to begin from, and offset, e.g. magic,offset.")
	flag.StringVar(&magic2, "magic2", "false", "Second magic number in a file to search for, no offset, e.g. magic.")
	flag.StringVar(&file, "file", "false", "File to find the distance between.")
	flag.BoolVar(&size, "size", false, "[Optional] Return size of file alongsize offset in bytes.")
	flag.BoolVar(&fname, "fname", false, "[Optional] Return filename alongside offset and size.")
}

func handleFile(fp *os.File, fi os.FileInfo) {
	var found bool
	var start, fileoff, offset1, offset2 int
	for {
		i, err := fp.Read(buf[start:])
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			return
		}
		fileoff += i
		if !found {
			// start is the number of bytes we've copied from the tail of the buffer on the previous loop
			// i is the length of the current read. Start + i will normally be the full length of the buffer. Except when we reach EOF.
			if off := bytes.Index(buf[:start+i], byteval1); off >= 0 {
				found = true
				offset1 = fileoff - len(buf[:start+i]) + off
				// copy remainder of buffer before looping in case the sequences are in same buffer
				copy(buf[:], buf[off+len(byteval1):start+i])
				start = len(buf[off+len(byteval1) : start+i])
				// immediately loop again. We may be at io.EOF already here but this is OK, the next read size will just be 0
				continue
			}
		} else {
			if off := bytes.Index(buf[:start+i], byteval2); off >= 0 {
				// Success, print response and return early
				offset2 = fileoff - len(buf[:start+i]) + off
				switch {
				case size && !fname:
					fmt.Fprintf(os.Stdout, "%d, %d\n", (offset2-offset1)-len(byteval1), fi.Size())
				case size && fname:
					fmt.Fprintf(os.Stdout, "%d, %d, \"%s\"\n", (offset2-offset1)-len(byteval1), fi.Size(), fi.Name())
				case fname && !size:
					fmt.Fprintf(os.Stdout, "%d, \"%s\"\n", (offset2-offset1)-len(byteval1), fi.Name())
				default:
					fmt.Fprintln(os.Stdout, (offset2-offset1)-len(byteval1))
				}
				return
			}
		}
		// have reached end of file without finding both sequences :(
		if err == io.EOF {
			if !found {
				fmt.Fprintf(os.Stderr, "INFO: Byte sequence one not found in file %s\n", fi.Name())
			} else {
				fmt.Fprintf(os.Stderr, "INFO: Byte sequence two not found following byte sequence one %s\n", fi.Name())
			}
			return
		}
		// copy the last bit of the buffer to the start so that we can find sequences that overlap the window we are searching
		copy(buf[:], buf[start+i-maxNeedle:start+i])
		start = maxNeedle
	}
}

//callback for walk needs to match the following:
//type WalkFunc func(path string, info os.FileInfo, err error) error
func readFile(path string, fi os.FileInfo, err error) error {
	f, err := os.Open(path)
	defer f.Close() // don't forget to close the file
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1) //should only exit if root is null, consider no-exit
	}
	switch mode := fi.Mode(); {
	case mode.IsRegular():
		handleFile(f, fi)
	case mode.IsDir():
		fmt.Fprintf(os.Stderr, "INFO: %s is a directory.\n", fi.Name())
	default:
		fmt.Fprintln(os.Stderr, "INFO: Something completely different.")
	}
	return nil
}

func main() {
	flag.Parse()
	if flag.NFlag() <= 2 { // can access args w/ len(os.Args[1:]) too
		fmt.Fprintln(os.Stderr, "Usage:  bindist [-magic1 ...] [-magic2 ...] [-file ...]")
		fmt.Fprintln(os.Stderr, "               [Optional -size] [Optional -fname]")
		fmt.Fprintln(os.Stderr, "Output: [CSV] 'offset','size','filename'")
		flag.Usage()
		os.Exit(0)
	}

	res, _ := regexp.MatchString("^[A-Fa-f\\d]+$", magic1)
	if res == false {
		fmt.Fprintln(os.Stderr, "INFO: Magic number one is not hexadecimal.")
		os.Exit(1)
	} else {
		if len(magic1)%2 != 0 {
			fmt.Fprintln(os.Stderr, "INFO: Magic number two contains uneven character count.")
			os.Exit(1)
		}
	}
	res, _ = regexp.MatchString("^[A-Fa-f\\d]+$", magic2)
	if res == false {
		fmt.Fprintln(os.Stderr, "INFO: Magic number two is not hexadecimal.")
		os.Exit(1)
	} else {
		if len(magic2)%2 != 0 {
			fmt.Fprintln(os.Stderr, "INFO: Magic number two contains uneven character count.")
			os.Exit(1)
		}
	}

	byteval1, _ = hex.DecodeString(magic1) // consider just using the errors returned here to test hex validity (rather than validate with regexes)? Would simplify this func
	byteval2, _ = hex.DecodeString(magic2)

	maxNeedle = len(byteval1) // store length of magic1 or magic2, whichever longer. This will be the length of the tail that we copy to the start of each buffer to cover overlaps.
	if len(byteval2) > maxNeedle {
		maxNeedle = len(byteval2)
	}

	filepath.Walk(file, readFile)
}
