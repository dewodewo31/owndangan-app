package entitlement

import (
	"encoding/json"

	"github.com/owndangan/backend/internal/model"
)

type Feature string

const (
	GuestMax          Feature = "guest.max"
	GalleryMax        Feature = "gallery.max"
	VideoEnabled      Feature = "video.enabled"
	MusicUpload       Feature = "music.upload"
	MusicPreset       Feature = "music.preset"
	CustomDomain      Feature = "custom_domain"
	WatermarkRemoved  Feature = "watermark.removed"
	WhatsappBulk      Feature = "whatsapp.bulk"
	GuestbookQR       Feature = "guestbook.qr"
	RSVPExport        Feature = "rsvp.export"
	DigitalGiftQRIS   Feature = "digital_gift.qris"
	TemplateCustom    Feature = "template.custom_request"
	EventMax          Feature = "event.max"
)

type Resolver struct {
	pkg *model.Package
}

func NewResolver(pkg *model.Package) *Resolver {
	return &Resolver{pkg: pkg}
}

func (r *Resolver) resolveInt(f Feature) *int {
	if r.pkg == nil || r.pkg.Features == nil {
		return nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(r.pkg.Features, &m); err != nil {
		return nil
	}
	v, ok := m[string(f)]
	if !ok {
		return nil
	}
	switch val := v.(type) {
	case float64:
		i := int(val)
		return &i
	case int:
		return &val
	case int64:
		i := int(val)
		return &i
	}
	return nil
}

func (r *Resolver) resolveBool(f Feature) bool {
	if r.pkg == nil || r.pkg.Features == nil {
		return false
	}
	var m map[string]interface{}
	if err := json.Unmarshal(r.pkg.Features, &m); err != nil {
		return false
	}
	v, ok := m[string(f)]
	if !ok {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val
	case string:
		return val == "true" || val == "1"
	}
	return false
}

func (r *Resolver) resolveString(f Feature) string {
	if r.pkg == nil || r.pkg.Features == nil {
		return ""
	}
	var m map[string]interface{}
	if err := json.Unmarshal(r.pkg.Features, &m); err != nil {
		return ""
	}
	v, ok := m[string(f)]
	if !ok {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func (r *Resolver) GuestMax() *int {
	return r.resolveInt(GuestMax)
}

func (r *Resolver) GalleryMax() *int {
	return r.resolveInt(GalleryMax)
}

func (r *Resolver) EventMax() *int {
	return r.resolveInt(EventMax)
}

func (r *Resolver) VideoEnabled() bool {
	return r.resolveBool(VideoEnabled)
}

func (r *Resolver) MusicUpload() bool {
	return r.resolveBool(MusicUpload)
}

func (r *Resolver) MusicPreset() bool {
	return r.resolveBool(MusicPreset)
}

func (r *Resolver) CustomDomain() bool {
	return r.resolveBool(CustomDomain)
}

func (r *Resolver) WatermarkRemoved() bool {
	return r.resolveBool(WatermarkRemoved)
}

func (r *Resolver) WhatsappBulk() bool {
	return r.resolveBool(WhatsappBulk)
}

func (r *Resolver) GuestbookQR() bool {
	return r.resolveBool(GuestbookQR)
}

func (r *Resolver) RSVPExport() bool {
	return r.resolveBool(RSVPExport)
}

func (r *Resolver) DigitalGiftQRIS() bool {
	return r.resolveBool(DigitalGiftQRIS)
}

func (r *Resolver) TemplateCustom() bool {
	return r.resolveBool(TemplateCustom)
}

func (r *Resolver) CanCreateGuest(currentCount int) bool {
	max := r.GuestMax()
	if max == nil {
		return true
	}
	return currentCount < *max
}

func (r *Resolver) CanCreateEvent(currentCount int) bool {
	max := r.EventMax()
	if max == nil {
		return true
	}
	return currentCount < *max
}

func (r *Resolver) IsUnlimitedGuests() bool {
	return r.GuestMax() == nil
}

func (r *Resolver) IsUnlimitedEvents() bool {
	return r.EventMax() == nil
}

func (r *Resolver) CanAccessPremiumTemplates() bool {
	premium := r.resolveBool("template.premium")
	if premium {
		return true
	}
	group := r.resolveString("template_group")
	return group == "premium" || group == "all"
}

func (r *Resolver) CanAccessAllTemplates() bool {
	all := r.resolveBool("template.all")
	if all {
		return true
	}
	group := r.resolveString("template_group")
	return group == "all"
}
