-- name: SearchDances :many
SELECT
    d.id,
    d.translation_id,
    COALESCE(
            CASE sqlc.arg('lang')::text
                WHEN 'en' THEN t.eng_name
                WHEN 'ru' THEN t.ru_name
                WHEN 'am' THEN t.arm_name
                ELSE t.ru_name
                END,
            d.name
    ) AS name,
    d.complexity,
    d.photo_key         AS photo_link,
    d.gender,
    d.paces,
    d.genres,
    d.handshakes,
    d.popularity,
    d.created_at,
    d.updated_at,
    COALESCE(
                    json_agg(
                    DISTINCT jsonb_build_object(
                            'id', r.id,
                            'translation_id', r.translation_id,
                            'name', COALESCE(
                                    CASE sqlc.arg('lang')::text
                                        WHEN 'en' THEN rt.eng_name
                                        WHEN 'ru' THEN rt.ru_name
                                        WHEN 'am' THEN rt.arm_name
                                        ELSE rt.ru_name
                                        END,
                                    r.name
                                    )
                             )
                            ) FILTER (WHERE r.id IS NOT NULL),
                    '[]'::json
    ) AS regions
FROM dances d
         LEFT JOIN translations t ON t.id = d.translation_id
         LEFT JOIN dance_region dr ON dr.dance_id = d.id
         LEFT JOIN regions r       ON r.id = dr.region_id
         LEFT JOIN translations rt ON rt.id = r.translation_id
WHERE d.deleted_at IS NULL
  AND (
    CASE
        WHEN sqlc.arg(search_text)::text IS NOT NULL
            AND sqlc.arg(search_text)::text <> ''
            THEN COALESCE(
                         CASE sqlc.arg('lang')::text
                             WHEN 'en' THEN t.eng_name
                             WHEN 'ru' THEN t.ru_name
                             WHEN 'am' THEN t.arm_name
                             ELSE t.ru_name
                             END,
                         d.name
                 ) ILIKE ('%' || sqlc.arg(search_text)::text || '%')
        ELSE TRUE
        END
    )
  AND (
    CASE
        WHEN array_length(sqlc.arg(genres_in)::text[], 1) > 0
            THEN d.genres && sqlc.arg(genres_in)::text[]
        ELSE TRUE
        END
    )
  AND (
    CASE
        WHEN array_length(sqlc.arg(region_ids_in)::bigint[], 1) > 0
            THEN dr.region_id = ANY(sqlc.arg(region_ids_in)::bigint[])
        ELSE TRUE
        END
    )
  AND (
    CASE
        WHEN array_length(sqlc.arg(complexities_in)::int[], 1) > 0
            THEN d.complexity = ANY(sqlc.arg(complexities_in)::int[])
        ELSE TRUE
        END
    )
  AND (
    CASE
        WHEN array_length(sqlc.arg(genders_in)::text[], 1) > 0
            THEN d.gender = ANY(sqlc.arg(genders_in)::text[])
        ELSE TRUE
        END
    )
  AND (
    CASE
        WHEN array_length(sqlc.arg(paces_in)::int[], 1) > 0
            THEN d.paces && sqlc.arg(paces_in)::int[]
        ELSE TRUE
        END
    )
  AND (
    CASE
        WHEN array_length(sqlc.arg(handshakes_in)::text[], 1) > 0
            THEN d.handshakes && sqlc.arg(handshakes_in)::text[]
        ELSE TRUE
        END
    )
GROUP BY
    d.id,
    d.translation_id,
    COALESCE(
            CASE sqlc.arg('lang')::text
                WHEN 'en' THEN t.eng_name
                WHEN 'ru' THEN t.ru_name
                WHEN 'am' THEN t.arm_name
                ELSE t.ru_name
                END,
            d.name
    ),
    d.complexity,
    d.photo_key,
    d.gender,
    d.paces,
    d.genres,
    d.handshakes,
    d.popularity,
    d.created_at,
    d.updated_at
ORDER BY
    CASE sqlc.arg('order_by')::text
        WHEN 'popularity' THEN d.popularity::numeric
        WHEN 'name'       THEN 0::numeric  -- сортируем по имени отдельным уровнем
        WHEN 'created_at' THEN EXTRACT(EPOCH FROM d.created_at)::numeric
        ELSE EXTRACT(EPOCH FROM d.created_at)::numeric
        END
        *
    CASE
        WHEN sqlc.arg('order_dir')::text = 'DESC' THEN -1::numeric
        ELSE 1::numeric
        END,
    -- для случая order_by = 'name' включаем второе поле сортировки
    CASE
        WHEN sqlc.arg('order_by')::text = 'name' THEN
            COALESCE(
                    CASE sqlc.arg('lang')::text
                        WHEN 'en' THEN t.eng_name
                        WHEN 'ru' THEN t.ru_name
                        WHEN 'am' THEN t.arm_name
                        ELSE t.ru_name
                        END,
                    d.name
            )
        ELSE NULL
        END
LIMIT  sqlc.arg('limit')::int
    OFFSET sqlc.arg('offset')::int;
