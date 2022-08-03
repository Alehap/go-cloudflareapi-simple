package main

import (
	"fmt"
	"os"
	"hdapi/cfapi"
)

func main() {
	// cf := cfapi.New("Your API KEY", "Your domain name")
	cf := cfapi.New(os.Getenv("CF_API_KEY"), "dlab.cyou")
	// fmt.Println(cf)
	// cf.GetIdIfExisted()
	status, id := cf.UpdateOrCreateDNSRecord("test.dlab.cyou", "A", "1.2.3.14")
	fmt.Println(status, id)
	fmt.Println(cf.DeleteDNSRecord("test.dlab.cyou", "A"))
// 	idnew := cf.CreateDNSRecord("rmas-1.rke", "A", "1.2.3.4")
// 	fmt.Println(idnew)
}