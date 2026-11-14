# ImgBB Image Upload Integration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add authenticated, backend-only ImgBB uploads for every existing image URL field without changing the database schema or existing JSON CRUD contracts.

**Architecture:** A shared ImgBB client validates image bytes, uploads via the configured endpoint, and returns the public URL. Resource-specific upload handlers call resource services/repositories to persist that URL after the resource already exists; existing CRUD endpoints remain unchanged and continue to read, update, and soft-delete URL fields. Review uploads replace the stored `photo_urls` array only after all files upload successfully.

**Tech Stack:** Go 1.26, Gin, `net/http`, `mime/multipart`, PostgreSQL via `database/sql`, `httptest`, existing repository/service/handler layering.

---

## File map

- Create `internal/config/image_upload.go`: typed image-upload environment configuration.
- Create `internal/services/imageupload/service.go`: validation and ImgBB HTTP client; no database access.
- Create `internal/services/imageupload/service_test.go`: deterministic HTTP and validation tests using `httptest.Server`.
- Modify `internal/services/event/event_service.go`: upload-and-persist methods for category, event poster/banner, and artist photo.
- Modify `internal/services/community/community_service.go`: upload-and-persist review photos.
- Modify `internal/services/user/user_services.go`: upload-and-persist authenticated profile avatar.
- Modify `internal/services/checkout/checkout_service.go`: upload-and-persist payment method icon.
- Modify the corresponding repositories (`event`, `community`, `user`, `checkout`): narrow URL update methods that return the updated row.
- Modify handlers for event, community, user, checkout: parse multipart files and map service errors to existing API responses.
- Modify `internal/server/routes.go`: register resource-specific authenticated upload routes with current role middleware.
- Keep `internal/config/config.go` unchanged; image configuration has its own typed loader.
- Add/update Swagger annotations in the touched handlers and regenerate `docs/swagger.*` with the repository’s `swag` tool.

## Task 1: Create the shared ImgBB client and configuration

**Files:**
- Create: `internal/config/image_upload.go`
- Create: `internal/services/imageupload/service.go`
- Create: `internal/services/imageupload/service_test.go`

- [ ] **Step 1: Write failing tests for configuration and successful upload**

Define a test server that checks the request is `POST`, receives the API key as the `key` form value, receives the image as the `image` form value, and returns:

```json
{"success":true,"data":{"url":"https://i.ibb.co/example/image.png"}}
```

Test that `Upload(ctx, file)` returns the exact URL and that an empty path/key returns an explicit configuration error.

- [ ] **Step 2: Run the focused tests and verify failure**

Run:

```bash
go test ./internal/services/imageupload -run 'TestUpload|TestConfig' -v
```

Expected: FAIL because the package and client do not exist yet.

- [ ] **Step 3: Implement typed configuration**

Add:

```go
type ImageUploadConfig struct {
    Path   string
    APIKey string
}

func LoadImageUploadConfig() ImageUploadConfig {
    return ImageUploadConfig{
        Path: os.Getenv("IMAGE_UPLOAD_PATH"),
        APIKey: os.Getenv("IMAGE_UPLOAD_API_KEY"),
    }
}
```

- [ ] **Step 4: Implement the client with explicit validation**

Use `http.Client` with a finite timeout. Read the multipart file into memory only after checking the declared size and enforce `5*1024*1024` bytes. Detect MIME from the first bytes using `http.DetectContentType`; accept only JPEG, PNG, WebP, and GIF. Build a multipart request with `key` and `image`, never include the API key in errors, reject non-2xx responses, decode the ImgBB JSON response, and reject a missing/empty `data.url`.

Expose a small interface:

```go
type Uploader interface {
    Upload(context.Context, multipart.FileHeader) (string, error)
}
```

The concrete service may use `multipart.FileHeader.Open()` and must close the opened file explicitly.

- [ ] **Step 5: Add failure-path tests**

