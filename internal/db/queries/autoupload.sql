-- name: TruncateAllTables :exec
TRUNCATE TABLE
    dance_region,
    dance_song,
    song_ensemble,
    dance_videos,
    regions,
    ensembles,
    songs,
    dances,
    videos,
    translations
    RESTART IDENTITY CASCADE;

-- name: GetTranslations :many
SELECT id, eng_name, ru_name, arm_name FROM translations;

-- name: InsertTranslations :many
INSERT INTO translations (eng_name, ru_name, arm_name)
SELECT unnest(@eng_names::text[]) as eng_name,
       unnest(@ru_names::text[])  as ru_name,
       unnest(@arm_names::text[]) as arm_name
    RETURNING id;

-- name: GetRegions :many
SELECT id, translation_id, name FROM regions;

-- name: InsertRegions :exec
INSERT INTO regions (id, translation_id, name)
SELECT unnest(@ids::bigint[])             as id,
       unnest(@translation_ids::bigint[]) as translation_id,
       unnest(@names::text[])             as name;

-- name: GetEnsembles :many
SELECT id, translation_id, name, link FROM ensembles;

-- name: InsertEnsembles :many
INSERT INTO ensembles (translation_id, name, link)
SELECT unnest(@translation_ids::bigint[]) as translation_id,
       unnest(@names::text[])             as name,
       unnest(@links::text[])             as link
    RETURNING id;

-- name: GetDances :many
SELECT id, translation_id, name, complexity, gender,
       paces, popularity, genres, handshakes, deleted_at
FROM dances;

-- name: InsertDance :exec
INSERT INTO dances (id, translation_id, name, complexity, gender,
                    paces, popularity, genres, handshakes, deleted_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: GetDanceRegions :many
SELECT dance_id, region_id FROM dance_region;

-- name: InsertDanceRegions :exec
INSERT INTO dance_region (dance_id, region_id)
SELECT unnest(@dance_ids::bigint[])  as dance_id,
       unnest(@region_ids::bigint[]) as region_id ON CONFLICT (dance_id, region_id) DO NOTHING;

-- name: GetSongs :many
SELECT id, translation_id, name, file_key FROM songs;

-- name: InsertSongs :exec
INSERT INTO songs (id, translation_id, name, file_key)
SELECT unnest(@ids::bigint[])             as id,
       unnest(@translation_ids::bigint[]) as translation_id,
       unnest(@names::text[])             as name,
       unnest(@file_keys::text[])         as file_key;

-- name: GetDanceSongs :many
SELECT dance_id, song_id FROM dance_song;

-- name: InsertDanceSongs :exec
INSERT INTO dance_song (dance_id, song_id)
SELECT unnest(@dance_ids::bigint[]) as dance_id,
       unnest(@song_ids::bigint[])  as song_id ON CONFLICT (dance_id, song_id) DO NOTHING;

-- name: GetVideos :many
SELECT id, link, translation_id, name, type FROM videos;

-- name: InsertVideos :many
INSERT INTO videos (link, translation_id, name, type)
SELECT unnest(@links::text[])             as link,
       unnest(@translation_ids::bigint[]) as translation_id,
       unnest(@names::text[])             as name,
       unnest(@types::text[])             as type
    RETURNING id;

-- name: GetDanceVideos :many
SELECT dance_id, video_id FROM dance_videos;

-- name: InsertDanceVideos :exec
INSERT INTO dance_videos (dance_id, video_id)
SELECT unnest(@dance_ids::bigint[]) as dance_id,
       unnest(@video_ids::bigint[]) as video_id ON CONFLICT (dance_id, video_id) DO NOTHING;