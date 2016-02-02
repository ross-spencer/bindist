package main

import (
      "os"
      "fmt"
      "flag"
      "regexp"
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

func readFile(fp *os.File, fi os.FileInfo) {
   var eof int64 = fi.Size()
   var pos int64 = 0

   // read file, control how we reach EOF
   for pos < eof {
      fmt.Println("Buffer required: ", getbfsize(pos, fi.Size()))

      b1 := make([]byte, bfsize)

      _, err := fp.Read(b1)
      if err != nil {
         fmt.Println("ERROR: Error reading bytes: ", err)
         break
      }
      //fmt.Printf("%d bytes: %s\n", n1, b1)

      //equivalent to ftell() in C
      pos, _ = fp.Seek(0, os.SEEK_CUR) 
      
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
      readFile(f, fi)
   default: 
      fmt.Println("INFO: Not a file.")
   }

}