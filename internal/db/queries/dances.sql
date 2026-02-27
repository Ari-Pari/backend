-- name: GetDanceByID :one
SELECT d.id, d.complexity, d.photo_key, d.gender, d.paces, d.genres, d.handshakes,
COALESCE(CASE WHEN sqlc.narg('lang')::text = 'ru' THEN t.ru_name WHEN sqlc.narg('lang')::text = 'en' THEN t.eng_name WHEN sqlc.narg('lang')::text = 'arm' THEN t.arm_name ELSE d.name END, d.name)::text AS name
FROM dances d
LEFT JOIN translations t ON d.translation_id = t.id
WHERE d.id = $1;

-- name: GetRegionsByDanceID :many
SELECT 
    r.id,
    COALESCE(
        CASE 
            WHEN sqlc.narg('lang')::text = 'ru' THEN t.ru_name
            WHEN sqlc.narg('lang')::text = 'en' THEN t.eng_name
            WHEN sqlc.narg('lang')::text = 'arm' THEN t.arm_name
            ELSE r.name
        END, 
        r.name
    )::text AS name
FROM regions r
LEFT JOIN translations t ON r.translation_id = t.id
JOIN dance_region dr ON dr.region_id = r.id
WHERE dr.dance_id = $1;

-- name: GetVideosByDanceID :many
SELECT v.id, v.name, v.link, v.type 
FROM videos v 
JOIN dance_videos dv ON dv.video_id = v.id 
WHERE dv.dance_id = $1;

-- name: GetSongsByDanceID :many
SELECT s.id, s.name, s.file_key 
FROM songs s 
JOIN dance_song ds ON ds.song_id = s.id 
WHERE ds.dance_id = $1;