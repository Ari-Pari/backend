CREATE TABLE translations (
    id BIGSERIAL PRIMARY KEY,
    eng_name VARCHAR NOT NULL DEFAULT '',
    ru_name VARCHAR NOT NULL DEFAULT '',
    arm_name VARCHAR NOT NULL DEFAULT ''
);

CREATE TABLE regions (
    id BIGSERIAL PRIMARY KEY,
    translation_id BIGINT REFERENCES translations(id),
    name VARCHAR NOT NULL
);

CREATE TABLE ensemble (
    id BIGSERIAL PRIMARY KEY,
    translation_id BIGINT REFERENCES translations(id),
    name VARCHAR NOT NULL,
    link VARCHAR NOT NULL DEFAULT ''
);

CREATE TABLE songs (
    id BIGSERIAL PRIMARY KEY,
    translation_id BIGINT REFERENCES translations(id),
    file_key VARCHAR NOT NULL,
    name VARCHAR NOT NULL
);

CREATE TABLE dances (
    id BIGSERIAL PRIMARY KEY,
    translation_id BIGINT REFERENCES translations(id),
    name VARCHAR NOT NULL,
    complexity INTEGER NOT NULL CHECK (complexity >= 1 AND complexity <= 5),
    photo_link VARCHAR,
    gender VARCHAR NOT NULL, 
    paces INTEGER[],
    popularity INTEGER NOT NULL DEFAULT 0,
    created_by BIGINT,
    genres TEXT[],
    handshakes TEXT[],
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Индексы GIN, так как создаем их на массивы
CREATE INDEX idx_dances_genres ON dances USING GIN (genres);
CREATE INDEX idx_dances_handshakes ON dances USING GIN (handshakes);
CREATE INDEX idx_dances_paces ON dances USING GIN (paces);

CREATE TABLE videos (
    id BIGSERIAL PRIMARY KEY,
    link VARCHAR NOT NULL,
    translation_id BIGINT REFERENCES translations(id),
    name VARCHAR NOT NULL,
    dance_id BIGINT NOT NULL REFERENCES dances(id) ON DELETE CASCADE,
    type VARCHAR NOT NULL
);

-- Связи Many-to-Many

CREATE TABLE dance_region (
    id BIGSERIAL PRIMARY KEY,
    dance_id BIGINT NOT NULL REFERENCES dances(id) ON DELETE CASCADE,
    region_id BIGINT NOT NULL REFERENCES regions(id) ON DELETE CASCADE,
    CONSTRAINT unique_dance_region UNIQUE (dance_id, region_id) -- Защита от дублей
);

CREATE TABLE dance_song (
    id BIGSERIAL PRIMARY KEY,
    dance_id BIGINT NOT NULL REFERENCES dances(id) ON DELETE CASCADE,
    song_id BIGINT NOT NULL REFERENCES songs(id) ON DELETE CASCADE,
    CONSTRAINT unique_dance_song UNIQUE (dance_id, song_id)
);

CREATE TABLE song_ensemble (
    id BIGSERIAL PRIMARY KEY,
    song_id BIGINT NOT NULL REFERENCES songs(id) ON DELETE CASCADE,
    ensemble_id BIGINT NOT NULL REFERENCES ensemble(id) ON DELETE CASCADE,
    CONSTRAINT unique_song_ensemble UNIQUE (song_id, ensemble_id)
);