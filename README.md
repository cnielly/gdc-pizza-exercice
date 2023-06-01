# The pizza exercise in Go

## Pre-requisites
- Go (I'm on version `1.20.4`)
- [https://github.com/codegangsta/gin]() to get the auto recompile

## Install the project and run it
- Clone the repository
- Run `docker compose up` to launch the postgres container and initialize the DB
- Go into the `src` folder and run `gin -i run main.go` to launch the app

## Play with the API

_Note that I used the [Gin Gonic web framework](https://github.com/gin-gonic/gin)_

### Ingredients

```json
// GET http://localhost:3231/ingredients
// returns 200
[
  {
    "slug": "tomato",
    "name": "Tomato",
    "price100": 50
  },
  {
    "slug": "sliced_mushroom",
    "name": "Sliced Mushrooms",
    "price100": 50
  },
  // ...
]

// GET http://localhost:3231/ingredients/tomato
// returns 200
{
    "slug": "tomato",
    "name": "Tomato",
    "price100": 50
}

// POST http://localhost:3231/ingredients
// with
{
  "slug": "ham",
  "name": "Ham",
  "price100": 150
}
// returns 201
{
  "slug": "ham",
  "name": "Ham",
  "price100": 150
}

// DELETE http://localhost:3231/ingredients/ham
// returns 204

// PUT http://localhost:3231/ingredients/ham
// with
{
  "name": "Ham new name",
  "price100": 250
}
// returns 200
{
  "name": "Ham new name",
  "price100": 250
}
```

### Pizzas

```json
// GET http://localhost:3231/pizzas
// returns 200
[
    {
        "slug": "super_mushroom",
        "name": "Super Mushroom",
        "price100": 350
    },
    {
        "slug": "fun",
        "name": "Fun",
        "price100": 500
    }
]

// GET http://localhost:3231/pizzas/fun
// returns 200
{
  "slug": "fun",
  "name": "Fun",
  "price100": 500
}

// POST http://localhost:3231/pizzas
// with 
{
  "slug": "hawaiian",
  "name": "Hawaiian"
}
// returns 201
{
  "slug": "hawaiian",
  "name": "Hawaiian"
  "price100": 500
}

// DELETE http://localhost:3231/pizzas/hawaiian
// returns 204
```

### Recipes
```json
// GET http://localhost:3231/recipes/fun
// returns 200
{
  "ingredients": [
    "tomato",
    "sliced_mushroom",
    "feta_cheese",
    "sausage",
    "sliced_onion",
    "mozzarella_cheese",
    "oregano"
  ]
}

// POST http://localhost:3231/recipes
// with
{
  "pizza_slug": "hawaiian",
  "ingredients": [
    "pineapple",
    "mozzarella_cheese",
    "ham"
  ]
}
// returns 201
{
  "ingredients": [
    "pineapple",
    "mozzarella_cheese",
    "ham"
  ]
}
```
