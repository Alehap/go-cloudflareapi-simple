package cfapi

import (
	"fmt"
	"log"
	// "encoding/json"
	"github.com/cloudflare/cloudflare-go"
	"context"
)

type account struct {
    apiKey		string
    zoneId		string
    zoneName	string
    cfAPI		*cloudflare.API
}

func New(apiKey string, zoneName string) account {  
    api, err := cloudflare.NewWithAPIToken(apiKey)
	if err != nil {
		log.Fatal(err)
	}
	// Fetch the zone ID
	zoneId, err := api.ZoneIDByName(zoneName)
	if err != nil {
		log.Fatal(err)
	}

	a := account {apiKey: apiKey, zoneId: zoneId, zoneName: zoneName, cfAPI: api}

    return a
}
func (a account) GetIdIfExisted(recordName string, recordType string) string {
    

	// Fetch all DNS records for example.org
	filter := cloudflare.DNSRecord{Type: recordType, Name: recordName}
	recs, err := a.cfAPI.DNSRecords(context.Background(), a.zoneId, filter)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(recs)
	for _, r := range recs {
		return r.ID
	}
	return ""
}

func (a account) CreateDNSRecord(recordName string, recordType string, recordIp string) (bool,string) {

	priority := uint16(10)
	proxied := false
	unicodeInput := cloudflare.DNSRecord{
		Type:     recordType,
		Name:     recordName,
		Content:  recordIp,
		TTL:      120,
		Priority: &priority,
		Proxied:  &proxied,
	}
	
	records, err := a.cfAPI.CreateDNSRecord(context.Background(), a.zoneId, unicodeInput)
	if err != nil {
		fmt.Println(err)
		return false, err.Error()
	}

	return true, "CreateDNSRecord OK with ID="+records.Result.ID
}

func (a account) UpdateDNSRecord(recordId string, recordName string, recordType string, recordIp string) (bool,string) {

	priority := uint16(10)
	proxied := false
	unicodeInput := cloudflare.DNSRecord{
		Type:     recordType,
		Name:     recordName,
		Content:  recordIp,
		TTL:      120,
		Priority: &priority,
		Proxied:  &proxied,
	}
	
	err := a.cfAPI.UpdateDNSRecord(context.Background(), a.zoneId, recordId, unicodeInput)
	if err != nil {
		fmt.Println(err)
		return false, err.Error()
	}

	return true, "UpdateDNSRecord OK"
}

func (a account) UpdateOrCreateDNSRecord(recordName string, recordType string, recordIp string) (bool,string) {

	rId := a.GetIdIfExisted(recordName, recordType)
	if len(rId) > 0 {
		return a.UpdateDNSRecord(rId, recordName, recordType, recordIp)
	} else {
		return a.CreateDNSRecord(recordName, recordType, recordIp)
	}
}

func (a account) DeleteDNSRecord(recordName string, recordType string) (bool,string) {
	rId := a.GetIdIfExisted(recordName, recordType)
	if len(rId) == 0 {
		return false, "recordName ["+recordName+"] isnot existed."
	}
	err := a.cfAPI.DeleteDNSRecord(context.Background(), a.zoneId, rId)
	if err != nil {
		log.Fatal(err)
		return false, err.Error()
	}


	return true, "DeleteDNSRecord: ["+recordName+"] OK."
}