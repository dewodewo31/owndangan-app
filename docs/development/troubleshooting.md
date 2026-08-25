# Troubleshooting

## Backend

### Server fails to start

**Error**: `listen tcp :8080: bind: address already in use`

**Solution**: Find and kill the process using the port:

```bash
lsof -i :8080
kill -9 <PID>
```

**Error**: `pq: role "postgres" does not exist`

**Solution**: Create the PostgreSQL role, or check the `DB_USER` / `DB_PASSWORD` / `DB_HOST` / `DB_PORT` values in `backend/.env`:

```bash
sudo -u postgres createuser --superuser $USER
# or
createdb owndangan
```

### Database connection issues

**Error**: `dial tcp 127.0.0.1:5432: connect: connection refused`

**Solution**: Ensure PostgreSQL is running:

```bash
sudo systemctl status postgresql    # Linux
pg_isready                          # Check connection
docker start owndangan-db           # If using Docker
```

**Error**: `pq: password authentication failed for user "postgres"`

**Solution**: Update the `DB_PASSWORD` value in `backend/.env` or configure `pg_hba.conf` to use `trust` for local connections.

### Migration errors

Migrations and package seeds run **automatically when the backend starts** (`db.AutoMigrate()` + `SeedPackages`). There is no separate migration command to run.

**Error**: Schema missing / tables don't exist after a fresh DB.

**Solution**: Restart the backend — it auto-migrates on boot. If the DB was created empty, the first server start populates it.

**Error**: `duplicate key value violates unique constraint`

**Solution**: The row already exists (e.g., a seeded package). This is harmless on restart since seeds are idempotent/guarded; otherwise drop and recreate the database.

### CORS errors

**Error**: Frontend gets `No 'Access-Control-Allow-Origin' header` in browser.

**Solution**: Ensure `CORS_ALLOWED_ORIGINS` in backend `.env` includes the frontend URL (e.g., `http://localhost:3000`).

## Frontend

### npm install fails

**Error**: `gyp: No Xcode or CLT version detected` (macOS)

**Solution**: Install Xcode Command Line Tools: `xcode-select --install`.

**Error**: `EACCES: permission denied`

**Solution**: Do not use `sudo`. Use a version manager (nvm, fnm) or reinstall Node without sudo.

### Dev server issues

**Error**: `Module not found: Can't resolve '...'`

**Solution**: Clear Next.js cache and reinstall:

```bash
rm -rf .next node_modules
npm install
```

**Error**: Port 3000 already in use:

```bash
npx kill-port 3000
# or
npm run dev -- --port 3001
```

### API calls fail

**Error**: `401 Unauthorized` on every request.

**Solution**: Your JWT token is missing or expired. Log out and log back in. Check `NEXT_PUBLIC_API_URL` is correct.

**Error**: `TypeError: Failed to fetch` / `NetworkError`.

**Solution**: Ensure the backend is running. Check for CORS issues. Check browser devtools network tab for the actual error.

## General

### Docker

**Error**: `Cannot connect to the Docker daemon`

**Solution**: Ensure Docker is running: `sudo systemctl start docker` or open Docker Desktop.

### Git

**Error**: `fatal: not a git repository`

**Solution**: Run from the project root directory.

### Environment

**Problem**: Changes to `.env` are not picked up.

**Solution**: Restart the backend server and frontend dev server. `.env` files are read at startup.

### Getting Help

If the issue persists:
1. Check the `#backend` or `#frontend` Slack/Discord channel.
2. Search GitHub issues for similar problems.
3. Tag the tech lead if blocking.