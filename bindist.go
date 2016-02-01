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

func init() {
	flag.StringVar(&magic1, "magic1", "false", "First magic number in a file to begin from, and offset, e.g. magic,offset.")
   flag.StringVar(&magic2, "magic2", "false", "second magic number in a file to search for, no offset, e.g. magic.")
}

func main() {
   flag.Parse()

   if flag.NFlag() <= 1 {    // can access args w/ len(os.Args[1:]) too
      flag.Usage()
      os.Exit(0)
   }

   byteArray1 := []byte(magic1)
   byteArray2 := []byte(magic2)
   fmt.Println(magic1)
   fmt.Println(magic2)
   fmt.Println(byteArray1)
   fmt.Println(byteArray2)
   fmt.Println(hex.EncodeToString(byteArray1))
   fmt.Println(hex.EncodeToString(byteArray2))

   res, _ := regexp.MatchString("^[A-Fa-f\\d]+$", magic1)
   if res == false {
      fmt.Println("INFO: Magic number one is not hexadecimal.")
      os.Exit(1)
   } else {
      if len(magic1) % 2 != 0 {
         fmt.Println("INFO: Magic number two contains uneven character count.")
         os.Exit(1)         
      }
   }

   res, _ = regexp.MatchString("^[A-Fa-f\\d]+$", magic2)
   if res == false {
      fmt.Println("INFO: Magic number two is not hexadecimal.")
      os.Exit(1)
   } else {
      if len(magic2) % 2 != 0 {
         fmt.Println("INFO: Magic number two contains uneven character count.")
         os.Exit(1)         
      }
   }

   var magic1len = len(magic1)
   var x = 0
   for x = 0; x < magic1len; x+=2 {
      fmt.Println(magic1[:2])
      magic1 = magic1[2:]
      //magic1 = append(magic1[:2], magic1[2+1:]...)
   }

}