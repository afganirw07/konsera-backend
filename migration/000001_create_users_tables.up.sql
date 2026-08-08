-- ============================================================================
-- KONSERA - Database Schema
-- PostgreSQL 15+
-- ============================================================================
-- Run order: extensions -> enums -> tables -> indexes -> triggers -> partitions
-- ============================================================================

-- ----------------------------------------------------------------------------
-- 0. EXTENSIONS
-- ----------------------------------------------------------------------------
CREATE EXTENSION IF NOT EXISTS "pgcrypto";      -- gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS "pg_trgm";       -- fuzzy text search
CREATE EXTENSION IF NOT EXISTS "btree_gin";     -- composite GIN indexes

-- ----------------------------------------------------------------------------
-- 1. ENUM TYPES
-- ----------------------------------------------------------------------------
CREATE TYPE user_status AS ENUM ('active', 'suspended', 'banned', 'pending_verification');
CREATE TYPE organizer_verification_status AS ENUM ('pending', 'in_review', 'approved', 'rejected');
CREATE TYPE event_status AS ENUM ('draft', 'pending_review', 'approved', 'rejected', 'published', 'ongoing', 'completed', 'cancelled');
CREATE TYPE venue_type AS ENUM ('indoor', 'outdoor', 'hybrid');
CREATE TYPE seat_status AS ENUM ('available', 'held', 'booked', 'blocked');
CREATE TYPE ticket_tier_type AS ENUM ('early_bird', 'regular', 'vip', 'vvip', 'group_package');
CREATE TYPE booking_status AS ENUM ('pending', 'awaiting_payment', 'paid', 'expired', 'cancelled', 'refunded', 'partially_refunded');
CREATE TYPE payment_status AS ENUM ('pending', 'processing', 'success', 'failed', 'expired', 'refunded', 'partially_refunded');
CREATE TYPE refund_status AS ENUM ('requested', 'approved', 'rejected', 'processed');
CREATE TYPE ticket_status AS ENUM ('issued', 'transferred', 'checked_in', 'expired', 'cancelled', 'void');
CREATE TYPE transfer_status AS ENUM ('pending', 'accepted', 'rejected', 'cancelled');
CREATE TYPE checkin_result AS ENUM ('success', 'duplicate', 'invalid', 'expired', 'not_yet_open', 'closed');
CREATE TYPE promo_discount_type AS ENUM ('percentage', 'fixed_amount');
CREATE TYPE review_status AS ENUM ('visible', 'hidden', 'reported');
CREATE TYPE payout_status AS ENUM ('scheduled', 'processing', 'paid', 'failed', 'on_hold');
CREATE TYPE audit_action AS ENUM ('create', 'update', 'delete', 'approve', 'reject', 'login', 'payout', 'refund');

-- ----------------------------------------------------------------------------
-- 2. USER MANAGEMENT
-- ----------------------------------------------------------------------------

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(32),
    password VARCHAR(255),              
    auth_provider varchar,
    provider_uid VARCHAR(255),                  
    status user_status NOT NULL DEFAULT 'pending_verification',
    email_verified_at TIMESTAMPTZ,
    phone_verified_at TIMESTAMPTZ,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_users_email UNIQUE (email),
    CONSTRAINT uq_users_phone UNIQUE (phone),
    CONSTRAINT uq_users_provider UNIQUE (provider_uid)
);
CREATE INDEX idx_users_status ON users (status) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_email_trgm ON users USING GIN (email gin_trgm_ops);

CREATE TABLE otp_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    profile_id UUID NOT NULL,
    code VARCHAR(10) NOT NULL,

    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_otp_codes_profile
        FOREIGN KEY (profile_id)
            REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_otp_codes_profile
    ON otp_codes (profile_id);

CREATE INDEX idx_otp_codes_expires_at
    ON otp_codes (expires_at);

CREATE TABLE user_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    full_name VARCHAR(150) NOT NULL,
    avatar_url TEXT,
    date_of_birth DATE,
    gender VARCHAR(20),
    address_line TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(10),
    country VARCHAR(2) DEFAULT 'ID',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user_profiles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT uq_user_profiles_user UNIQUE (user_id)
);
CREATE INDEX idx_user_profiles_city ON user_profiles (city);

