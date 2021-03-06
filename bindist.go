package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

var (
	version = "v2.0.1-beta"
	vers    bool

	magic1 string
	magic2 string
	file   string
	size   bool //bools initialize false
	fname  bool

	byteval1  []byte
	byteval2  []byte
	maxNeedle int

	//window we'll use to search for values
	bfsize = 2048
)

func init() {
	flag.StringVar(&magic1, "magic1", "false", "First magic number in a file to begin from, and offset, e.g. magic,offset.")
	flag.StringVar(&magic2, "magic2", "false", "Second magic number in a file to search for, no offset, e.g. magic.")
	flag.StringVar(&file, "file", "false", "File to find the distance between.")
	flag.BoolVar(&size, "size", false, "[Optional] Return size of file alongsize offset in bytes.")
	flag.BoolVar(&fname, "fname", false, "[Optional] Return filename alongside offset and size.")
	flag.BoolVar(&vers, "version", false, "[Optional] Return version of bindist.")
}

func getDistance(offset1 int, offset2 int) int {
	return (offset2 - offset1) - len(byteval1)
}

func outputResult(found bool, offset1, offset2 int, fi os.FileInfo) {

	//Have reached end of file without finding both sequences :(
	if found == false {
		fmt.Fprintln(os.Stderr, "INFO: Byte sequence one not found in file", fi.Name())
	} else if offset1 == 0 && offset2 == 0 {
		fmt.Fprintln(os.Stderr, "INFO: Byte sequence two not found following byte sequence one", fi.Name())
	} else {
		distance := getDistance(offset1, offset2)
		switch {
		case size && !fname:
			fmt.Fprintf(os.Stdout, "%d, %d\n", distance, fi.Size())
		case size && fname:
			fmt.Fprintf(os.Stdout, "%d, %d, \"%s\"\n", distance, fi.Size(), fi.Name())
		case fname && !size:
			fmt.Fprintf(os.Stdout, "%d, \"%s\"\n", distance, fi.Name())
		default:
			fmt.Fprintln(os.Stdout, distance)
		}
	}
}

func moveWindow(buf []byte, from, to int) (int, []byte) {
	var start int
	if from == 0 && to == 0 {
		start = maxNeedle
		copy(buf[:], buf[bfsize-start:])
	} else {
		start = to - from
		copy(buf[:], buf[from:to])
	}
	return start, buf
}

//return: found, off1, off2, errors
func handleFile(fp *os.File) (bool, int, int, error) {

	var found bool
	var fileoff, offset1, offset2 int
	var start int
	buf := make([]byte, bfsize)

	for {
		dataread, err := fp.Read(buf[start:])
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			return found, 0, 0, err //file read error return
		}

		fileoff += dataread //we'll see how many bytes are read from file

		if !found {
			if off := bytes.Index(buf, byteval1); off >= 0 {
				found = true
				//buf[:elements]  gives us a slice of only used values
				var elementsused = int(start) + dataread //only interested in used slices
				var copyfrom = off + len(byteval1)
				offset1 = fileoff - len(buf[:elementsused]) + off
				start, buf = moveWindow(buf, copyfrom, elementsused)
				continue
			}
		} else {
			if off := bytes.Index(buf, byteval2); off >= 0 {
				//buf[:elements]  gives us a slice of only used values
				elementsused := int(start) + dataread
				offset2 = fileoff - len(buf[:elementsused]) + off
				return found, offset1, offset2, nil
			}
		}
		//must call last to enable last iteration over buffer...
		if err == io.EOF {
			//we haven't returned so far and so we haven't found our values
			return found, 0, 0, nil
		}
		start, buf = moveWindow(buf, 0, 0)
	}
}

//callback for walk needs to match the following:
//type WalkFunc func(path string, info os.FileInfo, err error) error
func readFile(path string, fi os.FileInfo, err error) error {

	f, err := os.Open(path)
	defer f.Close() //closing the file
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1) //should only exit if root is null, consider no-exit
	}

	switch mode := fi.Mode(); {
	case mode.IsRegular():
		found, off1, off2, err := handleFile(f)
		if err == nil {
			outputResult(found, off1, off2, fi)
		}
	case mode.IsDir():
		fmt.Fprintln(os.Stderr, "INFO:", fi.Name(), "is a directory.")
	default:
		fmt.Fprintln(os.Stderr, "INFO: Something completely different.")
	}
	return nil
}

func validateHex(magic string) error {

	/*hex errors to return*/
	const NOTHEX string = "contains invalid hexadecimal characters."
	const UNEVEN string = "contains uneven character count."

	var regexString = "^[A-Fa-f\\d]+$"

	res, _ := regexp.MatchString(regexString, magic)
	if res == false {
		return errors.New(NOTHEX)
	}
	if len(magic)%2 != 0 {
		return errors.New(UNEVEN)
	}
	return nil
}

func getMaxNeedle() {
	// store length of magic1 or magic2, whichever longer.
	//This will be the length of the tail that we copy to the start of each buffer to cover overlaps.
	maxNeedle = len(byteval1)
	if len(byteval2) > maxNeedle {
		maxNeedle = len(byteval2)
	}

	//handle buffer needle size discovered through stress(?) unit tests
	if maxNeedle >= bfsize {
		//almost impossible scenario... os.Exit(1)?
		bfsize = (maxNeedle + 1)
	} else {
		maxNeedle = (maxNeedle - 1)
	}
}

func validateArgsAndGo() {

	err := validateHex(magic1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: magic1 %s \n", err)
		os.Exit(1)
	}

	err = validateHex(magic2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: magic2 %s \n", err)
		os.Exit(1)
	}

	//RL notes: maybe use errors from DecodeString
	//instead of validation above... will try
	byteval1, _ = hex.DecodeString(magic1)
	byteval2, _ = hex.DecodeString(magic2)

	getMaxNeedle()
	filepath.Walk(file, readFile)
}

func main() {
	flag.Parse()
	var verstring = "bindist version"
	if vers {
		fmt.Fprintf(os.Stderr, "%s %s \n", verstring, version)
		os.Exit(0)
	} else if flag.NFlag() <= 2 { // can access args w/ len(os.Args[1:]) too
		fmt.Fprintln(os.Stderr, "Usage:  bindist [-magic1 ...] [-magic2 ...] [-file ...]")
		fmt.Fprintln(os.Stderr, "                [Optional -size] [Optional -fname]")
		fmt.Fprintln(os.Stderr, "                [Optional -version]")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Output: [CSV] 'offset','size','filename'")
		fmt.Fprintf(os.Stderr, "Output: [STRING] '%s ...'\n\n", verstring)
		flag.Usage()
		os.Exit(0)
	}
	validateArgsAndGo()
}
