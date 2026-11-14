# ImgBB Image Upload Integration Design

## Scope

Audit identified image URL fields already represented in the database and existing CRUD
flows:

- `user_profiles.avatar_url`
- `event_categories.icon_url`
- `events.poster_url` and `events.banner_url`
- `artists.photo_url`
- `payment_methods.icon_url`
- `reviews.photo_urls`

`organizer_verifications.document_url` remains a general URL field and is not changed to
an image upload field.

## Architecture and data flow

Add a centralized ImgBB upload service that reads `IMAGE_UPLOAD_PATH` and
`IMAGE_UPLOAD_API_KEY` from the environment. Resource-specific upload endpoints accept
`multipart/form-data` through the backend, call the shared service, and return the
resulting public URL. Existing JSON CRUD endpoints remain available and continue to
persist URLs, preserving compatibility with existing clients.

Upload routes stay next to their resource routes and use the existing authentication
and role middleware. The service layer coordinates upload and database persistence;
repositories continue to store URLs in the existing columns. No database migration is
required.

For review photos, the upload endpoint accepts multiple files, uploads them in order,
and only persists the resulting URL array after every upload succeeds.

## Validation and error handling

Each file is validated in the backend before any external request:

- maximum size: 5 MiB per file;
- accepted content types: JPEG, PNG, WebP, and GIF;
- type detection uses file content, not only the filename extension;
- empty files are rejected.

Missing configuration, invalid input, ImgBB non-success responses, malformed ImgBB
responses, timeouts, and database errors are returned as explicit API errors. The API
key is never returned or logged.

Create uploads the supplied files before persisting their URLs. Update uploads a new
file before replacing the stored URL, so an upload failure leaves the existing value
unchanged. Reads return the stored URL(s). Deletes retain the repository's existing
soft-delete or delete behavior; remote ImgBB objects are not deleted because this
integration does not have a guaranteed delete operation.

## Testing

Add focused unit tests around the ImgBB client with `httptest.Server` covering:

1. successful upload and URL extraction;
2. MIME and size rejection before an external request;
3. non-success ImgBB responses;
4. malformed responses without a usable URL;
5. configuration and request error handling.

Run formatting, the focused tests, and `go test ./...` after implementation.

## Out of scope

- changing the database schema;
- changing existing JSON CRUD request formats;
- adding upload support for `document_url`;
- changing unrelated features or introducing a new storage abstraction.
