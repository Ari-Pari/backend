CREATE TABLE translations (
    id BIGSERIAL PRIMARY KEY,
    eng_name VARCHAR,
    ru_name VARCHAR,
    arm_name VARCHAR,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE TABLE regions (
    id BIGSERIAL PRIMARY KEY,
    translation_id BIGINT,
    name VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE TABLE ensemble (
    id BIGSERIAL PRIMARY KEY,
    translation_id BIGINT,
    name VARCHAR NOT NULL,
    link VARCHAR NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE TABLE songs (
    id BIGSERIAL PRIMARY KEY,
    translation_id BIGINT,
    file_key VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE TABLE dances (
    id BIGSERIAL PRIMARY KEY,
    translation_id BIGINT,
    name VARCHAR NOT NULL,
    complexity INTEGER NOT NULL CHECK (complexity >= 1 AND complexity <= 5),
    photo_key VARCHAR,
    gender VARCHAR NOT NULL, 
    paces INTEGER[],
    popularity INTEGER NOT NULL DEFAULT 0,
    created_by BIGINT,
    genres TEXT[],
    handshakes TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE TABLE videos (
    id BIGSERIAL PRIMARY KEY,
    link VARCHAR NOT NULL,
    translation_id BIGINT,
    name VARCHAR NOT NULL,
    dance_id BIGINT,
    type VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

-- Связи Many-to-Many

CREATE TABLE dance_region (
    id BIGSERIAL PRIMARY KEY,
    dance_id BIGINT NOT NULL,
    region_id BIGINT NOT NULL,
    CONSTRAINT unique_dance_region UNIQUE (dance_id, region_id) -- Защита от дублей
);

CREATE TABLE dance_song (
    id BIGSERIAL PRIMARY KEY,
    dance_id BIGINT NOT NULL,
    song_id BIGINT NOT NULL,
    CONSTRAINT unique_dance_song UNIQUE (dance_id, song_id)
);

CREATE TABLE song_ensemble (
    id BIGSERIAL PRIMARY KEY,
    song_id BIGINT NOT NULL,
    ensemble_id BIGINT NOT NULL,
    CONSTRAINT unique_song_ensemble UNIQUE (song_id, ensemble_id)
);