package service

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/mhsanaei/3x-ui/v2/database"
	"github.com/mhsanaei/3x-ui/v2/database/model"
	"github.com/mhsanaei/3x-ui/v2/logger"

	"github.com/google/uuid"
)

// InboundGeneratorService provides auto-generation of diverse inbound
// configurations for traffic obfuscation, inspired by 3dp-manager.
type InboundGeneratorService struct {
	InboundService
	SettingService
	ServerService
}

// GeneratedInbound holds a generated inbound config ready for creation.
type GeneratedInbound struct {
	Remark         string `json:"remark"`
	Port           int    `json:"port"`
	Protocol       string `json:"protocol"`
	Settings       string `json:"settings"`
	StreamSettings string `json:"streamSettings"`
	Sniffing       string `json:"sniffing"`
	Enable         bool   `json:"enable"`
}

// SniDomainService manages SNI domain whitelist for inbound generation.
type SniDomainService struct{}

// DefaultSniDomains returns the default whitelist of SNI domains.
func DefaultSniDomains() []string {
	return []string{
		"www.google.com",
		"www.microsoft.com",
		"www.apple.com",
		"www.amazon.com",
		"www.cloudflare.com",
		"www.github.com",
		"www.mozilla.org",
		"www.wikipedia.org",
		"www.reddit.com",
		"www.stackoverflow.com",
	}
}

// GetAllDomains returns all SNI domains from the database.
func (s *SniDomainService) GetAllDomains() ([]model.SniDomain, error) {
	db := database.GetDB()
	var domains []model.SniDomain
	err := db.Find(&domains).Error
	return domains, err
}

// AddDomain adds a new SNI domain to the whitelist.
func (s *SniDomainService) AddDomain(name string) (*model.SniDomain, error) {
	db := database.GetDB()
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("domain name cannot be empty")
	}

	var existing model.SniDomain
	err := db.Where("name = ?", name).First(&existing).Error
	if err == nil {
		return &existing, nil
	}

	domain := model.SniDomain{Name: name, Enabled: true}
	if err := db.Create(&domain).Error; err != nil {
		return nil, err
	}
	return &domain, nil
}

// AddDomainsBulk adds multiple SNI domains at once.
func (s *SniDomainService) AddDomainsBulk(names []string) (int, error) {
	db := database.GetDB()
	count := 0
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		var existing model.SniDomain
		if err := db.Where("name = ?", name).First(&existing).Error; err == nil {
			continue
		}
		domain := model.SniDomain{Name: name, Enabled: true}
		if err := db.Create(&domain).Error; err != nil {
			continue
		}
		count++
	}
	return count, nil
}

// DeleteDomain removes an SNI domain from the whitelist.
func (s *SniDomainService) DeleteDomain(id int) error {
	db := database.GetDB()
	return db.Delete(&model.SniDomain{}, id).Error
}

// DeleteAllDomains removes all SNI domains from the whitelist.
func (s *SniDomainService) DeleteAllDomains() error {
	db := database.GetDB()
	return db.Where("1 = 1").Delete(&model.SniDomain{}).Error
}

// SeedDefaultDomains seeds the database with default SNI domains if empty.
func (s *SniDomainService) SeedDefaultDomains() error {
	db := database.GetDB()
	var count int64
	db.Model(&model.SniDomain{}).Count(&count)
	if count > 0 {
		return nil
	}
	for _, name := range DefaultSniDomains() {
		domain := model.SniDomain{Name: name, Enabled: true}
		db.Create(&domain)
	}
	return nil
}

func randomHex(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to a fixed value if crypto/rand fails
		for i := range bytes {
			bytes[i] = byte(i)
		}
	}
	return hex.EncodeToString(bytes)
}

func randomPort(min, max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	if err != nil {
		return min
	}
	return int(n.Int64()) + min
}

func pickDomain(domains []model.SniDomain) string {
	if len(domains) == 0 {
		return "www.google.com"
	}
	enabled := make([]model.SniDomain, 0)
	for _, d := range domains {
		if d.Enabled {
			enabled = append(enabled, d)
		}
	}
	if len(enabled) == 0 {
		return domains[0].Name
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(enabled))))
	if err != nil {
		return enabled[0].Name
	}
	return enabled[n.Int64()].Name
}

