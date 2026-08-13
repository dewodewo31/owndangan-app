# Production Deployment Checklist

## Pre-Deployment

- [ ] All tests pass (`make test`)
- [ ] Build passes (`make build`)
- [ ] Lint passes (`make lint`)
- [ ] Security scan passes
- [ ] Database migrations reviewed
- [ ] Environment variables configured
- [ ] Secrets configured in GitHub
- [ ] SSL certificate obtained
- [ ] DNS configured

## Infrastructure

- [ ] Server provisioned
- [ ] Docker installed
- [ ] PostgreSQL installed/configured
- [ ] Nginx configured
- [ ] SSL/TLS configured
- [ ] Firewall configured
- [ ] Backup system configured

## Deployment Steps

1. **Database**
   - [ ] Create production database
   - [ ] Run migrations
   - [ ] Configure connection pooling
   - [ ] Set up automated backups

2. **Backend**
   - [ ] Build Docker image
   - [ ] Configure environment variables
   - [ ] Deploy container
   - [ ] Verify health check
   - [ ] Check logs

3. **Frontend**
   - [ ] Build Docker image
   - [ ] Configure API URL
   - [ ] Deploy container
   - [ ] Verify rendering

4. **Reverse Proxy**
   - [ ] Configure Nginx
   - [ ] Set up SSL/TLS
   - [ ] Configure rate limiting
   - [ ] Set up logging

## Post-Deployment

- [ ] Health check passes (`/health`)
- [ ] Authentication works
- [ ] Event creation works
- [ ] Guest management works
- [ ] Payment flow works (sandbox)
- [ ] Webhook reachable
- [ ] Admin panel accessible
- [ ] Logs flowing
- [ ] Monitoring active

## Rollback Plan

If deployment fails:
1. Run `./scripts/rollback.sh production`
2. Verify previous version is running
3. Investigate failure
4. Fix and redeploy

## Monitoring

- Application logs: `docker logs owndangan-backend-prod`
- Database logs: PostgreSQL logs
- Access logs: Nginx logs
- Health check: `curl https://owndangan.com/health`
- Error tracking: Configure Sentry (optional)

## Backup Strategy

- Database: Daily automated backups
- Retention: 30 days
- Storage: Separate from application server
- Test restoration monthly
