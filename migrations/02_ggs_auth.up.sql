ALTER TABLE sellers
    ADD COLUMN email text;

ALTER TABLE sellers
    ALTER COLUMN email
        SET DEFAULT ' ';