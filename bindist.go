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

func init() {
	flag.StringVar(&magic1, "magic1", "false", "First magic number in a file to begin from, and offset, e.g. magic,offset.")
   flag.StringVar(&magic2, "magic2", "false", "Second magic number in a file to search for, no offset, e.g. magic.")
   flag.StringVar(&file, "file", "false", "File to find the distance between.")
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

}