CREATE TABLE user_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    favorite_genres TEXT[] DEFAULT '{}',
    notify_push BOOLEAN NOT NULL DEFAULT TRUE,
    notify_email BOOLEAN NOT NULL DEFAULT TRUE,
    notify_sms BOOLEAN NOT NULL DEFAULT FALSE,
    marketing_opt_in BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user_preferences_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT uq_user_preferences_user UNIQUE (user_id)
);
CREATE INDEX idx_user_preferences_genres ON user_preferences USING GIN (favorite_genres);

CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_roles_name UNIQUE (name)
);
-- Seed: super_admin, admin_event, organizer, customer, gatekeeper

CREATE TABLE user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    role_id UUID NOT NULL,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT uq_user_roles UNIQUE (user_id, role_id)
);
CREATE INDEX idx_user_roles_user ON user_roles (user_id);
CREATE INDEX idx_user_roles_role ON user_roles (role_id);

-- ----------------------------------------------------------------------------
-- 3. ORGANIZER
-- ----------------------------------------------------------------------------

CREATE TABLE organizers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    company_name VARCHAR(200) NOT NULL,
    legal_entity_type VARCHAR(50),
    tax_id VARCHAR(50),
    bank_account_name VARCHAR(150),
    bank_account_number VARCHAR(50),
    bank_name VARCHAR(100),
    verification_status organizer_verification_status NOT NULL DEFAULT 'pending',
    commission_rate_override NUMERIC(5,2),      -- overrides platform default if set
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_organizers_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT uq_organizers_user UNIQUE (user_id)
);
CREATE INDEX idx_organizers_status ON organizers (verification_status);

CREATE TABLE organizer_verifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organizer_id UUID NOT NULL,
    document_type VARCHAR(50) NOT NULL,          -- ktp, npwp, siup, akta
    document_url TEXT NOT NULL,
    status organizer_verification_status NOT NULL DEFAULT 'pending',
    reviewed_by UUID,                            -- admin user_id
    reviewed_at TIMESTAMPTZ,
    rejection_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_org_verif_organizer FOREIGN KEY (organizer_id) REFERENCES organizers(id) ON DELETE CASCADE,
    CONSTRAINT fk_org_verif_reviewer FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX idx_org_verif_organizer ON organizer_verifications (organizer_id);

-- ----------------------------------------------------------------------------
-- 4. VENUE
-- ----------------------------------------------------------------------------

CREATE TABLE venues (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    type venue_type NOT NULL,
    address TEXT NOT NULL,
    city VARCHAR(100) NOT NULL,
    province VARCHAR(100),
    country VARCHAR(2) DEFAULT 'ID',
    latitude NUMERIC(9,6),
    longitude NUMERIC(9,6),
    total_capacity INTEGER NOT NULL CHECK (total_capacity > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_venues_city ON venues (city) WHERE deleted_at IS NULL;
CREATE INDEX idx_venues_geo ON venues (latitude, longitude);

CREATE TABLE venue_sections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    venue_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,                  -- e.g. VIP, Festival, Tribune A
    capacity INTEGER NOT NULL CHECK (capacity > 0),
    is_seated BOOLEAN NOT NULL DEFAULT TRUE,      -- false = standing/festival section
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_venue_sections_venue FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE
);
CREATE INDEX idx_venue_sections_venue ON venue_sections (venue_id);

CREATE TABLE seats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    venue_section_id UUID NOT NULL,
    row_label VARCHAR(10) NOT NULL,
    seat_number VARCHAR(10) NOT NULL,
    coord_x NUMERIC(8,2),                         -- for interactive seat map rendering
    coord_y NUMERIC(8,2),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_seats_section FOREIGN KEY (venue_section_id) REFERENCES venue_sections(id) ON DELETE CASCADE,
    CONSTRAINT uq_seats_position UNIQUE (venue_section_id, row_label, seat_number)
);
CREATE INDEX idx_seats_section ON seats (venue_section_id);

-- ----------------------------------------------------------------------------
-- 5. EVENT
-- ----------------------------------------------------------------------------

CREATE TABLE event_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(120) NOT NULL,
    icon_url TEXT,
    CONSTRAINT uq_event_categories_slug UNIQUE (slug)
);

CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organizer_id UUID NOT NULL,
    venue_id UUID NOT NULL,
    title VARCHAR(250) NOT NULL,
    slug VARCHAR(280) NOT NULL,
    description TEXT,
    poster_url TEXT,
    banner_url TEXT,
    status event_status NOT NULL DEFAULT 'draft',
    is_featured BOOLEAN NOT NULL DEFAULT FALSE,
    min_price NUMERIC(14,2),                      -- denormalized for filter/sort
    max_price NUMERIC(14,2),
    terms_and_conditions TEXT,
    refund_policy TEXT,
    approved_by UUID,
    approved_at TIMESTAMPTZ,
    rejection_reason TEXT,
    search_vector TSVECTOR,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_events_organizer FOREIGN KEY (organizer_id) REFERENCES organizers(id) ON DELETE RESTRICT,
    CONSTRAINT fk_events_venue FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE RESTRICT,
    CONSTRAINT fk_events_approver FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT uq_events_slug UNIQUE (slug)
);
CREATE INDEX idx_events_status ON events (status) WHERE deleted_at IS NULL;
CREATE INDEX idx_events_organizer ON events (organizer_id);
CREATE INDEX idx_events_venue ON events (venue_id);
CREATE INDEX idx_events_price ON events (min_price, max_price);
CREATE INDEX idx_events_search ON events USING GIN (search_vector);
CREATE INDEX idx_events_featured ON events (is_featured) WHERE status = 'published';

CREATE TABLE event_category_pivot (
    event_id UUID NOT NULL,
    category_id UUID NOT NULL,
    PRIMARY KEY (event_id, category_id),
    CONSTRAINT fk_ecp_event FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    CONSTRAINT fk_ecp_category FOREIGN KEY (category_id) REFERENCES event_categories(id) ON DELETE CASCADE
);
CREATE INDEX idx_ecp_category ON event_category_pivot (category_id);

CREATE TABLE event_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL,
    session_name VARCHAR(150),                    -- e.g. "Day 1", "Matinee"
    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ NOT NULL,
    gate_open_at TIMESTAMPTZ NOT NULL,             -- default start_at - 2h
    gate_close_at TIMESTAMPTZ NOT NULL,            -- default start_at + 1h
    status VARCHAR(20) NOT NULL DEFAULT 'scheduled',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_event_sessions_event FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    CONSTRAINT chk_session_times CHECK (end_at > start_at)
);
CREATE INDEX idx_event_sessions_event ON event_sessions (event_id);
CREATE INDEX idx_event_sessions_start ON event_sessions (start_at);

CREATE TABLE artists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(150) NOT NULL,
    bio TEXT,
    photo_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE event_artists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_session_id UUID NOT NULL,
    artist_id UUID NOT NULL,
    role VARCHAR(50) DEFAULT 'lineup',            -- headliner, guest_star, lineup
    performance_order INTEGER,
    stage_time TIMESTAMPTZ,
    CONSTRAINT fk_event_artists_session FOREIGN KEY (event_session_id) REFERENCES event_sessions(id) ON DELETE CASCADE,
    CONSTRAINT fk_event_artists_artist FOREIGN KEY (artist_id) REFERENCES artists(id) ON DELETE CASCADE
);
CREATE INDEX idx_event_artists_session ON event_artists (event_session_id);

-- ----------------------------------------------------------------------------
-- 6. TICKETING INVENTORY
-- ----------------------------------------------------------------------------

CREATE TABLE ticket_tiers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL,
    venue_section_id UUID,                        -- NULL if not seat-mapped
    name VARCHAR(100) NOT NULL,
    tier_type ticket_tier_type NOT NULL,
    price NUMERIC(14,2) NOT NULL CHECK (price >= 0),
    max_per_transaction INTEGER NOT NULL DEFAULT 6,
    sale_start_at TIMESTAMPTZ,
    sale_end_at TIMESTAMPTZ,
    is_seated BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_ticket_tiers_event FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    CONSTRAINT fk_ticket_tiers_section FOREIGN KEY (venue_section_id) REFERENCES venue_sections(id) ON DELETE SET NULL
);
CREATE INDEX idx_ticket_tiers_event ON ticket_tiers (event_id);

