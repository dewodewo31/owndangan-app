// Source of truth: backend/internal/database/seed.go
// If this test fails, marketing prices drifted from the backend seed — update BOTH.
import { pricingPlans } from "@/data/marketing"

const CANONICAL_SEED_PRICES: Record<string, number> = {
  starter: 99000,
  premium: 299000,
  "all-access": 999000,
}

describe("pricing display vs backend seed", () => {
  it("renders exactly the 3 canonical plans", () => {
    expect(pricingPlans).toHaveLength(3)
    expect(Object.keys(CANONICAL_SEED_PRICES)).toEqual(
      expect.arrayContaining(pricingPlans.map((p) => p.id))
    )
  })

  it("locks every displayed price to its seed value", () => {
    for (const plan of pricingPlans) {
      expect(plan.price).toBe(CANONICAL_SEED_PRICES[plan.id])
    }
  })
})
