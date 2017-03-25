package main

import "fmt"
import "os"
import "testing"

//var ExportContains = contains
var ExportHandleFile = handleFile
var ExportGetDistance = getDistance

type offsetTest struct {
	off1     int
	off2     int
	lenByte1 int
	expected int
}

//TODO: should we test for failures here too?
var offsetTests = []offsetTest{
	{0, 262, 6, 256},
	{0, 262, 2, 260},
	{2, 8, 2, 4},
	{0, 8002, 2, 8000},
	{0, 2, 2, 0},
	{0, 3, 2, 1},
	{0, 4, 2, 2},
	{0, 5, 2, 3},
	{3, 8, 2, 3},
	{0, 14319, 2, 14317},
	{2, 10, 4, 4},
}

type moveWindowsTest struct {
	path      string
	found     bool
	expected1 int
	expected2 int
	err       error
}

var bigLittleTests = []moveWindowsTest{
	{"skeleton-tests/big-little/coffee-ballad-cafe", true, 0, 262, nil},
}

var littleBigTests = []moveWindowsTest{
	{"skeleton-tests/little-big/cafe-coffee-ballad", true, 0, 262, nil},
}

var handleFileTests = []moveWindowsTest{
	{"skeleton-tests/jpg/12-byte-jpg", true, 2, 8, nil},
	{"skeleton-tests/jpg/8000-byte-jpg", true, 0, 8002, nil},
}

var failureTests = []moveWindowsTest{
	{"skeleton-tests/failures/one-only", true, 0, 0, nil},     //first sequence only
	{"skeleton-tests/failures/neither", false, 0, 0, nil},     //neither sequence
	{"skeleton-tests/failures/second-only", false, 0, 0, nil}, //second sequence only
}

var moveWindowTests = []moveWindowsTest{
	{"skeleton-tests/move-window-tests/coffee-one", true, 0, 2, nil},
	{"skeleton-tests/move-window-tests/coffee-two", true, 0, 3, nil},
	{"skeleton-tests/move-window-tests/coffee-three", true, 0, 4, nil},
	{"skeleton-tests/move-window-tests/coffee-four", true, 0, 5, nil},
	{"skeleton-tests/move-window-tests/coffee-five", true, 3, 8, nil},
}

//mock filesystem references (for future):
//https://talks.golang.org/2012/10things.slide#8
//https://github.com/mindreframer/golang-stuff/tree/master/github.com/globocom/tsuru/fs

func runHandleFile(newBfsize int, f *os.File, path string) moveWindowsTest {

	f.Seek(0, 0) //allow multiple uses of the function by resetting file pointer

	bfsize = newBfsize //only as big as the needle+1
	getMaxNeedle()

	found, off1, off2, err := ExportHandleFile(f)

	return moveWindowsTest{path, found, off1, off2, err}
}

