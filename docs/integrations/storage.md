# Storage Integration

## Overview

The platform uses **S3-compatible object storage** for file uploads. This can be AWS S3, DigitalOcean Spaces, MinIO (self-hosted), or any S3-compatible provider.

## Configuration

| Variable | Description | Example |
|---|---|---|
| `S3_ENDPOINT` | S3-compatible endpoint URL | `https://sgp1.digitaloceanspaces.com` |
| `S3_REGION` | Region (may be ignored by some providers) | `sgp1` |
| `S3_ACCESS_KEY` | Access key ID | `DO00ABC123DEF456` |
| `S3_SECRET_KEY` | Secret access key | `...` |
| `S3_BUCKET` | Bucket name | `owndangan-uploads` |
| `S3_PUBLIC_URL` | Public CDN URL (optional) | `https://uploads.owndangan.com` |

**Never expose `S3_SECRET_KEY` or `S3_ACCESS_KEY` to the frontend.**

## File Upload Handling

### Upload Flow

```
1. Client → Backend: POST /api/v1/upload (multipart form, JWT auth)
2. Backend: Validate file (type, size, magic bytes)
3. Backend: Generate unique filename
4. Backend: Upload to S3 (private by default)
5. Backend: Generate signed URL for access
6. Backend: Store file metadata in DB (original name, S3 key, MIME type, size)
7. Backend → Client: { url: "https://..." } (signed URL)
```

### Implementation

```go
func (s *StorageService) Upload(ctx context.Context, fileHeader *multipart.FileHeader, folder string) (*File, error) {
    // 1. Validate
    if err := validateUpload(fileHeader); err != nil {
        return nil, err
    }

    // 2. Generate unique key
    ext := filepath.Ext(fileHeader.Filename)
    key := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

    // 3. Open file
    file, err := fileHeader.Open()
    if err != nil {
        return nil, fmt.Errorf("open file: %w", err)
    }
    defer file.Close()

    // 4. Upload to S3
    _, err = s.s3Client.PutObject(ctx, &s3.PutObjectInput{
        Bucket:      aws.String(s.bucket),
        Key:         aws.String(key),
        Body:        file,
        ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
        ACL:         aws.String("private"), // Always private, use signed URLs
    })
    if err != nil {
        return nil, fmt.Errorf("s3 upload: %w", err)
    }

    // 5. Store metadata
    fileRecord := &File{
        ID:           uuid.New().String(),
        OriginalName: fileHeader.Filename,
        S3Key:        key,
        MimeType:     fileHeader.Header.Get("Content-Type"),
        Size:         fileHeader.Size,
        Folder:       folder,
        UploadedBy:   GetUserID(ctx),
    }
    s.fileRepo.Create(ctx, fileRecord)

    // 6. Generate signed URL
    fileRecord.URL = s.GenerateSignedURL(ctx, key, 24*time.Hour)

    return fileRecord, nil
}
```

## Signed URLs

### Why Signed URLs
- Files are stored as `private` in S3.
- Access is granted via pre-signed URLs with expiration.
- Prevents hotlinking, unauthorized access, and direct S3 URL guessing.

### Generation

```go
func (s *StorageService) GenerateSignedURL(ctx context.Context, key string, expiry time.Duration) string {
    req, _ := s.s3Client.GetObjectRequest(&s3.GetObjectInput{
        Bucket: aws.String(s.bucket),
        Key:    aws.String(key),
    })
    url, _ := req.Presign(expiry)
    return url
}
```

### Expiration
| Use Case | Expiry | Notes |
|---|---|---|
| Profile/cover photos | 24 hours | Cache on CDN |
| Invitation images | 7 days | Cache on CDN |
| Guest photo uploads | 24 hours | Temporary access |
| Download links | 1 hour | For PDF downloads |

**Long-term**: Files served via CDN with cache headers. Signed URLs are for write/update operations.

## Allowed File Types and Size Limits

### Images (Event galleries, profiles, invitation covers)

| Property | Limit |
|---|---|
| Max file size | 5 MB |
| Allowed MIME types | `image/jpeg`, `image/png`, `image/webp`, `image/gif` |
| Allowed extensions | `.jpg`, `.jpeg`, `.png`, `.webp`, `.gif` |
| Max dimensions | 4096 x 4096 pixels |
| Recommended format | WebP (smaller size, good quality) |

### Documents (Invoice PDFs, backup exports)

| Property | Limit |
|---|---|
| Max file size | 10 MB |
| Allowed MIME types | `application/pdf` |
| Allowed extensions | `.pdf` |

### Music (Background music for invitations)

| Property | Limit |
|---|---|
| Max file size | 10 MB |
| Allowed MIME types | `audio/mpeg`, `audio/ogg`, `audio/aac` |
| Allowed extensions | `.mp3`, `.ogg`, `.aac` |

## Folder Structure

```
s3://owndangan-uploads/
├── users/
│   └── {user_id}/
│       ├── profile/
│       │   └── {uuid}.webp
│       └── cover/
│           └── {uuid}.webp
├── events/
│   └── {event_id}/
│       ├── gallery/
│       │   ├── {uuid}.webp
│       │   └── {uuid}.webp
│       ├── invitations/
│       │   └── {uuid}.pdf
│       └── music/
│           └── {uuid}.mp3
├── templates/
│   └── {template_id}/
│       └── {uuid}.webp
└── temp/
    └── {uuid}.tmp
```

- `temp/` is for temporary files; cleaned up by a cron job every 24 hours.
- All other folders are permanent.
- Files are never moved once uploaded; the folder structure is determined at upload time.

## File Deletion

### On Event Deletion
- Soft-delete event → mark files as `pending_deletion`.
- After 30 days → hard-delete from S3 and DB.
- Use S3 lifecycle rules for automatic cleanup as a backup.

### Delete Implementation

```go
func (s *StorageService) Delete(ctx context.Context, key string) error {
    _, err := s.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
        Bucket: aws.String(s.bucket),
        Key:    aws.String(key),
    })
    return err
}
```

## CDN and Caching

- Use a CDN (e.g., Cloudflare, DigitalOcean CDN) for serving public files.
- Cache-Control: `public, max-age=31536000, immutable` for versioned files.
- For files that may change (profile pictures), use `max-age=86400` with a cache-busting query parameter.
- Signed URLs bypass CDN cache for private files.