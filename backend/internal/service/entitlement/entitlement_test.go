package entitlement_test

import (
	"testing"

	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/service/entitlement"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
)

func makePkg(features string) *model.Package {
	return &model.Package{
		Name:     "Test",
		Code:     "test",
		Price:    0,
		Features: datatypes.JSON(features),
	}
}

func TestResolver_FreePackage(t *testing.T) {
	r := entitlement.NewResolver(makePkg(`{"guest.max": 50, "event.max": 1, "video.enabled": false, "custom_domain": false}`))

	require.Equal(t, 50, *r.GuestMax())
	require.False(t, r.VideoEnabled())
	require.False(t, r.CustomDomain())
	require.True(t, r.CanCreateGuest(49))
	require.False(t, r.CanCreateGuest(50))
	require.True(t, r.CanCreateEvent(0))
	require.False(t, r.CanCreateEvent(1))
}

func TestResolver_ProPackage_Unlimited(t *testing.T) {
	r := entitlement.NewResolver(makePkg(`{"guest.max": null, "event.max": null, "video.enabled": true, "custom_domain": true}`))

	require.Nil(t, r.GuestMax())
	require.True(t, r.VideoEnabled())
	require.True(t, r.CustomDomain())
	require.True(t, r.IsUnlimitedGuests())
	require.True(t, r.IsUnlimitedEvents())
	require.True(t, r.CanCreateGuest(99999))
	require.True(t, r.CanCreateEvent(99999))
}

func TestResolver_MissingFeatures(t *testing.T) {
	r := entitlement.NewResolver(makePkg(`{}`))

	require.Nil(t, r.GuestMax())
	require.False(t, r.VideoEnabled())
	require.True(t, r.CanCreateGuest(100))
}

func TestResolver_NilPackage(t *testing.T) {
	r := entitlement.NewResolver(nil)

	require.Nil(t, r.GuestMax())
	require.False(t, r.VideoEnabled())
	require.False(t, r.CustomDomain())
}

func TestResolver_InvalidJSON(t *testing.T) {
	r := entitlement.NewResolver(makePkg(`invalid json`))

	require.Nil(t, r.GuestMax())
	require.False(t, r.VideoEnabled())
}

func TestResolver_AllFeatures(t *testing.T) {
	features := `{
		"guest.max": 200,
		"gallery.max": 20,
		"video.enabled": true,
		"music.upload": true,
		"music.preset": true,
		"custom_domain": false,
		"watermark.removed": true,
		"whatsapp.bulk": true,
		"guestbook.qr": true,
		"rsvp.export": true,
		"digital_gift.qris": false,
		"template.custom_request": true,
		"event.max": 3
	}`
	r := entitlement.NewResolver(makePkg(features))

	require.Equal(t, 200, *r.GuestMax())
	require.Equal(t, 20, *r.GalleryMax())
	require.Equal(t, 3, *r.EventMax())
	require.True(t, r.VideoEnabled())
	require.True(t, r.MusicUpload())
	require.True(t, r.MusicPreset())
	require.False(t, r.CustomDomain())
	require.True(t, r.WatermarkRemoved())
	require.True(t, r.WhatsappBulk())
	require.True(t, r.GuestbookQR())
	require.True(t, r.RSVPExport())
	require.False(t, r.DigitalGiftQRIS())
	require.True(t, r.TemplateCustom())
}

func TestResolver_BooleanStringTrue(t *testing.T) {
	r := entitlement.NewResolver(makePkg(`{"video.enabled": "true"}`))
	require.True(t, r.VideoEnabled())
}

func TestResolver_BooleanStringFalse(t *testing.T) {
	r := entitlement.NewResolver(makePkg(`{"video.enabled": "false"}`))
	require.False(t, r.VideoEnabled())
}