CREATE TABLE ticket_inventories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_tier_id UUID NOT NULL,
    event_session_id UUID NOT NULL,
    total_quota INTEGER NOT NULL CHECK (total_quota >= 0),
    sold_count INTEGER NOT NULL DEFAULT 0 CHECK (sold_count >= 0),
    held_count INTEGER NOT NULL DEFAULT 0 CHECK (held_count >= 0),
    version INTEGER NOT NULL DEFAULT 0,            -- optimistic locking
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_ticket_inv_tier FOREIGN KEY (ticket_tier_id) REFERENCES ticket_tiers(id) ON DELETE CASCADE,
    CONSTRAINT fk_ticket_inv_session FOREIGN KEY (event_session_id) REFERENCES event_sessions(id) ON DELETE CASCADE,
    CONSTRAINT uq_ticket_inv UNIQUE (ticket_tier_id, event_session_id),
    CONSTRAINT chk_ticket_inv_capacity CHECK (sold_count + held_count <= total_quota)
);
CREATE INDEX idx_ticket_inv_session ON ticket_inventories (event_session_id);

-- ----------------------------------------------------------------------------
-- 7. CART / BOOKING / PAYMENT
-- ----------------------------------------------------------------------------

CREATE TABLE carts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    ticket_tier_id UUID NOT NULL,
    event_session_id UUID NOT NULL,
    seat_id UUID,                                  -- NULL for non-seated tiers
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    held_until TIMESTAMPTZ NOT NULL,                -- now() + 10-15 min
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_carts_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_carts_tier FOREIGN KEY (ticket_tier_id) REFERENCES ticket_tiers(id) ON DELETE CASCADE,
    CONSTRAINT fk_carts_session FOREIGN KEY (event_session_id) REFERENCES event_sessions(id) ON DELETE CASCADE,
    CONSTRAINT fk_carts_seat FOREIGN KEY (seat_id) REFERENCES seats(id) ON DELETE SET NULL
);
CREATE INDEX idx_carts_user ON carts (user_id);
CREATE INDEX idx_carts_held_until ON carts (held_until);      -- for expiry sweep job
CREATE UNIQUE INDEX uq_carts_seat_active ON carts (seat_id) WHERE seat_id IS NOT NULL;

CREATE TABLE bookings (
    id UUID DEFAULT gen_random_uuid(),
    booking_code VARCHAR(20) NOT NULL,
    user_id UUID NOT NULL,
    event_id UUID NOT NULL,
    status booking_status NOT NULL DEFAULT 'pending',
    subtotal_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    discount_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    platform_fee_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    promo_code_id UUID,
    hold_expires_at TIMESTAMPTZ,
    paid_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

CONSTRAINT pk_bookings PRIMARY KEY (id),
CONSTRAINT uq_bookings_code UNIQUE (booking_code),
    CONSTRAINT fk_bookings_user
        FOREIGN KEY (user_id) REFERENCES users(id),

    CONSTRAINT fk_bookings_event
        FOREIGN KEY (event_id) REFERENCES events(id)
);

CREATE INDEX idx_bookings_user ON bookings(user_id);
CREATE INDEX idx_bookings_event ON bookings(event_id);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_created_at ON bookings(created_at);
-- ... additional monthly partitions created by a scheduled migration job

CREATE TABLE booking_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL,
    ticket_tier_id UUID NOT NULL,
    event_session_id UUID NOT NULL,
    unit_price NUMERIC(14,2) NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    line_total NUMERIC(14,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_booking_items_booking FOREIGN KEY (booking_id) REFERENCES bookings(id) ON DELETE CASCADE,
    CONSTRAINT fk_booking_items_tier FOREIGN KEY (ticket_tier_id) REFERENCES ticket_tiers(id) ON DELETE RESTRICT,
    CONSTRAINT fk_booking_items_session FOREIGN KEY (event_session_id) REFERENCES event_sessions(id) ON DELETE RESTRICT
);
CREATE INDEX idx_booking_items_booking ON booking_items (booking_id);

CREATE TABLE booking_seats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_item_id UUID NOT NULL,
    seat_id UUID NOT NULL,
    CONSTRAINT fk_booking_seats_item FOREIGN KEY (booking_item_id) REFERENCES booking_items(id) ON DELETE CASCADE,
    CONSTRAINT fk_booking_seats_seat FOREIGN KEY (seat_id) REFERENCES seats(id) ON DELETE RESTRICT,
    CONSTRAINT uq_booking_seats_seat UNIQUE (seat_id)   -- a seat can only be booked once (active bookings)
);
CREATE INDEX idx_booking_seats_item ON booking_seats (booking_item_id);

