package main

import "fmt"
import "os"
import "testing"

//var ExportContains = contains
var ExportHandleFile = handleFile

type moveWindowsTest struct {
        path      string
        found     bool
        expected1 int
        expected2 int
}

var moveWindowTests = []moveWindowsTest {
   {"skeleton-tests/move-window-tests/coffee-one", true, 0, 2}, 
   {"skeleton-tests/move-window-tests/coffee-two", true, 0, 3},
   {"skeleton-tests/move-window-tests/coffee-three", true, 0, 4},
   {"skeleton-tests/move-window-tests/coffee-four", true, 0, 5},
   {"skeleton-tests/move-window-tests/coffee-five", true, 3, 8},
}

//mock filesystem references (for future):
//https://talks.golang.org/2012/10things.slide#8
//https://github.com/mindreframer/golang-stuff/tree/master/github.com/globocom/tsuru/fs

func TestExportHandleFile(t *testing.T) {

   //without mocking the filesystem
   //Alternative negative lookup: if _, err := os.Stat("/path/to/whatever"); os.IsNotExist(err)
   if _, err := os.Stat("skeleton-tests/move-window-tests"); err == nil {

      byteval1 = []byte{0xC0, 0x1D}          //cold
      byteval2 = []byte{0xC0, 0xFF, 0xEE}    //coffee
      bfsize = 4

      for _, expected := range moveWindowTests {

         f, err := os.Open(expected.path)
         defer f.Close()   //closing the file
         if err != nil {
            fmt.Fprintln(os.Stderr, "ERROR:", err)
            os.Exit(1)  //should only exit if root is null, consider no-exit
         }

         found, off1, off2, err := ExportHandleFile(f)

         actual := moveWindowsTest{expected.path, found, off1, off2}

         if actual != expected {
            t.Errorf("FAIL: Got offsets, %v, %d, %d, expected, %v %d, %d", actual.found, actual.expected1, actual.expected2, expected.found, expected.expected1, expected.expected2)
         }
      }
   }

   if _, err := os.Stat("skeleton-tests/jpg"); err == nil {
      /*...*/
      // test a range of mock jpg files...
      /*...*/
   }  
}
