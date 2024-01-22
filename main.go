package main

import (
	"fmt"
	"gofiber/models"
	"log"
	"os"
	"time"

	"github.com/Pacific73/gorm-cache/cache"
	"github.com/Pacific73/gorm-cache/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/template/mustache/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	dsn := "host=efosa.me user=postgres password=5005227a52c02361b7e95a1f5acfc7f0 dbname=jobs_db port=44553 sslmode=disable TimeZone=America/Los_Angeles"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	cache, _ := cache.NewGorm2Cache(&config.CacheConfig{
		CacheLevel:           config.CacheLevelAll,
		CacheStorage:         config.CacheStorageMemory,
		InvalidateWhenUpdate: true,  // when u create/update/delete objects, invalidate cache
		CacheTTL:             10000, // 5000 ms
		CacheMaxItemCnt:      500,   // if length of objects retrieved one single time
		// exceeds this number, then don't cache
	})
	db.Use(cache)

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(30)
	engine := mustache.New("./views", ".mustache")

	app := fiber.New(fiber.Config{
		Views: engine,
		ETag:  true,
	})

	app.Use(compress.New())
	app.Use(helmet.New())

	app.Static("/public", "./public", fiber.Static{
		Compress:      true,
		CacheDuration: 10 * time.Second,
		MaxAge:        3600,
		ByteRange:     true,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title":       "Hello, World!",
			"Description": "Find the latest job posts in the tech industry.",
		}, "layouts/main")
	})

	app.Get("/monitor", monitor.New())

	app.Get("/posts", func(c *fiber.Ctx) error {
		posts := &[]models.Post{}
		db.Select("ID", "CompanyName", "Location", "Tags", "Thumbnail", "Title", "PublishedAt", "CreatedAt").Limit(10).Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true}).Find(&posts)

		return c.Render("posts", fiber.Map{
			"Title":       "Job Posts",
			"Posts":       posts,
			"Description": "Find the latest job posts in the tech industry.",
		}, "layouts/main")
	})

	app.Get("/partials/posts/details/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		post := &models.Post{}
		db.Select("ID", "Title", "Description").Where(&models.Post{ID: id}).First(&post)
		return c.Render("post_details", fiber.Map{
			"Title":       post.Title,
			"ID":          post.ID,
			"Description": post.GetDescription(),
		})
	})

	PORT := "localhost:8000"
	if os.Getenv("PORT") != "" {
		PORT = fmt.Sprintf(":%s", os.Getenv("PORT"))
	}

	log.Fatal(app.Listen(PORT))
}
