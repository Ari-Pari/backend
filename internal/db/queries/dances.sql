-- name: GetDanceByID :one
SELECT 
    d.id,
    d.complexity,
    d.photo_key,
    d.gender,
    d.paces,
    d.genres,
    d.handshakes,
    COALESCE(
        CASE 
            WHEN sqlc.narg('lang')::text = 'ru' THEN t.ru_name
            WHEN sqlc.narg('lang')::text = 'en' THEN t.eng_name
            WHEN sqlc.narg('lang')::text = 'arm' THEN t.arm_name
            ELSE d.name
        END, 
        d.name
    )::text AS name,
    (SELECT COALESCE(json_agg(json_build_object('id', reg.id, 'name', reg.name)), '[]'::json)
     FROM regions reg
     JOIN dance_region dr ON dr.region_id = reg.id
     WHERE dr.dance_id = d.id) AS regions_json,
    (SELECT COALESCE(json_agg(json_build_object('id', v.id, 'name', v.name, 'link', v.link, 'type', v.type)), '[]'::json)
     FROM videos v
     JOIN dance_videos dv ON dv.video_id = v.id
     WHERE dv.dance_id = d.id) AS videos_json
FROM dances d
LEFT JOIN translations t ON d.translation_id = t.id
WHERE d.id = $1;