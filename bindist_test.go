package main

import "testing"

var ExportContains = contains
var ExportAdd = add 

func TestExportContains(t *testing.T) {
   //
}

func TestAdd(t *testing.T) {
   v := ExportAdd(1,3)
   if v != 3 {
      t.Error("Expected 1.5, got ", v)
   }
}