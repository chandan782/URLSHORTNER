package routes

import (
	"urlshortner/api/database"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")
	redisClient := database.CreateRedisClient(0)
	defer redisClient.Close()

	value, err := redisClient.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "short not found in the database",
		})
	}

	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot connect to db",
		})
	}

	newRedisClient := database.CreateRedisClient(1)
	defer newRedisClient.Close()

	_ = newRedisClient.Incr(database.Ctx, "counter")

	return c.Redirect(value, 301)
}
