package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

const api_base_url = "https://api.r01.ru"

func GetDomains(accessToken string) (Domains, error) {
	url := api_base_url + "/api/v1/domains"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Domains{}, fmt.Errorf("Ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Domains{}, fmt.Errorf("Ошибка при отправке запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return Domains{}, fmt.Errorf("Статус ответа: %v, описание: %s", resp.Status, string(body))
	}

	var result Domains
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		return Domains{}, fmt.Errorf("Ошибка при декодировании ответа: %v", err)
	}

	return result, nil
}

func GetDNSRecords(accessToken string, domainID int) (RRecords, error) {
	url := fmt.Sprintf("%s/api/v1/domains/%d/rrecords", api_base_url, domainID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RRecords{}, fmt.Errorf("Ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return RRecords{}, fmt.Errorf("Ошибка при отправке запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return RRecords{}, fmt.Errorf("Статус ответа: %v, описание: %s", resp.Status, string(body))
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return RRecords{}, fmt.Errorf("Статус ответа: %v, описание: %s", resp.Status, string(body))
	}

	var result RRecords
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		return RRecords{}, fmt.Errorf("Ошибка при декодировании ответа: %v", err)
	}

	return result, nil
}

func AddDNSRecord(accessToken string, domainID int, name string, rtype string, ttl int, data string, info string) (int, error) {
	url := fmt.Sprintf("%s/api/v1/domains/%d/rrecords", api_base_url, domainID)

	body := map[string]interface{}{
		"name": name,
		"type": rtype,
		"ttl":  ttl,
		"data": data,
		"info": info,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return 0, fmt.Errorf("Ошибка формирования JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return 0, fmt.Errorf("Ошибка запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("Ошибка при отправке запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		data, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("API error: %s", data)
	}

	var result AddRecordResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("Ошибка при разборе ответа: %v", err)
	}

	return result.Content.Data.ID, nil
}

func DeleteDNSRecord(accessToken string, domainID int, rrId int) error {
	url := fmt.Sprintf("%s/api/v1/domains/%d/rrecords/%d", api_base_url, domainID, rrId)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("Ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Ошибка при отправке запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ошибка при удалении записи: статус %s", resp.Status)
	}
	return nil
}

func WaitForDNS(rrName, value string, checkInterval time.Duration) {
	log.Printf("Ожидаем появления TXT-записи %s...", rrName)
	for {
		txts, err := net.LookupTXT(rrName)
		if err == nil {
			for _, txt := range txts {
				if txt == value {
					log.Println("TXT-запись найдена.")
					return
				}
			}
		} else {
			log.Printf("Ошибка DNS-запроса: %v", err)
		}

		log.Printf("Запись пока не найдена, повтор через %v...", checkInterval)
		time.Sleep(checkInterval)
	}
}
