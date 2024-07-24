package routes

import (
	"first/controllers/users"
	"first/guards"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Router struct {
	app *fiber.App
	UserController *users.UserController
}

func NewRouter(user_controller *users.UserController) *Router {
	app := fiber.New()	
	
	// Middleware to check for excluded paths
	excludedPaths := map[string]bool{
		"/user": true,
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	//Setting multiple HTTP security headers
	app.Use(guards.Security)

	//logger https://docs.gofiber.io/api/middleware/logger
	app.Use(logger.New(logger.Config{
    	Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	//Cache layer using fiber inbuilt middleware using in-memory storage
	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return IsExcludedPath(c, excludedPaths) || c.Query("noCache") == "true"
		},
		Expiration: 10 * time.Second,
		CacheControl: true,
	}))

	//Rate limitting 100 requests per 10 seconds max
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return IsExcludedPath(c, excludedPaths)
		},
		Expiration: 10 * time.Second,
		Max:      100,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	return &Router{
		app: app,
		UserController: user_controller,
	}
}

func (router *Router) RegisterUserRoutes() {
	users := router.app.Group("/user")
	users.Post("/", router.UserController.CreateNewUser)
	users.Use(guards.JwtAuthGuard)
	users.Get("/protected", router.UserController.ProtectedUser)
	users.Get("/all", router.UserController.GetAllUserPaginated)
	users.All("/*", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "not found",
		})
	})
}

func (router *Router) RegisterRoutes() {
	router.RegisterUserRoutes();
	fmt.Println("ROUTES:= All Routes registered")
}

func (router *Router) StartServer() {
	router.RegisterUserRoutes()
	
	//for all unregistered route
	router.app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"code":    404,
			"message": "404: Not Found",
		})
	})

	router.app.Listen((os.Getenv("APP_PORT")))
}