CREATE TABLE payment_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,                   -- e.g. "BCA Virtual Account"
    type varchar NOT NULL,
    provider VARCHAR(50) NOT NULL,                 -- midtrans, xendit, stripe
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    icon_url TEXT
);

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL,
    payment_method_id UUID NOT NULL,
    provider_transaction_id VARCHAR(150),
    amount NUMERIC(14,2) NOT NULL,
    status payment_status NOT NULL DEFAULT 'pending',
    expires_at TIMESTAMPTZ,
    paid_at TIMESTAMPTZ,
    raw_callback_payload JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_payments_booking FOREIGN KEY (booking_id) REFERENCES bookings(id) ON DELETE CASCADE,
    CONSTRAINT fk_payments_method FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id) ON DELETE RESTRICT
);
CREATE INDEX idx_payments_booking ON payments (booking_id);
CREATE INDEX idx_payments_status ON payments (status);
CREATE INDEX idx_payments_provider_tx ON payments (provider_transaction_id);


CREATE TABLE refunds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL,
    payment_id UUID NOT NULL,
    amount NUMERIC(14,2) NOT NULL,
    reason TEXT,
    status refund_status NOT NULL DEFAULT 'requested',
    requested_by UUID NOT NULL,
    processed_by UUID,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_refunds_booking FOREIGN KEY (booking_id) REFERENCES bookings(id) ON DELETE CASCADE,
    CONSTRAINT fk_refunds_payment FOREIGN KEY (payment_id) REFERENCES payments(id) ON DELETE CASCADE,
    CONSTRAINT fk_refunds_requester FOREIGN KEY (requested_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT fk_refunds_processor FOREIGN KEY (processed_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX idx_refunds_booking ON refunds (booking_id);
CREATE INDEX idx_refunds_status ON refunds (status);

-- ----------------------------------------------------------------------------
-- 8. PROMO CODE
-- ----------------------------------------------------------------------------

CREATE TABLE promo_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL,
    event_id UUID,                                 -- NULL = platform-wide
    discount_type promo_discount_type NOT NULL,
    discount_value NUMERIC(14,2) NOT NULL,
    max_discount_amount NUMERIC(14,2),
    usage_limit INTEGER,
    usage_count INTEGER NOT NULL DEFAULT 0,
    per_user_limit INTEGER NOT NULL DEFAULT 1,
    valid_from TIMESTAMPTZ NOT NULL,
    valid_until TIMESTAMPTZ NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_promo_codes_event FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    CONSTRAINT uq_promo_codes_code UNIQUE (code)
);
CREATE INDEX idx_promo_codes_event ON promo_codes (event_id);
CREATE INDEX idx_promo_codes_active ON promo_codes (is_active, valid_from, valid_until);

ALTER TABLE bookings ADD CONSTRAINT fk_bookings_promo FOREIGN KEY (promo_code_id) REFERENCES promo_codes(id) ON DELETE SET NULL;

CREATE TABLE promo_code_usages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    promo_code_id UUID NOT NULL,
    user_id UUID NOT NULL,
    booking_id UUID NOT NULL,
    discount_applied NUMERIC(14,2) NOT NULL,
    used_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_promo_usage_promo FOREIGN KEY (promo_code_id) REFERENCES promo_codes(id) ON DELETE CASCADE,
    CONSTRAINT fk_promo_usage_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_promo_usage_booking FOREIGN KEY (booking_id) REFERENCES bookings(id) ON DELETE CASCADE
);
CREATE INDEX idx_promo_usage_user ON promo_code_usages (promo_code_id, user_id);

-- ----------------------------------------------------------------------------
-- 9. TICKET / TRANSFER / CHECK-IN
-- ----------------------------------------------------------------------------

