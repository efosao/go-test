package main

import (
	"gofiber/models"
	"strconv"
	"strings"

	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/mustache/v2"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Pacific73/gorm-cache/cache"
	"github.com/Pacific73/gorm-cache/config"
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
		// if length of objects retrieved one single time
		// exceeds this number, don't cache
		CacheMaxItemCnt: 20,
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
	})

	app.Use(compress.New())
	app.Use(helmet.New(
		helmet.Config{
			ContentSecurityPolicy: `
				default-src 'self';
				base-uri 'self';
				font-src 'self' https: data:;
				form-action 'self';
				frame-ancestors 'self';
				img-src 'self' https: data:;
				object-src 'none';
				script-src 'self' unpkg.com cdn.jsdelivr.net 'unsafe-eval' 'unsafe-inline';
				script-src-attr 'unsafe-inline';
				style-src 'self' https: 'unsafe-inline';
				upgrade-insecure-requests
			`,
		}))
	app.Use(etag.New())
	app.Use(recover.New())

	app.Use(func(c *fiber.Ctx) error {
		if isGet := (c.Method() == "GET"); isGet {
			theme := c.Cookies("theme")
			themeOptions := []models.ThemeOption{}
			themeOptions = append(themeOptions, models.ThemeOption{Value: "light", Label: "🌞", Selected: theme == "light"})
			themeOptions = append(themeOptions, models.ThemeOption{Value: "dark", Label: "🌘", Selected: theme == "dark"})
			themeOptions = append(themeOptions, models.ThemeOption{
				Value:    "system",
				Label:    "🌎",
				Selected: (theme != "light" && theme != "dark"),
			})
			c.Locals("ThemeOptions", themeOptions)
		}
		return c.Next()
	})

	app.Static("/public", "./public", fiber.Static{
		Compress:      true,
		CacheDuration: 60 * time.Second,
		ByteRange:     true,
		ModifyResponse: func(ctx *fiber.Ctx) error {
			ctx.Set(fiber.HeaderCacheControl, fmt.Sprintf("private, max-age=%d, stale-while-revalidate=3600", 3600))
			return nil
		},
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title":        "Hello, World!",
			"Description":  "Find the latest job posts in the tech industry.",
			"ThemeOptions": c.Locals("ThemeOptions"),
		}, "layouts/main")
	})

	app.Get("/monitor", monitor.New())

	app.Get("/posts", func(c *fiber.Ctx) error {
		cookie := new(models.Cookie)
		postsChan := make(chan []models.Post)
		tagsChan := make(chan []models.Tag)
		selectedTagsStr := c.Query("tags")
		selectedTags := strings.Split(selectedTagsStr, ",")

		go (func(p chan []models.Post) {
			posts := []models.Post{}
			if len(selectedTags) > 0 {
				queryInputTags := "{" + strings.Join(selectedTags, ",") + "}"
				db.Select("ID", "CompanyName", "Location", "Tags", "Thumbnail", "Title", "PublishedAt", "CreatedAt").Where("tags @> ?", queryInputTags).Where("published_at IS NOT NULL").Order(clause.OrderByColumn{Column: clause.Column{Name: "published_at"}, Desc: true}).Limit(10).Find(&posts)
				p <- posts
				return
			} else {
				db.Select("ID", "CompanyName", "Location", "Tags", "Thumbnail", "Title", "PublishedAt", "CreatedAt").Where("published_at IS NOT NULL").Order(clause.OrderByColumn{Column: clause.Column{Name: "published_at"}, Desc: true}).Limit(10).Find(&posts)
				p <- posts
			}
		})(postsChan)

		go (func(t chan []models.Tag) {
			tags := []models.Tag{}
			db.Raw(`
				SELECT unnest(tags) AS name, count(*)::text AS count
				FROM posts
				WHERE published_at IS NOT NULL
				GROUP by name
				ORDER BY count(*) DESC;
			`).Scan(&tags)
			t <- tags
		})(tagsChan)

		posts := <-postsChan
		tags := <-tagsChan

		selectedTagMap := map[string]bool{}
		for _, selectedTag := range selectedTags {
			selectedTagMap[selectedTag] = true
		}

		updatedTags := []models.Tag{}
		for _, tag := range tags {
			if selectedTagMap[tag.Name] {
				tag.Selected = true
			}
			updatedTags = append(updatedTags, tag)
		}

		if err := c.CookieParser(cookie); err != nil {
			return err
		}

		return c.Render("posts", fiber.Map{
			"Description":     "Find the latest job posts in the tech industry.",
			"Page":            "1",
			"Posts":           posts,
			"SelectedTagsStr": selectedTagsStr,
			"Theme":           cookie.Theme,
			"Tags":            updatedTags,
			"ThemeOptions":    c.Locals("ThemeOptions"),
			"Title":           "Job Posts",
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

	app.Post("/partials/posts/search/:page", func(c *fiber.Ctx) error {
		type Body struct{ Tags []string }
		if pageStr, err := strconv.Atoi(c.Params("page", "0")); err != nil {
			return err
		} else {
			body := Body{}
			c.BodyParser(&body)
			page := int(pageStr)
			nextPage := page + 1
			offset := (nextPage - 1) * 10
			posts := []models.Post{}
			selectedTagsStr := ""
			if len(body.Tags) > 0 {
				selectedTagsStr = strings.Join(body.Tags, ",")
			} else {
				selectedTagsStr = c.Query("tags")
			}

			queryInputTags := "{" + selectedTagsStr + "}"

			if selectedTagsStr == "" {
				if page == 0 {
					c.Response().Header.Set("HX-Push-Url", "/posts")
				}
				db.Select("ID", "CompanyName", "Location", "Tags", "Thumbnail", "Title", "PublishedAt", "CreatedAt").Where("published_at IS NOT NULL").Order(clause.OrderByColumn{Column: clause.Column{Name: "published_at"}, Desc: true}).Offset(offset).Limit(10).Find(&posts)

				if len(posts) == 0 {
					return c.Send([]byte(""))
				}

				return c.Render("post_list", fiber.Map{
					"Posts": posts,
					"Page":  nextPage,
				})
			} else {
				if page == 0 {
					c.Response().Header.Set("HX-Push-Url", fmt.Sprintf("/posts?tags=%s", selectedTagsStr))
				}
				db.Select("ID", "CompanyName", "Location", "Tags", "Thumbnail", "Title", "PublishedAt", "CreatedAt").Where("tags @> ?", queryInputTags).Where("published_at IS NOT NULL").Order(clause.OrderByColumn{Column: clause.Column{Name: "published_at"}, Desc: true}).Offset(offset).Limit(10).Find(&posts)

				if len(posts) == 0 {
					return c.Send([]byte(""))
				}

				return c.Render("post_list", fiber.Map{
					"Posts":           posts,
					"Page":            nextPage,
					"SelectedTagsStr": selectedTagsStr,
				})
			}
		}
	})

	PORT := "localhost:8000"
	if os.Getenv("PORT") != "" {
		PORT = fmt.Sprintf(":%s", os.Getenv("PORT"))
	}

	log.Fatal(app.Listen(PORT))
}
