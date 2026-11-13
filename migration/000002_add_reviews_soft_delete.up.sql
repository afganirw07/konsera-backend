ALTER TABLE reviews
    ADD COLUMN deleted_at TIMESTAMPTZ;

CREATE INDEX idx_reviews_event_active
    ON reviews (event_id, status)
    WHERE deleted_at IS NULL;