func defaultSniffing() string {
	s, _ := json.Marshal(map[string]any{
		"enabled":      false,
		"destOverride": []string{"http", "tls", "quic", "fakedns"},
		"metadataOnly": false,
		"routeOnly":    false,
	})
	return string(s)
}

func buildVlessRealityTCP(port int, clientUUID, domain, privateKey, publicKey string) GeneratedInbound {
	settings, _ := json.Marshal(map[string]any{
		"clients": []map[string]any{{
			"id": clientUUID, "flow": "xtls-rprx-vision", "email": clientUUID,
			"enable": true, "limitIp": 0, "totalGB": 0, "expiryTime": 0, "tgId": "", "subId": "", "reset": 0,
		}},
		"decryption": "none", "encryption": "none", "fallbacks": []any{},
	})
	stream, _ := json.Marshal(map[string]any{
		"network": "tcp", "security": "reality", "externalProxy": []any{},
		"realitySettings": map[string]any{
			"show": false, "xver": 0, "target": domain + ":443", "dest": domain + ":443",
			"serverNames": []string{domain}, "privateKey": privateKey,
			"shortIds": []string{randomHex(4), randomHex(4)},
			"settings": map[string]any{"publicKey": publicKey, "fingerprint": "random", "serverName": "", "spiderX": "/"},
		},
		"tcpSettings": map[string]any{"acceptProxyProtocol": false, "header": map[string]any{"type": "none"}},
	})
	return GeneratedInbound{
		Remark: "vless-tcp-reality", Port: port, Protocol: "vless", Enable: true,
		Settings: string(settings), StreamSettings: string(stream), Sniffing: defaultSniffing(),
	}
}

func buildVlessRealityXHTTP(port int, clientUUID, domain, privateKey, publicKey string) GeneratedInbound {
	settings, _ := json.Marshal(map[string]any{
		"clients": []map[string]any{{
			"id": clientUUID, "flow": "", "email": clientUUID,
			"enable": true, "limitIp": 0, "totalGB": 0, "expiryTime": 0, "tgId": "", "subId": "", "reset": 0,
		}},
		"decryption": "none", "encryption": "none", "fallbacks": []any{},
	})
	stream, _ := json.Marshal(map[string]any{
		"network": "xhttp", "security": "reality", "externalProxy": []any{},
		"realitySettings": map[string]any{
			"show": false, "xver": 0, "target": domain + ":443", "dest": domain + ":443",
			"serverNames": []string{domain}, "privateKey": privateKey,
			"shortIds": []string{randomHex(4), randomHex(4)},
			"settings": map[string]any{"publicKey": publicKey, "fingerprint": "random", "serverName": "", "spiderX": "/"},
		},
		"xhttpSettings": map[string]any{
			"host": domain, "path": "/", "mode": "auto", "noSSEHeader": false,
			"scMaxBufferedPosts": 30, "scMaxEachPostBytes": "1000000",
			"scStreamUpServerSecs": "20-80", "xPaddingBytes": "100-1000",
		},
	})
	return GeneratedInbound{
		Remark: "vless-xhttp-reality", Port: port, Protocol: "vless", Enable: true,
		Settings: string(settings), StreamSettings: string(stream), Sniffing: defaultSniffing(),
	}
}

func buildVlessRealityGRPC(port int, clientUUID, domain, privateKey, publicKey string) GeneratedInbound {
	settings, _ := json.Marshal(map[string]any{
		"clients": []map[string]any{{
			"id": clientUUID, "flow": "", "email": clientUUID,
			"enable": true, "limitIp": 0, "totalGB": 0, "expiryTime": 0, "tgId": "", "subId": "", "reset": 0,
		}},
		"decryption": "none", "encryption": "none", "fallbacks": []any{},
	})
	stream, _ := json.Marshal(map[string]any{
		"network": "grpc", "security": "reality", "externalProxy": []any{},
		"realitySettings": map[string]any{
			"show": false, "xver": 0, "target": domain + ":443", "dest": domain + ":443",
			"serverNames": []string{domain}, "privateKey": privateKey,
			"shortIds": []string{randomHex(4)},
			"settings": map[string]any{"publicKey": publicKey, "fingerprint": "random", "serverName": "", "spiderX": "/"},
		},
		"grpcSettings": map[string]any{
			"serviceName": "myservice", "authority": domain, "multiMode": false,
		},
	})
	return GeneratedInbound{
		Remark: "vless-grpc-reality", Port: port, Protocol: "vless", Enable: true,
		Settings: string(settings), StreamSettings: string(stream), Sniffing: defaultSniffing(),
	}
}

