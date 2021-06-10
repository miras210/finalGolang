ALTER TABLE comics ADD CONSTRAINT comics_year_check CHECK (year BETWEEN 1888 AND date_part('year', now()));