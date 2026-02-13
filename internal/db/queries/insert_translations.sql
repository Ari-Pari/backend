-- name: CreateTranslations :copyfrom
INSERT INTO translations (eng_name,
                          ru_name,
                          arm_name)
VALUES ($1,
        $2,
        $3);