Cover invalid MIME, file larger than 5 MiB, empty file, non-2xx ImgBB response, malformed JSON, and a successful JSON response without `data.url`. Assert that invalid local files do not reach the test server.

- [ ] **Step 6: Run tests and commit**

Run:

```bash
gofmt -w internal/config/image_upload.go internal/services/imageupload
go test ./internal/services/imageupload -v
```

Expected: PASS.

Commit:

```bash
git add internal/config/image_upload.go internal/services/imageupload
git commit -m "feat: add ImgBB image uploader"
```

## Task 2: Add resource persistence methods before handlers

**Files:**
- Modify: `internal/repository/event/event_repository.go`
- Modify: `internal/repository/community/community_repository.go`
- Modify: `internal/repository/user/user_repository.go`
- Modify: `internal/repository/checkout/payment_repository.go`
- Modify: `internal/services/event/event_service.go`
- Modify: `internal/services/community/community_service.go`
- Modify: `internal/services/user/user_services.go`
- Modify: `internal/services/checkout/checkout_service.go`

- [ ] **Step 1: Add narrow repository update methods**

Add methods that update only the relevant existing columns and return the normal model:

```go
UpdateCategoryIcon(ctx context.Context, id uuid.UUID, url string) (*models.EventCategory, error)
UpdateEventImages(ctx context.Context, id uuid.UUID, posterURL, bannerURL *string) (*models.Event, error)
UpdateArtistPhoto(ctx context.Context, id uuid.UUID, url string) (*models.Artist, error)
UpdateProfileAvatar(ctx context.Context, userID uuid.UUID, url string) (*models.UserProfile, error)
UpdatePaymentMethodIcon(ctx context.Context, id uuid.UUID, url string) (*models.PaymentMethod, error)
UpdateReviewPhotos(ctx context.Context, userID, reviewID uuid.UUID, urls []string) (*models.Review, error)
```

Use the same `deleted_at IS NULL` predicates and `RETURNING` scan patterns already used in each repository. Do not add a migration.

- [ ] **Step 2: Inject the shared uploader into services**

Extend constructors with an `imageupload.Uploader` dependency and store it on each service. Keep existing constructors’ behavior compatible in tests by using an explicit constructor argument at the single composition root (`internal/server/routes.go`) rather than global state.

- [ ] **Step 3: Add service methods**

Implement methods matching the upload routes:

```go
UploadCategoryIcon(ctx context.Context, id string, file multipart.FileHeader) (*models.EventCategory, error)
UploadEventImages(ctx context.Context, id string, poster, banner *multipart.FileHeader) (*models.Event, error)
UploadArtistPhoto(ctx context.Context, id string, file multipart.FileHeader) (*models.Artist, error)
UploadProfileAvatar(ctx context.Context, userID string, file multipart.FileHeader) (*models.UserProfile, error)
UploadPaymentMethodIcon(ctx context.Context, id string, file multipart.FileHeader) (*models.PaymentMethod, error)
UploadReviewPhotos(ctx context.Context, userID, reviewID string, files []*multipart.FileHeader) (*models.Review, error)
```

Parse UUIDs with the existing helpers. Upload each file first; call the repository only after upload success. For event images, allow either `poster` or `banner` and reject requests with neither. For review photos, upload all files in request order and pass the complete URL array to one repository update.

- [ ] **Step 4: Run formatting and package tests**

Run:

```bash
gofmt -w internal/repository internal/services
go test ./internal/repository/... ./internal/services/...
```

Expected: PASS, or only pre-existing integration failures that are documented without masking new failures.

- [ ] **Step 5: Commit**

```bash
git add internal/repository internal/services
git commit -m "feat: persist uploaded image URLs"
```

## Task 3: Add multipart handlers and authenticated routes

**Files:**
- Modify: `internal/handler/event/event_handler.go`
- Modify: `internal/handler/community/community_handler.go`
- Modify: `internal/handler/user/user_handler.go`
- Modify: `internal/handler/checkout/checkout_handler.go`
- Modify: `internal/server/routes.go`

