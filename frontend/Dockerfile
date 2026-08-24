# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# NEXT_PUBLIC_* vars are baked into the client bundle at build time
ARG NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL

# Server-side proxy destination (next.config.js rewrites); empty → falls back to NEXT_PUBLIC_API_URL
ARG API_BASE_URL=
ENV API_BASE_URL=$API_BASE_URL

# Install dependencies
COPY package*.json ./
RUN npm ci

# Copy source code
COPY . .

# Build
RUN npm run build

# Production stage
FROM node:20-alpine AS production

WORKDIR /app

# Install production dependencies only
COPY package*.json ./
RUN npm ci --only=production

# Copy built assets from builder
COPY --from=builder /app/.next ./.next

# Create non-root user (node:20-alpine already ships uid/gid 1000)
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Change ownership
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000 || exit 1

# Run
CMD ["npm", "start"]
