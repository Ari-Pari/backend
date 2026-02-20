-- name: TruncateAllTables :exec
TRUNCATE TABLE
    dance_region,
    dance_song,
    song_artist,
    dance_videos,
    regions,
    artists,
    songs,
    dances,
    videos,
    translations
    RESTART IDENTITY CASCADE;

-- name: GetTranslations :many
SELECT id, eng_name, ru_name, arm_name
FROM translations;

-- name: InsertTranslations :many
INSERT INTO translations (eng_name, ru_name, arm_name)
SELECT unnest(@eng_names::text[]) as eng_name,
       unnest(@ru_names::text[])  as ru_name,
       unnest(@arm_names::text[]) as arm_name
    RETURNING id;

-- name: GetRegions :many
SELECT id, translation_id, name
FROM regions;

-- name: InsertRegions :exec
INSERT INTO regions (id, translation_id, name)
SELECT unnest(@ids::bigint[])             as id,
       unnest(@translation_ids::bigint[]) as translation_id,
       unnest(@names::text[])             as name;

-- name: GetArtists :many
SELECT id, translation_id, name, link, deleted_at
FROM artists;

-- name: GetDances :many
SELECT id,
       translation_id,
       name,
       complexity,
       gender,
       paces,
       popularity,
       genres,
       handshakes,
       deleted_at
FROM dances;

-- name: InsertDance :exec
INSERT INTO dances (id, translation_id, name, complexity, gender,
                    paces, popularity, genres, handshakes, deleted_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: GetDanceRegions :many
SELECT dance_id, region_id
FROM dance_region;

-- name: InsertDanceRegions :exec
INSERT INTO dance_region (dance_id, region_id)
SELECT unnest(@dance_ids::bigint[])  as dance_id,
       unnest(@region_ids::bigint[]) as region_id ON CONFLICT (dance_id, region_id) DO NOTHING;

-- name: GetSongs :many
SELECT id, translation_id, name, file_key
FROM songs;

-- name: InsertSongs :exec
INSERT INTO songs (id, translation_id, name, file_key)
SELECT unnest(@ids::bigint[])             as id,
       unnest(@translation_ids::bigint[]) as translation_id,
       unnest(@names::text[])             as name,
       unnest(@file_keys::text[])         as file_key;

-- name: GetDanceSongs :many
SELECT dance_id, song_id
FROM dance_song;

-- name: InsertDanceSongs :exec
INSERT INTO dance_song (dance_id, song_id)
SELECT unnest(@dance_ids::bigint[]) as dance_id,
       unnest(@song_ids::bigint[])  as song_id ON CONFLICT (dance_id, song_id) DO NOTHING;

-- name: GetVideos :many
SELECT id, link, translation_id, name, type
FROM videos;

-- name: InsertVideos :many
INSERT INTO videos (link, translation_id, name, type)
SELECT unnest(@links::text[])             as link,
       unnest(@translation_ids::bigint[]) as translation_id,
       unnest(@names::text[])             as name,
       unnest(@types::text[])             as type
    RETURNING id;

-- name: GetDanceVideos :many
SELECT dance_id, video_id
FROM dance_videos;

-- name: InsertDanceVideos :exec
INSERT INTO dance_videos (dance_id, video_id)
SELECT unnest(@dance_ids::bigint[]) as dance_id,
       unnest(@video_ids::bigint[]) as video_id ON CONFLICT (dance_id, video_id) DO NOTHING;

-- name: InsertArtists :exec
INSERT INTO artists (id, translation_id, name, link, deleted_at)
SELECT unnest(@ids::bigint[])              as id,
       unnest(@translation_ids::bigint[])  as translation_id,
       unnest(@names::text[])              as name,
       unnest(@links::text[])              as link,
       unnest(@deleted_ats::timestamptz[]) as deleted_at;


-- name: InsertSongArtists :exec
INSERT INTO song_artist (song_id, artist_id)
SELECT unnest(@song_ids::bigint[])   as song_id,
       unnest(@artist_ids::bigint[]) as artist_id ON CONFLICT (song_id, artist_id) DO NOTHING;