package main

import (
      "os"
      "fmt"
      "flag"
      "regexp"
      "reflect"
      "encoding/hex"
   )

var magic1 string
var magic2 string
var file string

//window we'll use to search for values
var bfsize int64 = 2048

func init() {
   flag.StringVar(&magic1, "magic1", "false", "First magic number in a file to begin from, and offset, e.g. magic,offset.")
   flag.StringVar(&magic2, "magic2", "false", "Second magic number in a file to search for, no offset, e.g. magic.")
   flag.StringVar(&file, "file", "false", "File to find the distance between.")
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

   var offset int = 0
   var found bool = false

   for x := 0; x < xlen && found == false; x+=1 {
      if reflect.DeepEqual(needle, haystack[:nlen]) {    //check two slices are equal
         found = true
         break
      } else {
         //iterate through haystack comparing two by two...
         haystack = deletefromslice(1, haystack)
      }
      offset+=1
   }

   return found, offset
}

func readFile(fp *os.File, fi os.FileInfo, byteval1 []byte, byteval2 []byte) {
   var eof int64 = fi.Size()
   var pos int64 = 0

   var found1 = false
   //var found2 = false

   var tmpoff int = 0
   var offset1 int = 0

   // read file, control how we reach EOF
   for pos < eof {
      fmt.Println("Buffer required: ", getbfsize(pos, fi.Size()))

      buf := make([]byte, bfsize)

      _, err := fp.Read(buf)
      if err != nil {
         fmt.Println("ERROR: Error reading bytes: ", err)
         break
      }

      found, offset := contains(byteval1, buf)
      tmpoff += offset

      if found == true {
         //we don't need to look for byteval1 any more
         found1 = true
         offset1 = tmpoff
         fmt.Println("xxxx: ", offset1)
         break
      }
      
      //equivalent to ftell() in C
      pos, _ = fp.Seek(0, os.SEEK_CUR) 
   }

   if found1 == false {
      fmt.Println("Error: Byte sequence not found in file.")
   }
}

func main() {
   flag.Parse()

   if flag.NFlag() <= 2 {    // can access args w/ len(os.Args[1:]) too
      flag.Usage()
      os.Exit(0)
   }

   var magic1len = len(magic1)
   var magic2len = len(magic2)

   res, _ := regexp.MatchString("^[A-Fa-f\\d]+$", magic1)
   if res == false {
      fmt.Println("INFO: Magic number one is not hexadecimal.")
      os.Exit(1)
   } else {
      if magic1len % 2 != 0 {
         fmt.Println("INFO: Magic number two contains uneven character count.")
         os.Exit(1)         
      }
   }

   res, _ = regexp.MatchString("^[A-Fa-f\\d]+$", magic2)
   if res == false {
      fmt.Println("INFO: Magic number two is not hexadecimal.")
      os.Exit(1)
   } else {
      if magic2len % 2 != 0 {
         fmt.Println("INFO: Magic number two contains uneven character count.")
         os.Exit(1)         
      }
   }

   byteval1, _ := hex.DecodeString(magic1)
   fmt.Println(byteval1)

   byteval2, _ := hex.DecodeString(magic2)
   fmt.Println(byteval2)

   f, err := os.Open(file)
   if err != nil {
      fmt.Println("ERROR: ", err)
      os.Exit(1)
   }

   fi, err := f.Stat()
   if err != nil {
      fmt.Println("ERROR: ", err)
      os.Exit(1)
   }

   switch mode := fi.Mode(); {
   case mode.IsRegular():
      readFile(f, fi, byteval1, byteval2)
   default: 
      fmt.Println("INFO: Not a file.")
   }

}
