package slug_test

import (
	"strings"
	"testing"

	"github.com/owndangan/backend/internal/pkg/slug"
	"github.com/stretchr/testify/require"
)

func TestGenerate_Simple(t *testing.T) {
	slug.ClearCache()
	s := slug.Generate("Wedding Andi dan Siti")
	require.NotEmpty(t, s)
	require.Contains(t, s, "wedding")
	require.Contains(t, s, "andi")
	require.Contains(t, s, "siti")
}

func TestGenerate_Empty(t *testing.T) {
	slug.ClearCache()
	s := slug.Generate("")
	require.NotEmpty(t, s)
	require.Contains(t, s, "undangan")
}

func TestGenerate_SpecialChars(t *testing.T) {
	slug.ClearCache()
	s := slug.Generate("Wedding & Party! @Home")
	require.NotEmpty(t, s)
	require.NotContains(t, s, "&")
	require.NotContains(t, s, "!")
	require.NotContains(t, s, "@")
}

func TestGenerate_Accented(t *testing.T) {
	slug.ClearCache()
	s := slug.Generate("José & François")
	require.NotEmpty(t, s)
	require.Contains(t, s, "jose")
	require.Contains(t, s, "francois")
}

func TestGenerate_Reserved(t *testing.T) {
	slug.ClearCache()
	s := slug.Generate("admin")
	require.NotEmpty(t, s)
	require.True(t, slug.IsReserved("admin"))
}

func TestGenerate_Collision(t *testing.T) {
	slug.ClearCache()
	s1 := slug.Generate("Wedding Andi")
	s2 := slug.Generate("Wedding Andi")
	require.NotEqual(t, s1, s2)
}

func TestGenerate_MaxLength(t *testing.T) {
	slug.ClearCache()
	longName := strings.Repeat("a", 200)
	s := slug.Generate(longName)
	require.LessOrEqual(t, len(s), 103)
}

func TestValidate_Valid(t *testing.T) {
	require.NoError(t, slug.Validate("wedding-andi-siti"))
	require.NoError(t, slug.Validate("my-wedding-2026"))
	require.NoError(t, slug.Validate("abc"))
}

func TestValidate_TooShort(t *testing.T) {
	require.Error(t, slug.Validate("ab"))
}

func TestValidate_TooLong(t *testing.T) {
	require.Error(t, slug.Validate(strings.Repeat("a", 101)))
}

func TestValidate_Reserved(t *testing.T) {
	require.Error(t, slug.Validate("admin"))
	require.Error(t, slug.Validate("api"))
	require.Error(t, slug.Validate("login"))
}

func TestValidate_InvalidChars(t *testing.T) {
	require.Error(t, slug.Validate("wedding_andi"))
	require.Error(t, slug.Validate("wedding.andi"))
	require.Error(t, slug.Validate("WEDDING"))
}

func TestGenerateWithCollisionHandling_CustomSlug(t *testing.T) {
	slug.ClearCache()
	s := slug.GenerateWithCollisionHandling("", "my-custom-wedding")
	require.Equal(t, "my-custom-wedding", s)
}

func TestGenerateWithCollisionHandling_ReservedCustomSlug(t *testing.T) {
	slug.ClearCache()
	s := slug.GenerateWithCollisionHandling("Wedding", "admin")
	require.NotEqual(t, "admin", s)
}

func TestIsReserved(t *testing.T) {
	require.True(t, slug.IsReserved("admin"))
	require.True(t, slug.IsReserved("api"))
	require.True(t, slug.IsReserved("login"))
	require.False(t, slug.IsReserved("wedding-andi"))
}

func TestAddReserved(t *testing.T) {
	slug.AddReserved("my-custom-reserved")
	require.True(t, slug.IsReserved("my-custom-reserved"))
}