- [ ] **Step 1: Add handler tests for multipart parsing**

Use `gin` test contexts for handler-only validation paths. Cover missing single-file input and an empty review file list, and assert malformed multipart input returns HTTP 400 without invoking the service. Cover successful upload forwarding in the service/client tests from Task 1 and Task 2 rather than adding a mocking dependency solely for handlers.

- [ ] **Step 2: Implement resource-specific upload handlers**

Use `c.FormFile("file")` for single-image routes. Use `c.MultipartForm()` and read `files["files"]` for review photos. Use `c.FormFile("poster")` and `c.FormFile("banner")` for event images, tolerating one missing field but rejecting both missing. Always return the updated model with the existing `helpers.Success` response shape.

Add Swagger annotations with `@Accept mpfd`, `@Param file formData file true`, and resource-specific path parameters. Do not expose the ImgBB API key or internal error details.

- [ ] **Step 3: Register routes using existing middleware**

Add:

```text
POST /users/me/avatar
POST /event-categories/:category_id/icon
POST /events/:event_id/images
POST /artists/:artist_id/photo
POST /payment-methods/:payment_method_id/icon
POST /community/reviews/:review_id/photos
```

Place each route beside its existing CRUD group. Apply `Auth()` globally to the resource group, and retain the existing organizer/admin role middleware for organizer-managed event, category, artist, and payment method resources. The avatar and review routes remain under the existing authenticated user/community groups.

- [ ] **Step 4: Wire the uploader in the composition root**

In `internal/server/routes.go`, load `config.LoadImageUploadConfig()`, construct one `imageupload.Service`, pass it into the affected services, and keep the existing email/database wiring unchanged.

- [ ] **Step 5: Run handler and route tests**

Run:

```bash
gofmt -w internal/handler internal/server/routes.go
go test ./internal/handler/... ./internal/server/... ./internal/services/...
```

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add internal/handler internal/server/routes.go
git commit -m "feat: expose authenticated image upload routes"
```

## Task 4: Update API documentation and environment example

**Files:**
- Create: `.env.example` with placeholders only; never copy real credentials from `.env`.
- Modify: handler Swagger comments in the files from Task 3.
- Regenerate: `docs/swagger.json`, `docs/swagger.yaml`, and `docs/docs.go` only via the repository’s existing `swag` command.

- [ ] **Step 1: Add non-secret configuration documentation**

Document:

```dotenv
IMAGE_UPLOAD_PATH=https://api.imgbb.com/1/upload
IMAGE_UPLOAD_API_KEY=YOUR_IMGBB_API_KEY
```

Never copy the repository’s real `.env` values into tracked files or test fixtures.

- [ ] **Step 2: Regenerate Swagger**

Run:

```bash
swag init -g cmd/api/main.go -o docs
```

Expected: the upload routes advertise multipart form fields and do not expose secrets.

- [ ] **Step 3: Commit documentation**

```bash
git add .env.example internal/handler docs
git commit -m "docs: document image upload endpoints"
```

## Task 5: Full verification and audit result

**Files:**
- Modify only files needed to fix failures found in Tasks 1-4.

- [ ] **Step 1: Run formatting and compilation**

```bash
gofmt -w internal
go vet ./...
go test ./...
```

Expected: all packages compile and tests pass.

- [ ] **Step 2: Verify secret safety**

```bash
git grep -n "IMAGE_UPLOAD_API_KEY"
git grep -n "8bd0b93214c9b29c724c1c6ea1993e99" || true
```

Expected: only the placeholder/config access is present; no API key is added to tracked source or documentation.

- [ ] **Step 3: Verify the final diff**

```bash
git --no-pager diff --check HEAD~4..HEAD
git status --short
```

Expected: no whitespace errors, only intended tracked changes, and the pre-existing untracked `.commandcode/` directory remains untouched.

- [ ] **Step 4: Update task tracking**

Mark implementation and testing todos done only after the commands above pass. Record any environment-dependent database failures explicitly rather than treating them as successful tests.
