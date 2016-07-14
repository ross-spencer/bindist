package main

import (
      "io"
      "os"
      "fmt"
      "flag"
      "regexp"
      "errors"
      "bytes"
      "encoding/hex"
      "path/filepath"
   )

var (
   version string = "0.0.1"
   vers bool

	magic1  string
	magic2  string
	file    string
	size    bool    //bools initialize false
	fname   bool

	byteval1  []byte
	byteval2  []byte
	maxNeedle int

   //window we'll use to search for values
   bfsize int64 = 2048
)

func init() {
   flag.StringVar(&magic1, "magic1", "false", "First magic number in a file to begin from, and offset, e.g. magic,offset.")
   flag.StringVar(&magic2, "magic2", "false", "Second magic number in a file to search for, no offset, e.g. magic.")
   flag.StringVar(&file, "file", "false", "File to find the distance between.")
   flag.BoolVar(&size, "size", false, "[Optional] Return size of file alongsize offset in bytes.")
   flag.BoolVar(&fname, "fname", false, "[Optional] Return filename alongside offset and size.")
   flag.BoolVar(&vers, "version", false, "[Optional] Return version of bindist.")
}

func getbfsize(fsize int64, pos int64) int64 {
   newsize := (pos - fsize)
   if newsize > 0 && newsize < bfsize {
      bfsize = newsize
   }
   return bfsize
}

func contains(needle []byte, haystack []byte) (bool, int) {

	nlen := len(needle)
	xlen := len(haystack)

	var offset int

	for x := 0; x < xlen; x += 1 {
		if bytes.Equal(needle, haystack[:nlen]) { //check two slices are equal
			return true, offset
		} else {
			//iterate through haystack comparing two by two...
			haystack = deletefromslice(1, haystack)
		}
		offset += 1
	}

	return false, offset
}

func deletefromslice(n int, slice []byte) []byte {      //return false if no buffer left?
   //Slice Tricks: a = append(a[:i], a[i+1:]...) 
   for x:=0; x<n; x+=1 {
      slice = append(slice[:0], slice[0+1:]...)
   }
   return slice
}

func outputResult(found1, found2 bool, offset1, offset2 int, fi os.FileInfo) {
   if found1 == false {
      fmt.Fprintln(os.Stderr, "INFO: Byte sequence one not found in file", fi.Name())
   } else if found2 == false {
      fmt.Fprintln(os.Stderr, "INFO: Byte sequence two not found following byte sequence one", fi.Name())
   }

   if found1 && found2 {
      var offset = (offset2-offset1)-len(byteval1)
      if size == true && fname == false {
		   fmt.Fprintf(os.Stdout, "%d, %d\n", offset, fi.Size())
      } else if size == true && fname == true {
		   fmt.Fprintf(os.Stdout, "%d, %d, \"%s\"\n", offset, fi.Size(), fi.Name())
      } else if fname == true && size == false {
         fmt.Fprintf(os.Stdout, "%d, \"%s\"\n", offset, fi.Name())
      } else {
         fmt.Fprintln(os.Stdout, offset)
      }
   }
}

func moveWindow(buf []byte, byteval *[]byte) (int64, []byte) {
   start := int64(len(byteval1))
   copy(buf[:], buf[bfsize-start:])
   return start, buf
}

func handleFile(fp *os.File, fi os.FileInfo) {

   var found, found2 bool
   var fileoff, offset1, offset2 int   
   var start int64
   buf := make([]byte, bfsize)

   for {
      i, err := fp.Read(buf[start:])
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			return
		}
      
      fileoff+=i

      //only need one found var for first byteval
      if !found {
         if off := bytes.Index(buf, byteval1); off >= 0 {
            found = true
            offset1 = fileoff - len(buf[:int(start)+i]) + off
            start, buf = moveWindow(buf, &byteval2)
            continue
         }
         start, buf = moveWindow(buf, &byteval1)
      } else {
         if off := bytes.Index(buf, byteval2); off >= 0 {
            found2 = true
            fmt.Println("fileoff", fileoff)
            fmt.Println("index found", off)
            var dataread = int(start) + i
            fmt.Println("data read", dataread)
            fmt.Println("index minus data read", fileoff - int(dataread))
            offset2 = 8 //fileoff - off    //+ dataread
            break
         }
         start, buf = moveWindow(buf, &byteval2)
      }
      //must call last to enable last iteration of stream
      if err == io.EOF {
         //we haven't returned so far and so we haven't found our values
         break
      }
   }
   //204746 abc.jpg
   //4 test.jpg
   outputResult(found, found2, offset1, offset2, fi)  
}

/*func _handleFile(fp *os.File, fi os.FileInfo) {

   var found1, found2 bool
   var tmpoff, offset1, offset2 int
   
   var eof int64 = fi.Size()
   var pos int64 = 0

   // read file, control how we reach EOF
   for pos < eof {
      //fmt.Fprintln(os.Stderr, "Buffer required: ", getbfsize(pos, fi.Size()))

      buf := make([]byte, bfsize)

      _, err := fp.Read(buf)
      if err != nil {
         fmt.Fprintln(os.Stderr, "ERROR: Error reading bytes: ", err)
         break
      }

      if found1 == false {
         found, offset := contains(byteval1, buf)
         tmpoff += offset
         if found == true {
            //we don't need to look for byteval1 any more
            found1 = true
            offset1 = tmpoff
         }
      }

      if found1 == true {
         found, offset := contains(byteval2, buf)
         tmpoff += offset
         if found == true && found2 == false {
            found2 = true
            offset2 = tmpoff
            break
         }
      }

      //equivalent to ftell() in C
      pos, _ = fp.Seek(0, os.SEEK_CUR) 
   }
}*/

//callback for walk needs to match the following:
//type WalkFunc func(path string, info os.FileInfo, err error) error
func readFile (path string, fi os.FileInfo, err error) error {
   
   f, err := os.Open(path)
   defer f.Close()   //closing the file
   if err != nil {
      fmt.Fprintln(os.Stderr, "ERROR:", err)
      os.Exit(1)  //should only exit if root is null, consider no-exit
   }

   switch mode := fi.Mode(); {
   case mode.IsRegular():
      handleFile(f, fi)
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

   var re_string string = "^[A-Fa-f\\d]+$"

   res, _ := regexp.MatchString(re_string, magic)
   if res == false {
      return errors.New(NOTHEX)
   }
   if len(magic) % 2 != 0 {
      return errors.New(UNEVEN)         
   }
   return nil
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

   filepath.Walk(file, readFile)
}

func main() {
   flag.Parse()

   var verstring = "bindist version"

   if vers {
      fmt.Fprintf(os.Stderr, "%s %s \n", verstring, version)
      os.Exit(0)
   } else if flag.NFlag() <= 2 {    // can access args w/ len(os.Args[1:]) too
      fmt.Fprintln(os.Stderr, "Usage:  bindist [-magic1 ...] [-magic2 ...] [-file ...]")
      fmt.Fprintln(os.Stderr, "                [Optional -size] [Optional -fname]")
      fmt.Fprintln(os.Stderr, "                [Optional -version]")
      fmt.Fprintln(os.Stderr, "Output: [CSV] 'offset','size','filename'")
      fmt.Fprintf(os.Stderr, "Output: [STRING] '%s ...'\n", verstring)
      flag.Usage()
      os.Exit(0)
   }

   validateArgsAndGo()
}
