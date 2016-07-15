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

func outputResult(found bool, offset1, offset2 int, fi os.FileInfo) {
   
   //Have reached end of file without finding both sequences :(
   if found == false {
      fmt.Fprintln(os.Stderr, "INFO: Byte sequence one not found in file", fi.Name())
   } else if offset1 == 0 && offset2 == 0 {
      fmt.Fprintln(os.Stderr, "INFO: Byte sequence two not found following byte sequence one", fi.Name())
   } else {
      var offset = (offset2-offset1)-len(byteval1)
      switch {
         case size && !fname:
		      fmt.Fprintf(os.Stdout, "%d, %d\n", offset, fi.Size())
         case size && fname:
		      fmt.Fprintf(os.Stdout, "%d, %d, \"%s\"\n", offset, fi.Size(), fi.Name())
         case fname && !size:
            fmt.Fprintf(os.Stdout, "%d, \"%s\"\n", offset, fi.Name())
         default:
            fmt.Fprintln(os.Stdout, offset)
      }
   }
}

func moveWindow(buf []byte, from, to int) (int64, []byte) {
   var start int64
   var nullbuffer, buflen int64    //slice learning todo: delete

   if from == 0 && to == 0 {
      start = int64(maxNeedle)
      buflen = start
      copy(buf[:], buf[bfsize-start:])
   } else {
      start = int64(to - from)
      buflen = start
      copy(buf[:], buf[from:to])
   }

   nullbuffer = bfsize - buflen
   for i := int64(0) ; i < nullbuffer; i++ {
       buf[i+buflen] = 0
   } 

   return start, buf
}

func handleFile(fp *os.File, fi os.FileInfo) {

   var found bool
   var fileoff, offset1, offset2 int  
   var start int64
   buf := make([]byte, bfsize)

   for {
      dataread, err := fp.Read(buf[start:])
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			return
		}
      
      fileoff+=dataread    //we'll see how many bytes are read from file

      if !found {
         if off := bytes.Index(buf, byteval1); off >= 0 {
            found = true
            //buf[:elements]  gives us a slice of only used values
            var elementsused = int(start)+dataread     //only interested in used slices
            var copyfrom = off+len(byteval1)  
            offset1 = fileoff - len(buf[:elementsused]) + off
            start, buf = moveWindow(buf, copyfrom, elementsused)
            continue
         }
      } else {
         if off := bytes.Index(buf, byteval2); off >= 0 {
            //buf[:elements]  gives us a slice of only used values
            elementsused := int(start)+dataread
            offset2 = fileoff - len(buf[:elementsused]) + off
            outputResult(found, offset1, offset2, fi)
            return
         }
      }

      //must call last to enable last iteration over buffer...
      if err == io.EOF {
         //we haven't returned so far and so we haven't found our values
         outputResult(found, 0, 0, fi)
         return
      }

      start, buf = moveWindow(buf, 0, 0)
   }  
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

func getMaxNeedle() {
   maxNeedle = len(byteval1) // store length of magic1 or magic2, whichever longer. This will be the length of the tail that we copy to the start of each buffer to cover overlaps.
   if len(byteval2) > maxNeedle {
      maxNeedle = len(byteval2)
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
