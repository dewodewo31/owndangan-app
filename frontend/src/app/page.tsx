import Navbar from "@/components/marketing/navbar"
import Hero from "@/components/marketing/hero"
import TrustBar from "@/components/marketing/trust-bar"
import FeaturesSection from "@/components/marketing/features-section"
import TemplateShowcase from "@/components/marketing/template-showcase"
import HowItWorks from "@/components/marketing/how-it-works"
import PricingSection from "@/components/marketing/pricing-section"
import FAQSection from "@/components/marketing/faq-section"
import FinalCTA from "@/components/marketing/final-cta"
import { Footer } from "@/components/marketing/footer"

export const metadata = {
  title: "Owndangan — Undangan Pernikahan Digital yang Elegan",
  description:
    "Buat undangan pernikahan digital yang elegan, kelola tamu, RSVP, galeri, hadiah digital, dan bagikan undangan dengan mudah.",
}

export default function LandingPage() {
  return (
    <main className="min-h-screen">
      <Navbar />
      <Hero />
      <TrustBar />
      <FeaturesSection />
      <TemplateShowcase />
      <HowItWorks />
      <PricingSection />
      <FAQSection />
      <FinalCTA />
      <Footer />
    </main>
  )
}
