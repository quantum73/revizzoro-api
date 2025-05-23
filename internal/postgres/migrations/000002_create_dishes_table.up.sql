CREATE TABLE IF NOT EXISTS dishes
(
    id         BIGSERIAL    PRIMARY KEY,
    name       VARCHAR(256) NOT NULL,
    price      BIGINT       NOT NULL,
    score      INT          NOT NULL,
    restaurant_id BIGINT       REFERENCES restaurants (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY(restaurant_id) REFERENCES restaurants (id),
    CHECK ( price > 0),
    CHECK ( score BETWEEN 1 and 5)
);