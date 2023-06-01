package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"src/database"
)

type Env struct {
	db *database.DB
}

type ErrorMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func getErrorMessageFromTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	}
	return "Unknown error"
}

type PostIngredientInput struct {
	Slug     string `json:"slug" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Price100 int    `json:"price100" binding:"required"`
}

type PutIngredientInput struct {
	Name     string `json:"name" binding:"required"`
	Price100 int    `json:"price100" binding:"required"`
}

type PostPizzaInput struct {
	Slug string `json:"slug" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type PostRecipeInput struct {
	PizzaSlug   string   `json:"pizza_slug" binding:"required"`
	Ingredients []string `json:"ingredients" binding:"required"`
}

func (e *Env) getIngredients(c *gin.Context) {
	ingredients, err := e.db.GetIngredients()
	if err != nil {
		log.Fatalf("failed to retrieve ingredients: %v", err)
	}
	c.IndentedJSON(http.StatusOK, ingredients)
}

func (e *Env) getIngredient(c *gin.Context) {
	slug := c.Param("slug")
	ingredient, err := e.db.GetIngredient(slug)
	if err != nil {
		// Check if the error is a "not rows" error
		if errors.Is(err, sql.ErrNoRows) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Ingredient not found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Fatalf("failed to get ingredient: %v", err)
		}
		return
	}
	c.IndentedJSON(http.StatusOK, ingredient)
}

func (e *Env) createIngredient(c *gin.Context) {
	var newIngredient PostIngredientInput
	if err := c.BindJSON(&newIngredient); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ErrorMessage, len(ve))
			for i, fe := range ve {
				out[i] = ErrorMessage{fe.Field(), getErrorMessageFromTag(fe.Tag())}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
		}
		return
	}
	ingredient := database.Ingredient{
		Slug:     newIngredient.Slug,
		Name:     newIngredient.Name,
		Price100: newIngredient.Price100,
	}
	err := e.db.CreateIngredient(&ingredient)
	if err != nil {
		log.Fatalf("failed to create ingredient: %v", err)
	}
	c.IndentedJSON(http.StatusCreated, newIngredient)
}

func (e *Env) removeIngredient(c *gin.Context) {
	slug := c.Param("slug")
	err := e.db.RemoveIngredient(slug)
	if err != nil {
		log.Fatalf("failed to remove ingredient: %v", err)
	}
	c.IndentedJSON(http.StatusNoContent, nil)
}

func (e *Env) updateIngredient(c *gin.Context) {
	slug := c.Param("slug")
	var newIngredient PutIngredientInput
	if err := c.BindJSON(&newIngredient); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ErrorMessage, len(ve))
			for i, fe := range ve {
				out[i] = ErrorMessage{fe.Field(), getErrorMessageFromTag(fe.Tag())}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
		}
		return
	}
	ingredient := database.Ingredient{
		Slug:     slug,
		Name:     newIngredient.Name,
		Price100: newIngredient.Price100,
	}
	updatedIngredient, err := e.db.UpdateIngredient(&ingredient)
	if err != nil {
		// Check if the error is a "not rows" error
		if errors.Is(err, sql.ErrNoRows) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Fatalf("failed to update ingredient: %v", err)
		}
		return
	}
	c.IndentedJSON(http.StatusOK, updatedIngredient)
}

func (e *Env) getPizzas(c *gin.Context) {
	pizzas, err := e.db.GetPizzas()
	if err != nil {
		log.Fatalf("failed to retrieve pizzas: %v", err)
	}
	c.IndentedJSON(http.StatusOK, pizzas)
}

func (e *Env) getPizza(c *gin.Context) {
	slug := c.Param("slug")
	pizza, err := e.db.GetPizza(slug)
	if err != nil {
		// Check if the error is a "not rows" error
		if errors.Is(err, sql.ErrNoRows) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Pizza not found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Fatalf("failed to get pizza: %v", err)
		}
		return
	}
	c.IndentedJSON(http.StatusOK, pizza)
}

