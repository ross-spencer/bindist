package main

import "testing"

var ExportContains = contains

func TestExportContains(t *testing.T) {
   //n = needle, h = haystack
   var h1 = []byte{0x00, 0x00, 0x00, 0xca, 0xfe, 0xba, 0xbe, 0x00}
   var n1 = []byte{0xca, 0xfe, 0xba, 0xbe}

   //small needle 
   var h2 = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xab, 0x00, 0x00, 0x00, 0x00}
   var n2 = []byte{0xab}

   //needle isn't in the haystack
   var h3 = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
   var n3 = []byte{0xff}

   //needle at beginning of file
   var h4 = []byte{0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
   var n4 = []byte{0xff}

   //needle at end of file
   var h5 = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff}
   var n5 = []byte{0xff}

   found1, offset1 := ExportContains(n1,h1)
   if found1 != true && offset1 != 3 {
   	t.Error("Needle not found in haystack when it should have been.")
   }

   found2, offset2 := ExportContains(n2,h2)
   if found2 != true && offset2 != 7 {
   	t.Error("Needle not found in haystack when it should have been.")
   }

   found3, offset3 := ExportContains(n3,h3)
   if found3 != false && offset3 != 0 {
   	t.Error("Needle found in haystack when it shouldn't have been.")
   }

   found4, offset4 := ExportContains(n4,h4)
   if found4 != true && offset4 != 0 {
   	t.Error("Needle not found in haystack when it should have been.")
   }

   found5, offset5 := ExportContains(n5,h5)
   if found5 != true && offset5 != 11 {
   	t.Error("Needle not found in haystack when it should have been.")
   }
}