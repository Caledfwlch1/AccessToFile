package main

import (
	"atf"
	"fmt"
)

const (
	LUN1 = "lun1.bin"
	LUN2 = "lun2.bin"
)

func main() {
	var dev atf.Device

	Logger := atf.CreateLogger()

	srv := atf.NewServer(dev, &atf.ServerConfig{	"172.24.1.3",
						"iqn.2016-04.npp.sit-1920:storage",
						[]string{LUN1, LUN2},
						Logger})

	fmt.Println(srv)
	for i:=0; i<100; i++ {
		fmt.Fprint(srv.LUNs[0].F, "1")
		fmt.Fprint(srv.LUNs[1].F, "2")
	}
	srv.LUNs[0].WriteAt([]byte{0x30, 0x30,0x30,0x30,0x30,0x30}, 50)
	srv.Close()

	return
}
