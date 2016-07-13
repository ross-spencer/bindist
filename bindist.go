package main

import (
      "os"
      "fmt"
      "flag"
      "regexp"
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

func deletefromslice(n int, slice []byte) []byte {      //return false if no buffer left?
   //Slice Tricks: a = append(a[:i], a[i+1:]...) 
   for x:=0; x<n; x+=1 {
      slice = append(slice[:0], slice[0+1:]...)
   }
   return slice
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

func outputResult(found1, found2 bool, offset1, offset2 int, fi os.FileInfo) {
   if found1 == false {
      fmt.Fprintln(os.Stderr, "INFO: Byte sequence one not found in file", fi.Name())
   } else if found2 == false {
      fmt.Fprintln(os.Stderr, "INFO: Byte sequence two not found following byte sequence one", fi.Name())
   }

   if found1 && found2 {
      if size == true && fname == false {
		   fmt.Fprintf(os.Stdout, "%d, %d\n", (offset2-offset1)-len(byteval1), fi.Size())
      } else if size == true && fname == true {
		   fmt.Fprintf(os.Stdout, "%d, %d, \"%s\"\n", (offset2-offset1)-len(byteval1), fi.Size(), fi.Name())
      } else if fname == true && size == false {
         fmt.Fprintf(os.Stdout, "%d, \"%s\"\n", (offset2-offset1)-len(byteval1), fi.Name())
      } else {
         fmt.Fprintln(os.Stdout, (offset2-offset1)-len(byteval1))
      }
   }
}

func handleFile(fp *os.File, fi os.FileInfo) {
   var eof int64 = fi.Size()
   var pos int64 = 0

   var found1 = false
   var found2 = false

   var tmpoff int = 0
   var offset1 int = 0
   var offset2 int = 0

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

   outputResult(found1, found2, offset1, offset2, fi)
}

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

func main() {
   flag.Parse()

   if vers {
      fmt.Fprintf(os.Stderr, "bindist version %s \n", version)
      os.Exit(0)
   } else if flag.NFlag() <= 2 {    // can access args w/ len(os.Args[1:]) too
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
      if len(magic1) % 2 != 0 {
         fmt.Fprintln(os.Stderr, "INFO: Magic number two contains uneven character count.")
         os.Exit(1)         
      }
   }

   res, _ = regexp.MatchString("^[A-Fa-f\\d]+$", magic2)
   if res == false {
      fmt.Fprintln(os.Stderr, "INFO: Magic number two is not hexadecimal.")
      os.Exit(1)
   } else {
      if len(magic2) % 2 != 0 {
         fmt.Fprintln(os.Stderr, "INFO: Magic number two contains uneven character count.")
         os.Exit(1)         
      }
   }

   //RL notes: maybe use errors from DecodeString 
   //instead of validation above... will try
   byteval1, _ = hex.DecodeString(magic1)
   byteval2, _ = hex.DecodeString(magic2)

   filepath.Walk(file, readFile)
}
