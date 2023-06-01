package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Ingredient struct {
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	Price100 int    `json:"price100"`
}

type Pizza struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type PizzaWithPrice100 struct {
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	Price100 int    `json:"price100"`
}

type Recipe struct {
	Ingredients []string `json:"ingredients"`
}

type DB struct {
	*sql.DB
}

func OpenDBConnection(dbSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dbSourceName)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) CloseDBConnection() error {
	return db.DB.Close()
}

func (db *DB) GetIngredients() ([]*Ingredient, error) {
	rows, err := db.Query("SELECT * FROM ingredient")
	if err != nil {
		return nil, err
	}
	var ingredients []*Ingredient
	for rows.Next() {
		ingredient := Ingredient{}
		if err := rows.Scan(&ingredient.Slug, &ingredient.Name, &ingredient.Price100); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, &ingredient)
	}

	return ingredients, nil
}

func (db *DB) GetIngredient(slug string) (*Ingredient, error) {
	row := db.QueryRow("SELECT * FROM ingredient WHERE slug = $1", slug)
	ingredient := Ingredient{}
	err := row.Scan(&ingredient.Slug, &ingredient.Name, &ingredient.Price100)
	if err != nil {
		return nil, err
	}

	return &ingredient, nil
}

func (db *DB) CreateIngredient(ingredient *Ingredient) error {
	_, err := db.Exec("INSERT INTO ingredient (slug, name, price100) VALUES ($1, $2, $3)", ingredient.Slug, ingredient.Name, ingredient.Price100)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) RemoveIngredient(slug string) error {
	_, err := db.Exec("DELETE FROM ingredient WHERE slug = $1", slug)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) UpdateIngredient(ingredient *Ingredient) (*Ingredient, error) {
	_, err := db.GetIngredient(ingredient.Slug)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("UPDATE ingredient SET name = $1, price100 = $2 WHERE slug = $3", ingredient.Name, ingredient.Price100, ingredient.Slug)
	if err != nil {
		return nil, err
	}

	return ingredient, nil
}

func (db *DB) GetPizzas() ([]*PizzaWithPrice100, error) {
	rows, err := db.Query(`SELECT p.slug, p.name, SUM(i.price100) 
			FROM pizza p
			JOIN recipe r on p.slug = r.pizza
			JOIN ingredient i on i.slug = r.ingredient
			GROUP BY p.slug, p.name`)
	if err != nil {
		return nil, err
	}
	var pizzas []*PizzaWithPrice100
	for rows.Next() {
		pizza := PizzaWithPrice100{}
		if err := rows.Scan(&pizza.Slug, &pizza.Name, &pizza.Price100); err != nil {
			return nil, err
		}
		pizzas = append(pizzas, &pizza)
	}

	return pizzas, nil
}

func (db *DB) GetPizza(slug string) (*PizzaWithPrice100, error) {
	row := db.QueryRow(
		`SELECT p.slug, p.name, SUM(i.price100) 
			FROM pizza p
			JOIN recipe r on p.slug = r.pizza
			JOIN ingredient i on i.slug = r.ingredient
			WHERE p.slug = $1
			GROUP BY p.slug, p.name`,
		slug)
	pizza := PizzaWithPrice100{}
	err := row.Scan(&pizza.Slug, &pizza.Name, &pizza.Price100)
	if err != nil {
		return nil, err
	}

	return &pizza, nil
}

func (db *DB) CreatePizza(pizza *Pizza) error {
	_, err := db.Exec("INSERT INTO pizza (slug, name) VALUES ($1, $2)", pizza.Slug, pizza.Name)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) RemovePizza(slug string) error {
	_, err := db.Exec("DELETE FROM pizza WHERE slug = $1", slug)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetRecipe(pizzaSlug string) (*Recipe, error) {
	rows, err := db.Query("SELECT ingredient FROM recipe WHERE pizza = $1 ORDER BY position", pizzaSlug)
	if err != nil {
		return nil, err
	}
	var ingredients []string
	for rows.Next() {
		var ingredient string
		if err := rows.Scan(&ingredient); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, ingredient)
	}
	if len(ingredients) == 0 {
		return nil, sql.ErrNoRows
	}

	return &Recipe{Ingredients: ingredients}, nil
}

func (db *DB) CreateRecipe(pizzaSlug string, ingredientSlugs []string) error {
	var values string
	for i := 1; i <= len(ingredientSlugs); i++ {
		values += fmt.Sprintf("('%s', '%s', %d), ", pizzaSlug, ingredientSlugs[i-1], i)
	}
	sqlQuery := fmt.Sprintf("INSERT INTO recipe (pizza, ingredient, position) VALUES %s;", values[:len(values)-2])
	fmt.Println(sqlQuery)
	_, err := db.Exec(sqlQuery)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) IngredientExists(ingredientSlug string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM ingredient WHERE slug = $1)", ingredientSlug).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (db *DB) PizzaExists(pizzaSlug string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM pizza WHERE slug = $1)", pizzaSlug).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
