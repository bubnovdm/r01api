package internal

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const recordIDFile = "/tmp/certbot_r01_record_id"

func RunAuth(accessToken, domain, validation string) error {
	rrName := "_acme-challenge." + domain

	// Ищем домен
	domains, err := GetDomains(accessToken)
	if err != nil {
		return fmt.Errorf("Ошибка при получении доменов: %w", err)
	}

	var domainID int
	for _, d := range domains.Content.Data {
		if d.Domain == domain || strings.HasSuffix(domain, "."+d.Domain) {
			domainID = d.ID
			break
		}
	}
	if domainID == 0 {
		return fmt.Errorf("Домен %s не найден", domain)
	}

	// Удаляем старые записи
	records, err := GetDNSRecords(accessToken, domainID)
	if err != nil {
		return fmt.Errorf("Ошибка при получении записей: %w", err)
	}
	for _, rec := range records.Content.Data {
		if rec.Name == rrName && rec.Type == "TXT" {
			_ = DeleteDNSRecord(accessToken, domainID, rec.ID)
		}
	}

	// Добавляем новую
	id, err := AddDNSRecord(accessToken, domainID, rrName, "TXT", 300, validation, "letsencrypt-challenge")
	if err != nil {
		return fmt.Errorf("Ошибка при добавлении записи: %w", err)
	}

	if err := os.WriteFile(recordIDFile, []byte(strconv.Itoa(id)), 0644); err != nil {
		log.Printf("Не удалось сохранить ID записи: %v", err)
	}

	// Ждём распространения
	WaitForDNS(rrName, validation, 30*time.Second)

	return nil
}

func RunCleanup(accessToken, domain string) error {
	// Получаем id домена(-ов)
	domains, err := GetDomains(accessToken)
	if err != nil {
		return fmt.Errorf("Ошибка при получении доменов: %w", err)
	}

	var domainID int
	for _, d := range domains.Content.Data {
		if d.Domain == domain || strings.HasSuffix(domain, "."+d.Domain) {
			domainID = d.ID
			break
		}
	}
	if domainID == 0 {
		return fmt.Errorf("Домен %s не найден", domain)
	}

	data, err := os.ReadFile(recordIDFile)
	if err != nil {
		return fmt.Errorf("Не удалось прочитать ID записи: %w", err)
	}
	id, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return fmt.Errorf("Ошибка разбора ID: %w", err)
	}

	if err := DeleteDNSRecord(accessToken, domainID, id); err != nil {
		return fmt.Errorf("Ошибка удаления записи: %w", err)
	}

	_ = os.Remove(recordIDFile)
	return nil
}
