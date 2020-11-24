package pharm

import (
	"github.com/gofiber/fiber/v2"
	"github.com/urfave/cli/v2"
	"github.com/warete/pharm/cmd/pharm/database"
	"github.com/warete/pharm/cmd/pharm/models/product"

	log "github.com/sirupsen/logrus"
	"strconv"
)

var app *fiber.App

func Init() {
	app = fiber.New()

	err := database.Init("bin/pharm.db")
	if err != nil {
		log.Fatal(err)
	}

	database.DB.Connection.AutoMigrate(&product.Product{})
}

func Cmd(c *cli.Context) error {
	api := app.Group("/api/v0.1/")

	api.Get("/products", func(c *fiber.Ctx) error {
		products, _ := product.GetAll()
		return c.JSON(products)
	})
	api.Get("/products/:id", func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		prod, _ := product.GetById(id)
		return c.JSON(prod)
	})
	api.Post("/products", func(c *fiber.Ctx) error {
		prod := new(product.Product)
		if err := c.BodyParser(prod); err != nil {
			return c.Status(503).SendString("error")
		}
		product.Add(prod)
		return c.JSON(prod)
	})

	app.Listen(":3000")

	return nil
}
