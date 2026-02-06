package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mhsanaei/3x-ui/v2/logger"
)

// CloudflareService provides integration with Cloudflare DNS API
// for managing DNS records and zones.
type CloudflareService struct {
	SettingService
}

// CloudflareZone represents a Cloudflare DNS zone.
type CloudflareZone struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// CloudflareDNSRecord represents a Cloudflare DNS record.
type CloudflareDNSRecord struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
	Proxied bool   `json:"proxied"`
}

type cfListResponse struct {
	Success bool              `json:"success"`
	Result  []json.RawMessage `json:"result"`
	Errors  []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type cfSingleResponse struct {
	Success bool            `json:"success"`
	Result  json.RawMessage `json:"result"`
	Errors  []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func (s *CloudflareService) doRequest(method, url string, body io.Reader) ([]byte, error) {
	apiToken, err := s.GetCloudflareAPIToken()
	if err != nil || apiToken == "" {
		return nil, fmt.Errorf("cloudflare API token not configured")
	}

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("cloudflare API error (HTTP %d): %s", resp.StatusCode, string(data))
	}

	return data, nil
}

// ListZones retrieves all DNS zones from the Cloudflare account.
func (s *CloudflareService) ListZones() ([]CloudflareZone, error) {
	data, err := s.doRequest("GET", "https://api.cloudflare.com/client/v4/zones?per_page=50", nil)
	if err != nil {
		return nil, err
	}

	var resp cfListResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	if !resp.Success {
		msg := "unknown error"
		if len(resp.Errors) > 0 {
			msg = resp.Errors[0].Message
		}
		return nil, fmt.Errorf("cloudflare API: %s", msg)
	}

	var zones []CloudflareZone
	for _, raw := range resp.Result {
		var z CloudflareZone
		if err := json.Unmarshal(raw, &z); err != nil {
			continue
		}
		zones = append(zones, z)
	}
	return zones, nil
}

// ListDNSRecords retrieves DNS records for a specific zone.
func (s *CloudflareService) ListDNSRecords(zoneID string) ([]CloudflareDNSRecord, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records?per_page=100", zoneID)
	data, err := s.doRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var resp cfListResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	if !resp.Success {
		msg := "unknown error"
		if len(resp.Errors) > 0 {
			msg = resp.Errors[0].Message
		}
		return nil, fmt.Errorf("cloudflare API: %s", msg)
	}

	var records []CloudflareDNSRecord
	for _, raw := range resp.Result {
		var r CloudflareDNSRecord
		if err := json.Unmarshal(raw, &r); err != nil {
			continue
		}
		records = append(records, r)
	}
	return records, nil
}

// CreateDNSRecord creates a new DNS record in the specified zone.
func (s *CloudflareService) CreateDNSRecord(zoneID string, record CloudflareDNSRecord) (*CloudflareDNSRecord, error) {
	payload := map[string]any{
		"type":    record.Type,
		"name":    record.Name,
		"content": record.Content,
		"ttl":     record.TTL,
		"proxied": record.Proxied,
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneID)
	data, err := s.doRequest("POST", url, strings.NewReader(string(bodyBytes)))
	if err != nil {
		return nil, err
	}

	var resp cfSingleResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	if !resp.Success {
		msg := "unknown error"
		if len(resp.Errors) > 0 {
			msg = resp.Errors[0].Message
		}
		return nil, fmt.Errorf("cloudflare API: %s", msg)
	}

	var created CloudflareDNSRecord
	if err := json.Unmarshal(resp.Result, &created); err != nil {
		return nil, err
	}

	logger.Info("Cloudflare DNS record created:", created.Name, created.Type, created.Content)
	return &created, nil
}

// DeleteDNSRecord deletes a DNS record from the specified zone.
func (s *CloudflareService) DeleteDNSRecord(zoneID, recordID string) error {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zoneID, recordID)
	data, err := s.doRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	var resp cfSingleResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}
	if !resp.Success {
		msg := "unknown error"
		if len(resp.Errors) > 0 {
			msg = resp.Errors[0].Message
		}
		return fmt.Errorf("cloudflare API: %s", msg)
	}

	logger.Info("Cloudflare DNS record deleted:", recordID)
	return nil
}
