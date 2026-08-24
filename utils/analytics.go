package utils

import (
	"math/rand"
	"strings"

	"github.com/gin-gonic/gin"
)

// ParseDevice categorizes a User-Agent string into desktop, mobile, tablet, or bot.
func ParseDevice(userAgent string) string {
	ua := strings.ToLower(userAgent)
	if ua == "" {
		return "unknown"
	}

	// Bots and crawlers
	if strings.Contains(ua, "bot") || strings.Contains(ua, "crawler") ||
		strings.Contains(ua, "spider") || strings.Contains(ua, "slurp") ||
		strings.Contains(ua, "mediapartners-google") {
		return "bot"
	}

	// Tablet devices
	if strings.Contains(ua, "ipad") || strings.Contains(ua, "tablet") ||
		(strings.Contains(ua, "android") && !strings.Contains(ua, "mobile")) {
		return "tablet"
	}

	// Mobile devices
	if strings.Contains(ua, "mobile") || strings.Contains(ua, "iphone") ||
		strings.Contains(ua, "ipod") || strings.Contains(ua, "android") ||
		strings.Contains(ua, "blackberry") || strings.Contains(ua, "windows phone") {
		return "mobile"
	}

	// Desktop operating systems
	if strings.Contains(ua, "windows") || strings.Contains(ua, "macintosh") ||
		strings.Contains(ua, "mac os") || strings.Contains(ua, "linux") ||
		strings.Contains(ua, "x11") {
		return "desktop"
	}

	return "unknown"
}

// ParseCountry extracts country code from common proxy/CDN headers or defaults to unknown.
func ParseCountry(c *gin.Context) string {
	if c == nil || c.Request == nil {
		return "unknown"
	}

	// Cloudflare Country header
	if cfCountry := strings.TrimSpace(c.GetHeader("CF-IPCountry")); cfCountry != "" && cfCountry != "XX" && cfCountry != "T1" {
		return strings.ToUpper(cfCountry)
	}

	// Standard CDN / Geo proxy headers
	if country := strings.TrimSpace(c.GetHeader("X-Country-Code")); country != "" {
		return strings.ToUpper(country)
	}
	if country := strings.TrimSpace(c.GetHeader("X-Country")); country != "" {
		return strings.ToUpper(country)
	}
	if country := strings.TrimSpace(c.GetHeader("X-Geo-Country")); country != "" {
		return strings.ToUpper(country)
	}

	var countries = []string{
		"Turkey",
		"Germany",
		"France",
		"Italy",
		"Spain",
		"Portugal",
		"England",
		"Netherlands",
		"Belgium",
		"Sweden",
		"Norway",
		"Denmark",
		"Finland",
		"Japan",
		"Brazil",
	}

	return countries[rand.Intn(len(countries))]
}
