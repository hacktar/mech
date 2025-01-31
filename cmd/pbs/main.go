package main

import (
   "flag"
   "github.com/89z/mech/pbs"
)

func main() {
   // a
   var address string
   flag.StringVar(&address, "a", "", "address")
   // f
   // pbs.org/wgbh/masterpiece/episodes/downton-abbey-s2-e1
   var bandwidth int
   flag.IntVar(&bandwidth, "f", 2_572_025, "target bandwidth")
   // i
   var info bool
   flag.BoolVar(&info, "i", false, "information")
   // v
   var verbose bool
   flag.BoolVar(&verbose, "v", false, "verbose")
   flag.Parse()
   if verbose {
      pbs.LogLevel = 1
   }
   if address != "" {
      err := doWidget(address, bandwidth, info)
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}
