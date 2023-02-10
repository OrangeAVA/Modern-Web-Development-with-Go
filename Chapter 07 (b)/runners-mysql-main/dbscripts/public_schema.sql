-- runners
CREATE TABLE runners (
    id int NOT NULL AUTO_INCREMENT,
    first_name varchar(100) NOT NULL,
    last_name varchar(100) NOT NULL,
    age integer,
    is_active boolean DEFAULT TRUE,
    country varchar(60) NOT NULL,
    personal_best time,
    season_best time,
    CONSTRAINT runners_pk PRIMARY KEY (id)
)
ENGINE = InnoDB;

CREATE INDEX runners_country
ON runners (country);

CREATE INDEX runners_season_best
ON runners (season_best);

-- results
CREATE TABLE results (
    id int NOT NULL AUTO_INCREMENT,
    runner_id int NOT NULL,
    race_result time NOT NULL,
    location varchar(100) NOT NULL,
    position integer,
    result_year integer NOT NULL,
    CONSTRAINT results_pk PRIMARY KEY (id),
    CONSTRAINT fk_results_runner_id FOREIGN KEY (runner_id)
        REFERENCES runners (id)
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)
ENGINE = InnoDB;