func TestExportHandleFile(t *testing.T) {

	//without mocking the filesystem
	//Alternative negative lookup: if _, err := os.Stat("/path/to/whatever"); os.IsNotExist(err)
	if _, err := os.Stat("skeleton-tests/move-window-tests"); err == nil {

		/*Test the window moving contiguously, based on four+ sample files
		  /*incrementing by max needle each time, but no more.*/

		byteval1 = []byte{0xC0, 0x1D}       //cold
		byteval2 = []byte{0xC0, 0xFF, 0xEE} //coffee

		for _, expected := range moveWindowTests {

			f, err := os.Open(expected.path)
			defer f.Close() //closing the file
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR:", err)
				os.Exit(1) //should only exit if root is null, consider no-exit
			}

			var newBfsize = 4 //one bigger than the needle
			actual := runHandleFile(newBfsize, f, expected.path)
			if actual != expected {
				t.Errorf("FAIL: Got offsets, %v, %d, %d, expected, %v %d, %d, bfsize: %d", actual.found, actual.expected1, actual.expected2, expected.found, expected.expected1, expected.expected2, newBfsize)
			}

			newBfsize = 3 //same size as the needle
			actual = runHandleFile(newBfsize, f, expected.path)
			if actual != expected {
				t.Errorf("FAIL: Got offsets, %v, %d, %d, expected, %v %d, %d, bfsize: %d", actual.found, actual.expected1, actual.expected2, expected.found, expected.expected1, expected.expected2, newBfsize)
			}

			newBfsize = 2 //too small for the needle
			actual = runHandleFile(newBfsize, f, expected.path)
			if actual != expected {
				t.Errorf("FAIL: Got offsets, %v, %d, %d, expected, %v %d, %d, bfsize: %d", actual.found, actual.expected1, actual.expected2, expected.found, expected.expected1, expected.expected2, newBfsize)
			}

			newBfsize = 2040 //random buffer size, not base2 not relevant
			actual = runHandleFile(newBfsize, f, expected.path)
			if actual != expected {
				t.Errorf("FAIL: Got offsets, %v, %d, %d, expected, %v %d, %d, bfsize: %d", actual.found, actual.expected1, actual.expected2, expected.found, expected.expected1, expected.expected2, newBfsize)
			}
		}
	}

	if _, err := os.Stat("skeleton-tests/jpg"); err == nil {

		/*Just some JPG files to test with easy enough magix*/

		byteval1 = []byte{0xFF, 0xD8} //cold
		byteval2 = []byte{0xFF, 0xD9} //coffee

		for _, expected := range handleFileTests {

			f, err := os.Open(expected.path)
			defer f.Close() //closing the file
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR:", err)
				os.Exit(1) //should only exit if root is null, consider no-exit
			}

			actual := runHandleFile(bfsize, f, expected.path)
			if actual != expected {
				t.Errorf("FAIL: Got offsets, %v, %d, %d, expected, %v %d, %d, bfsize: %d", actual.found, actual.expected1, actual.expected2, expected.found, expected.expected1, expected.expected2, bfsize)
			}
		}
	}

	if _, err := os.Stat("skeleton-tests/little-big"); err == nil {

		/*Look for a bigger second sequence*/

		byteval1 = []byte{0xCA, 0xFE}                         //cafe
		byteval2 = []byte{0xC0, 0xFF, 0xEE, 0xBA, 0x11, 0xAD} //coffee-ballad

		for _, expected := range littleBigTests {

			f, err := os.Open(expected.path)
			defer f.Close() //closing the file
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR:", err)
				os.Exit(1) //should only exit if root is null, consider no-exit
			}

			var newBfsize = 50
			actual := runHandleFile(newBfsize, f, expected.path)
			if actual != expected {
				t.Errorf("FAIL: Got offsets, %v, %d, %d, expected, %v %d, %d, bfsize: %d", actual.found, actual.expected1, actual.expected2, expected.found, expected.expected1, expected.expected2, newBfsize)
			}
		}
	}

	if _, err := os.Stat("skeleton-tests/big-little"); err == nil {

		/*Look for a smaller second sequence*/

		byteval1 = []byte{0xC0, 0xFF, 0xEE, 0xBA, 0x11, 0xAD} //coffee-ballad
		byteval2 = []byte{0xCA, 0xFE}                         //cafe

		for _, expected := range bigLittleTests {

			f, err := os.Open(expected.path)
			defer f.Close() //closing the file
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR:", err)
				os.Exit(1) //should only exit if root is null, consider no-exit
			}

			var newBfsize = 50
			actual := runHandleFile(newBfsize, f, expected.path)
			if actual != expected {
				t.Errorf("FAIL: Got offsets, %v, %d, %d, expected, %v %d, %d, bfsize: %d", actual.found, actual.expected1, actual.expected2, expected.found, expected.expected1, expected.expected2, newBfsize)
			}
		}
	}

	if _, err := os.Stat("skeleton-tests/failures"); err == nil {

		/*Look for failures*/

		byteval1 = []byte{0xD0, 0x0D, 0x1E} //coffee-ballad
		byteval2 = []byte{0xF1, 0x1E}       //cafe

		for _, expected := range failureTests {

			f, err := os.Open(expected.path)
			defer f.Close() //closing the file
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR:", err)
				os.Exit(1) //should only exit if root is null, consider no-exit
			}

			var newBfsize = 50
			actual := runHandleFile(newBfsize, f, expected.path)
			if actual != expected {
				t.Errorf("FAIL: Got offsets, %v, %d, %d, expected, %v %d, %d, bfsize: %d", actual.found, actual.expected1, actual.expected2, expected.found, expected.expected1, expected.expected2, newBfsize)
			}
		}
	}
}

//create an arbitrary byte slice for byteval1 length calc
func createBytes(length int) []byte {
	buf := make([]byte, length)
	var onebyte byte = 0x00

	for i := 0; i < length; i++ {
		buf[i] = onebyte
	}
	return buf
}

func TestExportGetDistance(t *testing.T) {
	for _, expected := range offsetTests {
		byteval1 = createBytes(expected.lenByte1)
		actual := ExportGetDistance(expected.off1, expected.off2)
		if actual != expected.expected {
			t.Errorf("FAIL: Got distance, %d, expected, %d, off1: %d, off2: %d", actual, expected.expected, expected.off1, expected.off2)
		}
	}
}
