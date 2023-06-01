CREATE TABLE IF NOT EXISTS ingredient (
   slug VARCHAR UNIQUE NOT NULL PRIMARY KEY,
   name VARCHAR NOT NULL,
   price100 INTEGER NOT NULL
);

INSERT INTO ingredient
VALUES ('tomato', 'Tomato', 50),
       ('sliced_mushroom', 'Sliced Mushrooms', 50),
       ('feta_cheese', 'Feta Cheeze', 100),
       ('sausage', 'Sausage', 100),
       ('sliced_onion', 'Sliced Onion', 50),
       ('mozzarella_cheese', 'Mozzarella Cheeze', 50),
       ('oregano', 'Oregano', 100),
       ('bacon', 'Bacon', 100)
;

CREATE TABLE IF NOT EXISTS pizza (
    slug VARCHAR UNIQUE NOT NULL PRIMARY KEY,
    name VARCHAR NOT NULL
);

INSERT INTO pizza
VALUES ('fun', 'Fun'),
       ('super_mushroom', 'Super Mushroom');

CREATE TABLE IF NOT EXISTS recipe (
    pizza VARCHAR NOT NULL REFERENCES pizza(slug) ON DELETE CASCADE,
    ingredient VARCHAR NOT NULL REFERENCES ingredient(slug) ON DELETE CASCADE,
    position INTEGER,
    PRIMARY KEY (pizza, ingredient)
);

INSERT INTO recipe
VALUES ('fun', 'tomato', 1),
       ('fun', 'sliced_mushroom', 2),
       ('fun', 'feta_cheese', 3),
       ('fun', 'sausage', 4),
       ('fun', 'sliced_onion', 5),
       ('fun', 'mozzarella_cheese', 6),
       ('fun', 'oregano', 7),
       ('super_mushroom', 'tomato', 1),
       ('super_mushroom', 'bacon', 2),
       ('super_mushroom', 'mozzarella_cheese', 3),
       ('super_mushroom', 'sliced_mushroom', 4),
       ('super_mushroom', 'oregano', 5)
;