/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  swcMinify: true,
  images: {
    domains: ['localhost', 'owndangan.s3.amazonaws.com'],
  },
  async rewrites() {
    // Proxy backend-relative /uploads/* through this server so images resolve on the frontend origin.
    const base =
      process.env.API_BASE_URL ||
      process.env.NEXT_PUBLIC_API_URL ||
      'http://localhost:8080/api/v1'
    const origin = base.replace(/\/api\/v1\/?$/, '')
    return [{ source: '/uploads/:path*', destination: `${origin}/uploads/:path*` }]
  },
}

module.exports = nextConfig