func (e *Env) createPizza(c *gin.Context) {
	var newPizza PostPizzaInput
	if err := c.BindJSON(&newPizza); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ErrorMessage, len(ve))
			for i, fe := range ve {
				out[i] = ErrorMessage{fe.Field(), getErrorMessageFromTag(fe.Tag())}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
		}
		return
	}
	pizza := database.Pizza{
		Slug: newPizza.Slug,
		Name: newPizza.Name,
	}
	err := e.db.CreatePizza(&pizza)
	if err != nil {
		log.Fatalf("failed to create pizza: %v", err)
	}
	c.IndentedJSON(http.StatusCreated, newPizza)
}

func (e *Env) removePizza(c *gin.Context) {
	slug := c.Param("slug")
	err := e.db.RemovePizza(slug)
	if err != nil {
		log.Fatalf("failed to remove pizza: %v", err)
	}
	c.IndentedJSON(http.StatusNoContent, nil)
}

func (e *Env) getRecipe(c *gin.Context) {
	pizzaSlug := c.Param("pizza_slug")
	recipe, err := e.db.GetRecipe(pizzaSlug)
	if err != nil {
		// Check if the error is a "not rows" error
		if errors.Is(err, sql.ErrNoRows) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Fatalf("failed to get recipe: %v", err)
		}
		return
	}
	c.IndentedJSON(http.StatusOK, recipe)
}

func (e *Env) createRecipe(c *gin.Context) {
	var newRecipe PostRecipeInput
	if err := c.BindJSON(&newRecipe); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ErrorMessage, len(ve))
			for i, fe := range ve {
				out[i] = ErrorMessage{fe.Field(), getErrorMessageFromTag(fe.Tag())}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
		}
		return
	}
	// check the provided pizza TODO: refaire
	pizzaSlug := newRecipe.PizzaSlug
	exists, err := e.db.PizzaExists(pizzaSlug)
	if err != nil {
		fmt.Printf("Error checking if pizza %s exists: %v\n", pizzaSlug, err)
		return
	}
	if !exists {
		out := fmt.Sprintf("Provided pizza '%s' does not exist", pizzaSlug)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
		return
	}
	// check the provided list of ingredients
	for _, ingredientSlug := range newRecipe.Ingredients {
		exists, err := e.db.IngredientExists(ingredientSlug)
		if err != nil {
			fmt.Printf("Error checking if ingredient %s exists: %v\n", ingredientSlug, err)
			return
		}
		if !exists {
			out := fmt.Sprintf("Provided ingredient '%s' does not exist", ingredientSlug)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
			return
		}
	}
	// check if recipe already exists
	existingRecipe, err := e.db.GetRecipe(newRecipe.PizzaSlug)
	if existingRecipe != nil {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"errors": "A recipe already exists for this pizza"})
		return
	}

	err = e.db.CreateRecipe(newRecipe.PizzaSlug, newRecipe.Ingredients)
	if err != nil {
		log.Fatalf("failed to create recipe: %v", err)
	}
	c.IndentedJSON(http.StatusCreated, database.Recipe{Ingredients: newRecipe.Ingredients})
}

func main() {
	db, err := database.OpenDBConnection("postgresql://root:root@localhost?sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
	env := &Env{db: db}
	router := gin.Default()

	// routes for ingredients
	router.GET("/ingredients", env.getIngredients)
	router.GET("/ingredients/:slug", env.getIngredient)
	router.POST("/ingredients", env.createIngredient)
	router.DELETE("/ingredients/:slug", env.removeIngredient)
	router.PUT("/ingredients/:slug", env.updateIngredient)

	// routes for pizzas
	router.GET("/pizzas", env.getPizzas)
	router.GET("/pizzas/:slug", env.getPizza)
	router.POST("/pizzas", env.createPizza)
	router.DELETE("/pizzas/:slug", env.removePizza)

	// routes for recipes
	router.GET("/recipes/:pizza_slug", env.getRecipe)
	router.POST("/recipes", env.createRecipe)

	router.Run(":3231")
}
