package main

import (
	"gofiber/models"
	"log"

	"github.com/Pacific73/gorm-cache/cache"
	"github.com/Pacific73/gorm-cache/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/template/mustache/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	engine := mustache.New("./views", ".mustache")

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
		CacheMaxItemCnt:      5,     // if length of objects retrieved one single time
		// exceeds this number, then don't cache
	})

	db.Use(cache)

	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(30)

	if err != nil {
		panic("failed to connect database")
	}

	// Or from an embedded system
	//   Note that with an embedded system the partials included from template files must be
	//   specified relative to the filesystem's root, not the current working directory
	// engine := mustache.NewFileSystem(http.Dir("./views", ".mustache"), ".mustache")

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(compress.New())

	app.Get("/", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		}, "layouts/main")
	})

	app.Get("/posts", func(c *fiber.Ctx) error {
		// readPost := &Post{}
		// db.First(&readPost, "id = ?", "cll03wjh700hvo485rzqiyjy2")
		// // Render index within layouts/main
		// return c.Render("index", fiber.Map{
		// 	"Title": readPost.Title,
		// }, "layouts/main")

		posts := &[]models.Post{}
		db.Select("Id", "Title", "CreatedAt").Limit(10).Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true}).Find(&posts)
		return c.Render("posts", fiber.Map{
			"Title": "Job Posts",
			"Posts": posts,
		}, "layouts/main")
	})

	app.Get("/single", func(c *fiber.Ctx) error {
		// Render w/o layout
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	log.Fatal(app.Listen(":8000"))
}
