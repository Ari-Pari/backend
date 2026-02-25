-- name: SearchDances :many
SELECT
    d.id,
    d.translation_id,
    COALESCE(
            CASE sqlc.arg('lang')::text
                WHEN 'en' THEN t.eng_name
                WHEN 'ru' THEN t.ru_name
                WHEN 'hy' THEN t.arm_name
                ELSE t.eng_name
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
    -- Отдельный массив для ID регионов
    COALESCE(
                    array_agg(DISTINCT r.id) FILTER (WHERE r.id IS NOT NULL),
                    ARRAY[]::bigint[]
    ) AS region_ids,
    -- Отдельный массив для названий регионов
    COALESCE(
                    array_agg(
                    DISTINCT COALESCE(
                            CASE sqlc.arg('lang')::text
                                WHEN 'en' THEN rt.eng_name
                                WHEN 'ru' THEN rt.ru_name
                                WHEN 'hy' THEN rt.arm_name
                                ELSE rt.eng_name
                                END,
                            r.name
                             )
                             ) FILTER (WHERE r.id IS NOT NULL),
                    ARRAY[]::text[]
    ) AS region_names
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
            THEN CONCAT_WS(' ', t.eng_name, t.ru_name, t.arm_name, d.name) ILIKE ('%' || sqlc.arg(search_text)::text || '%')
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
                WHEN 'hy' THEN t.arm_name
                ELSE t.eng_name
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
    CASE WHEN sqlc.arg(order_by_popularity)::boolean = true AND sqlc.arg(reverse_order)::boolean = false THEN d.popularity END ASC,
    CASE WHEN sqlc.arg(order_by_popularity)::boolean = true AND sqlc.arg(reverse_order)::boolean = true THEN d.popularity END DESC,

    CASE WHEN sqlc.arg(order_by_name)::boolean = true AND sqlc.arg(reverse_order)::boolean = false THEN COALESCE(
            CASE sqlc.arg('lang')::text
                WHEN 'en' THEN t.eng_name
                WHEN 'ru' THEN t.ru_name
                WHEN 'hy' THEN t.arm_name
                ELSE t.eng_name
                END,
            d.name
                                                                                                        ) END ASC,
    CASE WHEN sqlc.arg(order_by_name)::boolean = true AND sqlc.arg(reverse_order)::boolean = true THEN COALESCE(
            CASE sqlc.arg('lang')::text
                WHEN 'en' THEN t.eng_name
                WHEN 'ru' THEN t.ru_name
                WHEN 'hy' THEN t.arm_name
                ELSE t.eng_name
                END,
            d.name
                                                                                                       ) END DESC,

    CASE WHEN sqlc.arg(order_by_created_at)::boolean = true AND sqlc.arg(reverse_order)::boolean = false THEN d.created_at END ASC,
    CASE WHEN sqlc.arg(order_by_created_at)::boolean = true AND sqlc.arg(reverse_order)::boolean = true THEN d.created_at END DESC,
    d.created_at DESC
LIMIT  sqlc.arg('limit')::int
    OFFSET sqlc.arg('offset')::int;