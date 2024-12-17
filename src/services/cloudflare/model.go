package cloudflare

import (
	"bytes"
	"ddnsu/v2/src/global"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const RequestRoot string = "https://api.cloudflare.com/client/v4"

var client = &http.Client{}

func TestToken(token string) (bool, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%v/user/tokens/verify", RequestRoot), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	defer resp.Body.Close()

	var obj map[string]any

	json.NewDecoder(resp.Body).Decode(&obj)

	var validToken bool = obj["success"].(bool)
	return validToken, nil

}

func ReturnZoneIdFromDomain(domain string, token string) (string, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%v/zones", RequestRoot), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, errResp := client.Do(req)

	if errResp != nil {
		fmt.Printf("errResp: %v\n", errResp)
	}

	var obj global.CloudflareZoneResponse
	json.NewDecoder(resp.Body).Decode(&obj)

	defer resp.Body.Close()

	correctId := ""

	for _, result := range obj.Result {
		if result.Name == domain {
			correctId = result.Id
		}
	}

	if correctId == "" {
		return "", fmt.Errorf("could not find id for domain: %v", domain)
	} else {
		return correctId, nil
	}
}

func AddDnsRecord(zoneId string, rtype string, name string, ttl string, comment string, value string, token string) (string, error) {

	var ttlNum int
	if ttl == "auto" {
		ttlNum = 0
	} else {
		ttlNumT, ttlNumErr := strconv.Atoi(ttl)
		ttlNum = ttlNumT

		if ttlNumErr != nil {
			return "", fmt.Errorf("invalid conversion of %v to int", ttl)
		}
	}

	body := map[string]any{
		"type":    rtype,
		"name":    name,
		"ttl":     ttlNum,
		"comment": comment,
		"content": value,
	}
	bodyMarshal, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%v/zones/%v/dns_records", RequestRoot, zoneId), bytes.NewBuffer(bodyMarshal))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, errResp := client.Do(req)

	defer req.Body.Close()

	if errResp != nil {
		return "", fmt.Errorf("attempting to add dns record failed with error: %v", errResp)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("attempting to update dns record failed with message: %v", resp.Body)
	}

	var result global.CloudflareZoneRecordResponseSingle
	decodeErr := json.NewDecoder(resp.Body).Decode(&result)

	if decodeErr != nil {
		return "", fmt.Errorf("error while parsing the body")
	}

	recordId := result.Result.Id
	return recordId, nil
}

func ListDnsRecords(zoneId string, token string) ([]global.DDNSURecord, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%v/zones/%v/dns_records", RequestRoot, zoneId), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, errorResp := client.Do(req)

	if errorResp != nil {
		return []global.DDNSURecord{}, fmt.Errorf("error in attempting to list dns records: %v", errorResp)
	}
	defer resp.Body.Close()

	var dnsRecords global.CloudflareZoneRecordResponseMulti
	decodeErr := json.NewDecoder(resp.Body).Decode(&dnsRecords)
	var records []global.DDNSURecord = make([]global.DDNSURecord, len(dnsRecords.Result))

	if decodeErr != nil {
		return []global.DDNSURecord{}, fmt.Errorf("could not decode the records json: %v", decodeErr)
	}

	for i, record := range dnsRecords.Result {
		// talk about this annoying fucking part
		var fixedName string
		nameSplitArray := strings.Split(record.Name, ".")

		if len(nameSplitArray) == 2 {
			fixedName = "@"
		} else {
			fixedName = nameSplitArray[0]
		}

		var ddnsuRecord global.DDNSURecord = global.DDNSURecord{
			Name:    fixedName,
			Comment: record.Comment,
			Ttl:     record.Ttl,
			Content: record.Content,
			Type:    record.Type,
			Id:      record.Id,
		}
		/*
			fmt.Printf("record: %v\n", record.Name)
			fmt.Printf("type: %v\n", record.Type)
			fmt.Printf("ttl: %v\n", record.Ttl)
			fmt.Printf("value: %v\n", record.Content)
			fmt.Printf("comment: %v\n", record.Comment)
			fmt.Println("-----------------------------")
		*/
		records[i] = ddnsuRecord

	}

	return records, nil
}

func UpdateDnsRecord(recordId string, zoneId string, newvalue string, token string) (string, error) {
	body := map[string]string{
		"content": newvalue,
	}
	bodyMarshal, _ := json.Marshal(body)

	req, _ := http.NewRequest("PATCH", fmt.Sprintf("%v/zones/%v/dns_records/%v", RequestRoot, zoneId, recordId), bytes.NewBuffer(bodyMarshal))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, errResp := client.Do(req)

	defer req.Body.Close()

	if errResp != nil {
		return "", fmt.Errorf("attempting to update dns record failed with error: %v", errResp)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("attempting to update dns record failed with message: %v", resp.Body)
	}

	var dnsRecord global.CloudflareZoneRecordResponseSingle
	json.NewDecoder(resp.Body).Decode(&dnsRecord)

	return "", nil
}

func DeleteDnsRecord(recordId string, zoneId string, token string) (bool, error) {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%v/zones/%v/dns_records/%v", RequestRoot, zoneId, recordId), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, errorResp := client.Do(req)

	if errorResp != nil || resp.StatusCode >= 400 {
		return false, fmt.Errorf("could not delete record: %v. error: %v", recordId, zoneId)
	}

	defer resp.Body.Close()

	return true, nil
}
