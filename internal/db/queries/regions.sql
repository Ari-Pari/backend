-- name: ListRegions :many
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
ORDER BY r.id;