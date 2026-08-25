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
    const rewrites = [
      { source: '/uploads/:path*', destination: `${origin}/uploads/:path*` },
    ]
    // DEV/ngrok only: when API_BASE_URL is set, proxy /api/v1/* to the backend at
    // runtime. This keeps the browser on the frontend (ngrok) origin — same-origin,
    // so no CORS config is needed and ngrok URL rotations require zero rebuilds.
    // In production API_BASE_URL is unset, so this rewrite is omitted and the
    // existing cross-origin behaviour is unchanged.
    if (process.env.API_BASE_URL) {
      rewrites.push({
        source: '/api/v1/:path*',
        destination: `${process.env.API_BASE_URL}/:path*`,
      })
    }
    return rewrites
  },
}

module.exports = nextConfig
