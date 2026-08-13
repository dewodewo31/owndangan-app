package slug

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"unicode"

	"github.com/google/uuid"
)

var reservedSlugs = map[string]bool{
	"admin": true, "api": true, "auth": true, "login": true, "logout": true,
	"register": true, "user": true, "users": true, "event": true, "events": true,
	"invitation": true, "invitations": true, "wedding": true, "weddings": true,
	"dashboard": true, "settings": true, "profile": true, "help": true, "about": true,
	"contact": true, "support": true, "blog": true, "news": true, "shop": true,
	"store": true, "cart": true, "checkout": true, "payment": true, "payments": true,
	"webhook": true, "health": true, "status": true, "public": true, "private": true,
	"static": true, "assets": true, "images": true, "css": true, "js": true,
	"fonts": true, "uploads": true, "media": true, "files": true, "download": true,
	"upload": true, "search": true, "feed": true, "sitemap": true, "robots": true,
	"favicon": true, "manifest": true, "test": true, "testing": true, "demo": true,
	"example": true, "sample": true, "null": true, "undefined": true, "true": true,
	"false": true, "www": true, "mail": true, "ftp": true, "localhost": true,
	"domain": true, "root": true, "system": true,
}

var (
	slugCache = make(map[string]bool)
	cacheMu   sync.RWMutex
)

var nonAlphanumeric = regexp.MustCompile(`[^a-z0-9-]+`)
var multipleHyphens = regexp.MustCompile(`-+`)
var leadingTrailingHyphens = regexp.MustCompile(`^-+|-+$`)

func Generate(name string) string {
	return GenerateWithCollisionHandling(name, "")
}

func GenerateWithCollisionHandling(name string, existingSlug string) string {
	base := toSlug(name)
	if !isValidSlug(base) || len(base) < 3 {
		base = "undangan"
	}

	if isReserved(base) {
		base = base + "-" + uuid.New().String()[:4]
	}

	candidate := base
	if existingSlug != "" && isValidSlug(existingSlug) && !isReserved(existingSlug) {
		candidate = existingSlug
	}

	candidate = ensureUnique(candidate)
	return candidate
}

func Validate(s string) error {
	if len(s) < 3 {
		return fmt.Errorf("slug must be at least 3 characters")
	}
	if len(s) > 100 {
		return fmt.Errorf("slug must be at most 100 characters")
	}
	if isReserved(s) {
		return fmt.Errorf("slug is reserved")
	}
	if !isValidSlug(s) {
		return fmt.Errorf("slug contains invalid characters")
	}
	return nil
}

func isValidSlug(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-') {
			return false
		}
	}
	return true
}

func isReserved(s string) bool {
	return reservedSlugs[s]
}

func toSlug(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))

	var result strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(toASCII(r))
		} else if r == ' ' || r == '-' || r == '_' {
			result.WriteRune('-')
		}
	}

	s = result.String()
	s = nonAlphanumeric.ReplaceAllString(s, "-")
	s = multipleHyphens.ReplaceAllString(s, "-")
	s = leadingTrailingHyphens.ReplaceAllString(s, "")

	if len(s) < 3 {
		s = "undangan-" + uuid.New().String()[:8]
	}
	if len(s) > 100 {
		s = s[:100]
		s = leadingTrailingHyphens.ReplaceAllString(s, "")
	}

	return s
}

func toASCII(r rune) rune {
	switch r {
	case 'á', 'à', 'ã', 'â', 'ä', 'å':
		return 'a'
	case 'é', 'è', 'ê', 'ë':
		return 'e'
	case 'í', 'ì', 'î', 'ï':
		return 'i'
	case 'ó', 'ò', 'õ', 'ô', 'ö':
		return 'o'
	case 'ú', 'ù', 'û', 'ü':
		return 'u'
	case 'ñ':
		return 'n'
	case 'ç':
		return 'c'
	}
	return r
}

func ensureUnique(candidate string) string {
	cacheMu.RLock()
	exists := slugCache[candidate]
	cacheMu.RUnlock()

	if !exists {
		cacheMu.Lock()
		slugCache[candidate] = true
		cacheMu.Unlock()
		return candidate
	}

	for i := 1; i < 1000; i++ {
		variant := fmt.Sprintf("%s-%d", candidate, i)
		cacheMu.RLock()
		exists := slugCache[variant]
		cacheMu.RUnlock()
		if !exists {
			cacheMu.Lock()
			slugCache[variant] = true
			cacheMu.Unlock()
			return variant
		}
	}

	return candidate + "-" + uuid.New().String()[:8]
}

func ClearCache() {
	cacheMu.Lock()
	slugCache = make(map[string]bool)
	cacheMu.Unlock()
}

func IsReserved(s string) bool {
	return reservedSlugs[s]
}

func AddReserved(slugs ...string) {
	for _, s := range slugs {
		reservedSlugs[s] = true
	}
}