CREATE TABLE tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_item_id UUID NOT NULL,
    owner_user_id UUID NOT NULL,                   -- current holder (may differ after transfer)
    ticket_code VARCHAR(30) NOT NULL,
    qr_token VARCHAR(255) NOT NULL,                 -- rotating dynamic token seed
    qr_token_updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status ticket_status NOT NULL DEFAULT 'issued',
    watermark_user_hash VARCHAR(64),                -- anti-screenshot watermark payload
    issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_tickets_booking_item FOREIGN KEY (booking_item_id) REFERENCES booking_items(id) ON DELETE CASCADE,
    CONSTRAINT fk_tickets_owner FOREIGN KEY (owner_user_id) REFERENCES users(id) ON DELETE RESTRICT,
    CONSTRAINT uq_tickets_code UNIQUE (ticket_code)
);
CREATE INDEX idx_tickets_owner ON tickets (owner_user_id);
CREATE INDEX idx_tickets_booking_item ON tickets (booking_item_id);
CREATE INDEX idx_tickets_status ON tickets (status);

CREATE TABLE ticket_transfers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID NOT NULL,
    from_user_id UUID NOT NULL,
    to_user_id UUID NOT NULL,
    status transfer_status NOT NULL DEFAULT 'pending',
    requested_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMPTZ,
    CONSTRAINT fk_transfers_ticket FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE,
    CONSTRAINT fk_transfers_from FOREIGN KEY (from_user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_transfers_to FOREIGN KEY (to_user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX idx_transfers_ticket ON ticket_transfers (ticket_id);

CREATE TABLE check_ins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID NOT NULL,
    event_session_id UUID NOT NULL,
    scanned_by UUID NOT NULL,                       -- gatekeeper user_id
    result checkin_result NOT NULL,
    device_id VARCHAR(100),
    scanned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_offline_sync BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT fk_checkins_ticket FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE,
    CONSTRAINT fk_checkins_session FOREIGN KEY (event_session_id) REFERENCES event_sessions(id) ON DELETE CASCADE,
    CONSTRAINT fk_checkins_scanner FOREIGN KEY (scanned_by) REFERENCES users(id) ON DELETE RESTRICT
);
CREATE INDEX idx_checkins_ticket ON check_ins (ticket_id);
CREATE INDEX idx_checkins_session ON check_ins (event_session_id, result);

-- ----------------------------------------------------------------------------
-- 10. NOTIFICATION
-- ----------------------------------------------------------------------------

CREATE TABLE notification_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(80) NOT NULL,                      -- e.g. BOOKING_CONFIRMED
    channel varchar NOT NULL,
    subject VARCHAR(200),
    body_template TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_notif_template UNIQUE (code, channel)
);

CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    template_id UUID,
    channel varchar NOT NULL,
    title VARCHAR(200),
    body TEXT,
    status varchar,
    metadata JSONB,
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_notifications_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_notifications_template FOREIGN KEY (template_id) REFERENCES notification_templates(id) ON DELETE SET NULL
);
CREATE INDEX idx_notifications_user ON notifications (user_id, status);
CREATE INDEX idx_notifications_created ON notifications (created_at);

-- ----------------------------------------------------------------------------
-- 11. REVIEW & COMMUNITY
-- ----------------------------------------------------------------------------

CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL,
    user_id UUID NOT NULL,
    booking_id UUID NOT NULL,                        -- proof of attendance
    rating SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment TEXT,
    photo_urls TEXT[],
    status review_status NOT NULL DEFAULT 'visible',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_reviews_event FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    CONSTRAINT fk_reviews_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_reviews_booking FOREIGN KEY (booking_id) REFERENCES bookings(id) ON DELETE CASCADE,
    CONSTRAINT uq_reviews_booking UNIQUE (booking_id)
);
CREATE INDEX idx_reviews_event ON reviews (event_id, status);

CREATE TABLE favorites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    event_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_favorites_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_favorites_event FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    CONSTRAINT uq_favorites UNIQUE (user_id, event_id)
);
CREATE INDEX idx_favorites_event ON favorites (event_id);

-- ----------------------------------------------------------------------------
-- 12. PLATFORM FINANCE / ADMIN
-- ----------------------------------------------------------------------------

CREATE TABLE commission_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scope VARCHAR(20) NOT NULL DEFAULT 'global',      -- global | organizer | event
    scope_ref_id UUID,                                 -- organizer_id or event_id when scoped
    rate_percentage NUMERIC(5,2) NOT NULL CHECK (rate_percentage BETWEEN 0 AND 100),
    effective_from TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_commission_scope ON commission_settings (scope, scope_ref_id);

