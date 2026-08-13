# Frontend Deployment

## Build

Next.js supports two deployment modes: **Static Export** and **Node.js Server**.

### Static Export (SSG)

Produces a static HTML/JS/CSS bundle that can be served by any web server (Nginx, S3, Cloudflare Pages).

```bash
cd frontend

# Build
npm run build

# Output is in frontend/out/
# Serve with any static file server
```

**Pros**: No server needed, cheap to host, fast CDN distribution.
**Cons**: No API routes, no server-side rendering, no middleware.

### Node.js Server (SSR/ISR)

Runs as a Node.js process with full Next.js features.

```bash
cd frontend

# Build
npm run build

# Start
npm start
```

**Pros**: Full Next.js features (SSR, ISR, API routes, middleware, rewrites).
**Cons**: Requires a Node.js runtime, more expensive to host.

**Recommendation**: Use Node.js server mode for this project (needs API routes for auth and payment integration).

## Environment Variables

All variables must be set at build time (for static export) or runtime (for Node.js server):

```bash
# Build-time vars (baked into the bundle)
NEXT_PUBLIC_MIDTRANS_CLIENT_KEY=...

# Runtime vars (read on the server)
DATABASE_URL=...   # Only if using API routes
```

## Deployment Options

### Option 1: Vercel (Recommended)

```bash
# Install Vercel CLI
npm i -g vercel

# Deploy
vercel --prod
```

Vercel automatically:
- Detects Next.js and applies optimal settings.
- Provides preview deployments for every PR.
- Handles SSL, CDN, and environment variables.
- Supports serverless functions for API routes.

### Option 2: Self-Hosted (Docker)

```dockerfile
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM node:20-alpine AS runner
WORKDIR /app
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/node_modules ./node_modules
EXPOSE 3000
CMD ["npm", "start"]
```

```bash
docker build -t owndangan-frontend:latest -f Dockerfile .
docker run -d --name owndangan-web \
  -p 3000:3000 \
  -e NEXT_PUBLIC_API_URL=https://api.owndangan.com \
  -e NEXT_PUBLIC_MIDTRANS_CLIENT_KEY=... \
  owndangan-frontend:latest
```

### Option 3: Static Export to CDN

```bash
npm run build
# Copy frontend/out/ to S3, Cloudflare R2, or Nginx

# Nginx config
server {
  listen 80;
  server_name owndangan.com;
  root /var/www/owndangan/out;
  index index.html;
  try_files $uri $uri.html $uri/ =404;
}
```

## Post-Deployment

1. Verify the site loads at the deployed URL.
2. Run smoke tests against the production API.
3. Verify Midtrans Snap popup loads correctly.
4. Check that all environment variables are set correctly.
5. Verify SSL certificate is valid.
6. Test the complete user flow (register → login → create invitation).