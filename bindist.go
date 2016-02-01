package main

import (
      "os"
      "fmt"
      "flag"
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

   fmt.Println(magic1)
   fmt.Println(magic2)
}
