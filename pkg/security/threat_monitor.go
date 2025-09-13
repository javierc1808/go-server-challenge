package security

import (
	"net/http"
	"sync"
	"time"
)

// ThreatMonitor monitors threats in real time
type ThreatMonitor struct {
	suspiciousIPs  map[string]*IPThreat
	mutex          sync.RWMutex
	securityLogger *SecurityLogger
}

// IPThreat represents threat information for an IP
type IPThreat struct {
	IP              string
	FirstSeen       time.Time
	LastSeen        time.Time
	RequestCount    int
	SuspiciousCount int
	Blocked         bool
	ThreatLevel     string
}

// NewThreatMonitor creates a new instance of ThreatMonitor
func NewThreatMonitor(securityLogger *SecurityLogger) *ThreatMonitor {
	tm := &ThreatMonitor{
		suspiciousIPs:  make(map[string]*IPThreat),
		securityLogger: securityLogger,
	}

	// Clean up old IPs every hour
	go tm.cleanup()

	return tm
}

// AnalyzeRequest analyzes a request for threats
func (tm *ThreatMonitor) AnalyzeRequest(r *http.Request) bool {
	clientIP := tm.getClientIP(r)

	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Get or create information for the IP
	threat, exists := tm.suspiciousIPs[clientIP]
	if !exists {
		threat = &IPThreat{
			IP:          clientIP,
			FirstSeen:   time.Now(),
			LastSeen:    time.Now(),
			ThreatLevel: "LOW",
		}
		tm.suspiciousIPs[clientIP] = threat
	}

	threat.LastSeen = time.Now()
	threat.RequestCount++

	// Analyze suspicious patterns
	isSuspicious := tm.detectSuspiciousPatterns(r, threat)

	if isSuspicious {
		threat.SuspiciousCount++
		tm.updateThreatLevel(threat)

		// Log the suspicious event
		tm.securityLogger.LogSuspiciousActivity(r, "Suspicious pattern detected", map[string]interface{}{
			"threat_level":     threat.ThreatLevel,
			"request_count":    threat.RequestCount,
			"suspicious_count": threat.SuspiciousCount,
		})
	}

	// Block IPs with high threat
	return threat.ThreatLevel == "HIGH" && threat.Blocked
}

// detectSuspiciousPatterns detects suspicious patterns
func (tm *ThreatMonitor) detectSuspiciousPatterns(r *http.Request, threat *IPThreat) bool {
	suspicious := false

	// Detect suspicious User-Agent
	if tm.isSuspiciousUserAgent(r.UserAgent()) {
		suspicious = true
	}

	// Detect requests to sensitive routes
	if tm.isSensitivePath(r.URL.Path) {
		suspicious = true
	}

	// Detect too many requests in a short time
	if threat.RequestCount > 100 && time.Since(threat.FirstSeen) < 5*time.Minute {
		suspicious = true
	}

	// Detect common attack patterns
	if tm.detectAttackPatterns(r) {
		suspicious = true
	}

	return suspicious
}

// isSuspiciousUserAgent checks if the User-Agent is suspicious
func (tm *ThreatMonitor) isSuspiciousUserAgent(userAgent string) bool {
	suspiciousPatterns := []string{
		"sqlmap",
		"nikto",
		"nmap",
		"masscan",
		"zap",
		"burp",
		"w3af",
		"havij",
		"acunetix",
		"nessus",
		"openvas",
		"wget",
		"curl",
		"python-requests",
		"go-http-client",
	}

	for _, pattern := range suspiciousPatterns {
		if containsIgnoreCase(userAgent, pattern) {
			return true
		}
	}

	return false
}

// isSensitivePath checks if the path is sensitive
func (tm *ThreatMonitor) isSensitivePath(path string) bool {
	sensitivePaths := []string{
		"/admin",
		"/login",
		"/api/auth",
		"/.env",
		"/config",
		"/backup",
		"/database",
		"/phpmyadmin",
		"/wp-admin",
		"/.git",
		"/.svn",
		"/robots.txt",
		"/sitemap.xml",
	}

	for _, sensitivePath := range sensitivePaths {
		if containsIgnoreCase(path, sensitivePath) {
			return true
		}
	}

	return false
}

// detectAttackPatterns detects common attack patterns
func (tm *ThreatMonitor) detectAttackPatterns(r *http.Request) bool {
	// Detect SQL injection attempts
	query := r.URL.RawQuery
	if containsIgnoreCase(query, "union") ||
		containsIgnoreCase(query, "select") ||
		containsIgnoreCase(query, "drop") ||
		containsIgnoreCase(query, "insert") {
		return true
	}

	// Detect XSS attempts
	if containsIgnoreCase(query, "<script") ||
		containsIgnoreCase(query, "javascript:") ||
		containsIgnoreCase(query, "onload=") {
		return true
	}

	// Detect path traversal attempts
	if containsIgnoreCase(query, "../") ||
		containsIgnoreCase(query, "..\\") ||
		containsIgnoreCase(query, "%2e%2e") {
		return true
	}

	return false
}

// updateThreatLevel updates the threat level
func (tm *ThreatMonitor) updateThreatLevel(threat *IPThreat) {
	if threat.SuspiciousCount >= 10 {
		threat.ThreatLevel = "HIGH"
		threat.Blocked = true
	} else if threat.SuspiciousCount >= 5 {
		threat.ThreatLevel = "MEDIUM"
	} else if threat.SuspiciousCount >= 2 {
		threat.ThreatLevel = "LOW"
	}
}

// getClientIP gets the real client IP
func (tm *ThreatMonitor) getClientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return forwarded
	}
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}
	return r.RemoteAddr
}

// cleanup cleans up old IPs
func (tm *ThreatMonitor) cleanup() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		tm.mutex.Lock()
		now := time.Now()
		for ip, threat := range tm.suspiciousIPs {
			if now.Sub(threat.LastSeen) > 24*time.Hour {
				delete(tm.suspiciousIPs, ip)
			}
		}
		tm.mutex.Unlock()
	}
}

// GetThreatStats returns threat statistics
func (tm *ThreatMonitor) GetThreatStats() map[string]interface{} {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	highThreats := 0
	mediumThreats := 0
	lowThreats := 0
	blockedIPs := 0

	for _, threat := range tm.suspiciousIPs {
		switch threat.ThreatLevel {
		case "HIGH":
			highThreats++
		case "MEDIUM":
			mediumThreats++
		case "LOW":
			lowThreats++
		}
		if threat.Blocked {
			blockedIPs++
		}
	}

	return map[string]interface{}{
		"total_suspicious_ips": len(tm.suspiciousIPs),
		"high_threats":         highThreats,
		"medium_threats":       mediumThreats,
		"low_threats":          lowThreats,
		"blocked_ips":          blockedIPs,
	}
}

// containsIgnoreCase checks if a string contains another (case insensitive)
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsSubstring(s, substr)))
}

// containsSubstring checks if s contains substr (case insensitive)
func containsSubstring(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}

	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