func buildVlessWS(port int, clientUUID, domain string) GeneratedInbound {
	settings, _ := json.Marshal(map[string]any{
		"clients": []map[string]any{{
			"id": clientUUID, "flow": "", "email": clientUUID,
			"enable": true, "limitIp": 0, "totalGB": 0, "expiryTime": 0, "tgId": "", "subId": "", "reset": 0,
		}},
		"decryption": "none", "encryption": "none", "fallbacks": []any{},
	})
	stream, _ := json.Marshal(map[string]any{
		"network": "ws", "security": "none", "externalProxy": []any{},
		"wsSettings": map[string]any{
			"host": domain, "path": "/", "acceptProxyProtocol": false, "heartbeatPeriod": 0,
		},
	})
	return GeneratedInbound{
		Remark: "vless-ws", Port: port, Protocol: "vless", Enable: true,
		Settings: string(settings), StreamSettings: string(stream), Sniffing: defaultSniffing(),
	}
}

func buildVmessTCP(port int, clientUUID string) GeneratedInbound {
	settings, _ := json.Marshal(map[string]any{
		"clients": []map[string]any{{
			"id": clientUUID, "flow": "", "email": clientUUID,
			"enable": true, "limitIp": 0, "totalGB": 0, "expiryTime": 0,
			"tgId": "", "subId": "0", "alterId": "0", "reset": 0,
		}},
	})
	stream, _ := json.Marshal(map[string]any{
		"network": "tcp", "security": "none",
		"tcpSettings": map[string]any{"acceptProxyProtocol": false, "header": map[string]any{"type": "none"}},
	})
	return GeneratedInbound{
		Remark: "vmess-tcp", Port: port, Protocol: "vmess", Enable: true,
		Settings: string(settings), StreamSettings: string(stream), Sniffing: defaultSniffing(),
	}
}

func buildShadowsocksTCP(port int, clientUUID string) GeneratedInbound {
	serverPass := randomHex(32)
	clientPass := randomHex(32)
	settings, _ := json.Marshal(map[string]any{
		"clients": []map[string]any{{
			"id": "", "flow": "", "email": clientUUID, "password": clientPass,
			"enable": true, "limitIp": 0, "totalGB": 0, "expiryTime": 0, "tgId": "", "subId": "", "reset": 0,
		}},
		"ivCheck": false, "method": "2022-blake3-aes-256-gcm", "network": "tcp", "password": serverPass,
	})
	stream, _ := json.Marshal(map[string]any{
		"network": "tcp", "security": "none",
		"tcpSettings": map[string]any{"acceptProxyProtocol": false, "header": map[string]any{"type": "none"}},
	})
	return GeneratedInbound{
		Remark: "ss-tcp", Port: port, Protocol: "shadowsocks", Enable: true,
		Settings: string(settings), StreamSettings: string(stream), Sniffing: defaultSniffing(),
	}
}

func buildTrojanRealityTCP(port int, clientUUID, domain, privateKey, publicKey string) GeneratedInbound {
	trojanPass := randomHex(8)
	settings, _ := json.Marshal(map[string]any{
		"clients": []map[string]any{{
			"id": clientUUID, "email": clientUUID, "password": trojanPass,
			"enable": true, "flow": "", "limitIp": 0, "totalGB": 0, "expiryTime": 0, "tgId": "", "subId": "", "reset": 0,
		}},
		"fallbacks": []any{},
	})
	stream, _ := json.Marshal(map[string]any{
		"network": "tcp", "security": "reality", "externalProxy": []any{},
		"realitySettings": map[string]any{
			"show": false, "xver": 0, "target": domain + ":443", "dest": domain + ":443",
			"serverNames": []string{domain}, "privateKey": privateKey,
			"shortIds": []string{randomHex(4), randomHex(3), randomHex(8), randomHex(2)},
			"settings": map[string]any{"publicKey": publicKey, "fingerprint": "random", "serverName": "", "spiderX": "/"},
		},
		"tcpSettings": map[string]any{"acceptProxyProtocol": false, "header": map[string]any{"type": "none"}},
	})
	return GeneratedInbound{
		Remark: "trojan-tcp-reality", Port: port, Protocol: "trojan", Enable: true,
		Settings: string(settings), StreamSettings: string(stream), Sniffing: defaultSniffing(),
	}
}