CREATE TABLE payouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organizer_id UUID NOT NULL,
    event_id UUID NOT NULL,
    gross_amount NUMERIC(14,2) NOT NULL,
    commission_amount NUMERIC(14,2) NOT NULL,
    net_amount NUMERIC(14,2) NOT NULL,
    status payout_status NOT NULL DEFAULT 'scheduled',
    scheduled_at TIMESTAMPTZ NOT NULL,                  -- event end + 3 days
    paid_at TIMESTAMPTZ,
    bank_reference VARCHAR(150),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_payouts_organizer FOREIGN KEY (organizer_id) REFERENCES organizers(id) ON DELETE RESTRICT,
    CONSTRAINT fk_payouts_event FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE RESTRICT
);
CREATE INDEX idx_payouts_organizer ON payouts (organizer_id);
CREATE INDEX idx_payouts_status ON payouts (status);

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_user_id UUID,
    action audit_action NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID,
    before_state JSONB,
    after_state JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_audit_logs_actor FOREIGN KEY (actor_user_id) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX idx_audit_logs_entity ON audit_logs (entity_type, entity_id);
CREATE INDEX idx_audit_logs_actor ON audit_logs (actor_user_id);


-- ----------------------------------------------------------------------------
-- 13. TRIGGERS
-- ----------------------------------------------------------------------------

-- 13.1 Generic updated_at trigger
CREATE OR REPLACE FUNCTION trg_set_updated_at() RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at_users BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION trg_set_updated_at();
CREATE TRIGGER set_updated_at_events BEFORE UPDATE ON events FOR EACH ROW EXECUTE FUNCTION trg_set_updated_at();
CREATE TRIGGER set_updated_at_bookings BEFORE UPDATE ON bookings FOR EACH ROW EXECUTE FUNCTION trg_set_updated_at();
CREATE TRIGGER set_updated_at_payments BEFORE UPDATE ON payments FOR EACH ROW EXECUTE FUNCTION trg_set_updated_at();
CREATE TRIGGER set_updated_at_ticket_tiers BEFORE UPDATE ON ticket_tiers FOR EACH ROW EXECUTE FUNCTION trg_set_updated_at();

-- 13.2 Deduct inventory when payment succeeds (sold_count += qty, held_count -= qty)
CREATE OR REPLACE FUNCTION trg_deduct_inventory_on_payment_success() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status = 'success' AND OLD.status <> 'success' THEN
        UPDATE ticket_inventories ti
        SET sold_count = sold_count + bi.quantity,
            held_count = GREATEST(held_count - bi.quantity, 0),
            version = version + 1,
            updated_at = NOW()
        FROM booking_items bi
        WHERE bi.booking_id = NEW.booking_id
          AND ti.ticket_tier_id = bi.ticket_tier_id
          AND ti.event_session_id = bi.event_session_id;

        UPDATE bookings SET status = 'paid', paid_at = NOW() WHERE id = NEW.booking_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER deduct_inventory_after_payment
    AFTER UPDATE ON payments
    FOR EACH ROW EXECUTE FUNCTION trg_deduct_inventory_on_payment_success();

-- 13.3 Full-text search vector maintenance for events
CREATE OR REPLACE FUNCTION trg_events_search_vector() RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('simple', coalesce(NEW.title, '')), 'A') ||
        setweight(to_tsvector('simple', coalesce(NEW.description, '')), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER events_search_vector_update
    BEFORE INSERT OR UPDATE OF title, description ON events
    FOR EACH ROW EXECUTE FUNCTION trg_events_search_vector();

-- ----------------------------------------------------------------------------
-- 14. SEED DATA (minimal)
-- ----------------------------------------------------------------------------
INSERT INTO roles (name, description) VALUES
    ('super_admin', 'Full platform access'),
    ('admin_event', 'Event moderation and approval'),
    ('organizer', 'Can create and manage own events'),
    ('customer', 'Can browse and purchase tickets'),
    ('gatekeeper', 'Can scan tickets at venue');

INSERT INTO commission_settings (scope, rate_percentage) VALUES ('global', 7.5);