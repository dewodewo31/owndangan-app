/**
 * Guest personalization helpers.
 *
 * Guest name and token come from the invitation URL, e.g.
 *   /wedding-test-a?to=Budi+dan+Keluarga&token=ABC12345
 *
 * Everything here is rendered as plain text (never dangerouslySetInnerHTML),
 * but we still sanitize aggressively: strip tags, collapse whitespace, cap
 * length — so a malicious ?to= value can never break layout or inject HTML.
 */

export const DEFAULT_GUEST_NAME = "Tamu Undangan"

const MAX_NAME_LENGTH = 120

/** Strip HTML-ish input down to plain text, collapse whitespace, cap length. */
export function sanitizeGuestName(raw: string | null | undefined): string {
  if (!raw) return DEFAULT_GUEST_NAME
  const cleaned = raw
    .replace(/<[^>]*>/g, " ")
    .replace(/[<>"'`]/g, "")
    .replace(/\s+/g, " ")
    .trim()
    .slice(0, MAX_NAME_LENGTH)
  return cleaned || DEFAULT_GUEST_NAME
}

/**
 * Read the guest name from search params.
 * Accepts `to` (the standard param). Falls back to "Tamu Undangan".
 */
export function getGuestName(searchParams: URLSearchParams): string {
  return sanitizeGuestName(searchParams.get("to"))
}

/**
 * Read the guest RSVP token from search params.
 * Guests created by the dashboard carry an 8-char token used to submit RSVP.
 */
export function getGuestToken(searchParams: URLSearchParams): string | null {
  const token = searchParams.get("token")
  return token && /^[A-Za-z0-9]{8}$/.test(token) ? token : null
}

/** Build a shareable URL that preserves the guest param. */
export function buildInvitationUrl(slug: string, guestName?: string): string {
  const base =
    typeof window !== "undefined"
      ? window.location.origin
      : process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000"
  const url = new URL(`/${slug}`, base)
  if (guestName && guestName !== DEFAULT_GUEST_NAME) {
    url.searchParams.set("to", guestName)
  }
  return url.toString()
}
