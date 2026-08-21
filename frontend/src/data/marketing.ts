import {
  Camera,
  Calendar,
  ClipboardList,
  Gift,
  Link,
  MessageCircle,
  Palette,
  PenLine,
  Phone,
  Smartphone,
  Users,
  type LucideIcon,
} from "lucide-react"

export const features: { id: string; icon: LucideIcon; title: string; description: string }[] = [
  {
    id: "templates",
    icon: Palette,
    title: "Beautiful Templates",
    description: "Pilih dari berbagai template undangan yang elegan dan modern.",
  },
  {
    id: "editor",
    icon: PenLine,
    title: "Invitation Editor",
    description: "Sesuaikan setiap bagian undangan dengan editor yang mudah digunakan.",
  },
  {
    id: "guests",
    icon: Users,
    title: "Guest Management",
    description: "Kelola seluruh tamu dalam satu tempat dengan mudah.",
  },
  {
    id: "rsvp",
    icon: ClipboardList,
    title: "RSVP",
    description: "Terima konfirmasi kehadiran secara real-time.",
  },
  {
    id: "gallery",
    icon: Camera,
    title: "Gallery",
    description: "Bagikan foto dan momen spesial kalian.",
  },
  {
    id: "digital-gift",
    icon: Gift,
    title: "Digital Gift",
    description: "Terima hadiah digital dengan mudah dan aman.",
  },
  {
    id: "whatsapp",
    icon: Smartphone,
    title: "WhatsApp Sharing",
    description: "Bagikan undangan melalui WhatsApp dengan satu klik.",
  },
  {
    id: "guestbook",
    icon: MessageCircle,
    title: "Guestbook",
    description: "Kumpulkan pesan dan doa dari tamu undangan.",
  },
]

export const templates = [
  {
    id: 1,
    name: "Elegant Gold",
    style: "Elegant",
    image: "/templates/elegant-gold.jpg",
  },
  {
    id: 2,
    name: "Minimal White",
    style: "Minimal",
    image: "/templates/minimal-white.jpg",
  },
  {
    id: 3,
    name: "Modern Navy",
    style: "Modern",
    image: "/templates/modern-navy.jpg",
  },
  {
    id: 4,
    name: "Floral Pink",
    style: "Floral",
    image: "/templates/floral-pink.jpg",
  },
  {
    id: 5,
    name: "Classic Ivory",
    style: "Classic",
    image: "/templates/classic-ivory.jpg",
  },
  {
    id: 6,
    name: "Romantic Rose",
    style: "Elegant",
    image: "/templates/romantic-rose.jpg",
  },
]

export const templateCategories = ["All", "Elegant", "Minimal", "Modern", "Floral", "Classic"]

export const pricingPlans = [
  {
    id: "basic",
    name: "Basic",
    price: 99000,
    duration: "3 bulan",
    description: "Untuk pasangan yang baru memulai",
    features: [
      "1 Undangan",
      "100 Tamu",
      "RSVP Basic",
      "Gallery",
      "Guestbook",
    ],
    cta: "Mulai Sekarang",
    popular: false,
  },
  {
    id: "premium",
    name: "Premium",
    price: 299000,
    duration: "6 bulan",
    description: "Pilihan terbaik untuk pernikahan impian",
    features: [
      "3 Undangan",
      "500 Tamu",
      "RSVP Advanced",
      "Gallery + Video",
      "Digital Gift",
      "WhatsApp Sharing",
    ],
    cta: "Mulai Sekarang",
    popular: true,
  },
  {
    id: "pro",
    name: "Pro",
    price: 350000,
    duration: "Lifetime",
    description: "Untuk unlimited kemungkinan",
    features: [
      "Unlimited Undangan",
      "Unlimited Tamu",
      "Semua Fitur",
      "Custom Domain",
      "Priority Support",
    ],
    cta: "Mulai Sekarang",
    popular: false,
  },
]

export const faqItems = [
  {
    question: "Apa itu Owndangan?",
    answer: "Owndangan adalah platform undangan pernikahan digital yang memungkinkan Anda membuat, mengelola, dan membagikan undangan dengan mudah dan elegan.",
  },
  {
    question: "Apakah saya perlu coding?",
    answer: "Tidak sama sekali. Owndangan dirancang agar mudah digunakan oleh siapa saja tanpa perlu pengetahuan teknis.",
  },
  {
    question: "Apakah undangan bisa dibagikan melalui WhatsApp?",
    answer: "Ya, Owndangan mendukung berbagi undangan melalui WhatsApp dengan satu klik.",
  },
  {
    question: "Apakah tersedia RSVP?",
    answer: "Ya, Owndangan memiliki fitur RSVP yang memungkinkan tamu mengkonfirmasi kehadiran secara real-time.",
  },
  {
    question: "Apakah tersedia digital gift?",
    answer: "Ya, Owndangan mendukung fitur digital gift melalui transfer bank dan QRIS.",
  },
  {
    question: "Berapa lama undangan aktif?",
    answer: "Masa aktif undangan tergantung pada paket yang Anda pilih, mulai dari 3 bulan hingga lifetime.",
  },
  {
    question: "Apakah bisa menggunakan custom slug?",
    answer: "Ya, Anda bisa menggunakan custom slug untuk undangan Anda pada paket Premium dan Pro.",
  },
]

export const problems: { icon: LucideIcon; title: string; description: string }[] = [
  {
    icon: ClipboardList,
    title: "Daftar tamu berantakan",
    description: "Kesulitan mengelola daftar tamu secara manual",
  },
  {
    icon: Phone,
    title: "Konfirmasi RSVP manual",
    description: "Menghubungi satu per satu tamu memakan waktu",
  },
  {
    icon: Link,
    title: "Sulit membagikan undangan",
    description: "Berbagi undangan fisik memakan biaya dan waktu",
  },
  {
    icon: Calendar,
    title: "Informasi acara tersebar",
    description: "Detail acara sulit diakses oleh tamu",
  },
]

export const howItWorks = [
  {
    step: "01",
    title: "Pilih Template",
    description: "Pilih dari berbagai template undangan yang elegan",
  },
  {
    step: "02",
    title: "Isi Detail Pernikahan",
    description: "Tambahkan informasi pasangan, acara, dan lainnya",
  },
  {
    step: "03",
    title: "Bagikan Undangan",
    description: "Bagikan undangan melalui WhatsApp atau link",
  },
]

export const trustIndicators = [
  "Mobile Friendly",
  "RSVP Management",
  "WhatsApp Sharing",
  "Digital Gift",
  "Custom Invitation",
]

// Placeholder business number — replace with the real Owndangan WA number.
export const contactWhatsApp = "https://wa.me/6281234567890"