// GenerateInbounds generates 10 diverse inbound configurations
// using different protocols, transports, and ports.
func (s *InboundGeneratorService) GenerateInbounds() ([]GeneratedInbound, error) {
	var domainService SniDomainService
	domains, err := domainService.GetAllDomains()
	if err != nil || len(domains) == 0 {
		logger.Warning("No SNI domains found, seeding defaults")
		_ = domainService.SeedDefaultDomains()
		domains, _ = domainService.GetAllDomains()
	}

	// Get X25519 keys for Reality
	keysRaw, err := s.ServerService.GetNewX25519Cert()
	if err != nil {
		return nil, fmt.Errorf("failed to generate X25519 keys: %w", err)
	}
	keys, ok := keysRaw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("unexpected X25519 key format")
	}
	privateKey := fmt.Sprint(keys["privateKey"])
	publicKey := fmt.Sprint(keys["publicKey"])

	usedPorts := make(map[int]bool)
	getPort := func(preferred int) int {
		if preferred > 0 && !usedPorts[preferred] {
			usedPorts[preferred] = true
			return preferred
		}
		for {
			p := randomPort(10000, 60000)
			if !usedPorts[p] {
				usedPorts[p] = true
				return p
			}
		}
	}

	inbounds := []GeneratedInbound{
		buildVlessRealityTCP(getPort(8443), uuid.New().String(), pickDomain(domains), privateKey, publicKey),
		buildVlessRealityXHTTP(getPort(443), uuid.New().String(), pickDomain(domains), privateKey, publicKey),
		buildVlessRealityGRPC(getPort(0), uuid.New().String(), pickDomain(domains), privateKey, publicKey),
		buildVlessWS(getPort(0), uuid.New().String(), pickDomain(domains)),
		buildVlessRealityTCP(getPort(0), uuid.New().String(), pickDomain(domains), privateKey, publicKey),
		buildVlessRealityTCP(getPort(0), uuid.New().String(), pickDomain(domains), privateKey, publicKey),
		buildVlessRealityTCP(getPort(0), uuid.New().String(), pickDomain(domains), privateKey, publicKey),
		buildVmessTCP(getPort(0), uuid.New().String()),
		buildShadowsocksTCP(getPort(0), uuid.New().String()),
		buildTrojanRealityTCP(getPort(0), uuid.New().String(), pickDomain(domains), privateKey, publicKey),
	}

	return inbounds, nil
}

// GenerateAndApplyInbounds generates inbounds and adds them to the panel.
func (s *InboundGeneratorService) GenerateAndApplyInbounds(userID int) ([]*model.Inbound, error) {
	generated, err := s.GenerateInbounds()
	if err != nil {
		return nil, err
	}

	var created []*model.Inbound
	for _, gen := range generated {
		inbound := &model.Inbound{
			UserId:         userID,
			Enable:         gen.Enable,
			Remark:         gen.Remark,
			Port:           gen.Port,
			Protocol:       model.Protocol(gen.Protocol),
			Settings:       gen.Settings,
			StreamSettings: gen.StreamSettings,
			Sniffing:       gen.Sniffing,
			Tag:            fmt.Sprintf("inbound-%v", gen.Port),
		}

		result, _, err := s.InboundService.AddInbound(inbound)
		if err != nil {
			logger.Warning("Failed to create generated inbound:", gen.Remark, err)
			continue
		}
		created = append(created, result)
	}

	logger.Infof("Generated and applied %d inbounds", len(created))
	return created, nil
}
