package routes

import (
	"first/controllers"
	"fmt"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

type Router struct {
	app *fiber.App
	UserController *controllers.UserController
}

func NewRouter(user_controller *controllers.UserController) *Router {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	//Setting multiple HTTP security headers
	app.Use(helmet.New())
	
	// Middleware to check for excluded paths
	excludedPaths := map[string]bool{
		"/api/v1/user/:name": true,
		"/api/v1/no-rate-limit": true,
		"/api/v1/no-cache": true,
	}

	isExcludedPath := func(c *fiber.Ctx) bool {
		path := c.Path()
		for k, v := range excludedPaths {
			fmt.Print(v)
			regex := RoutePathToRegex(k)
			if regex.MatchString(path) {
				return true
			}
		}
		return false
	}

	//Cache layer using fiber inbuilt middleware using in-memory storage
	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return isExcludedPath(c) || c.Query("noCache") == "true"
		},
		Expiration: 30 * time.Minute,
		CacheControl: true,
	}))

	//Rate limitting 20 requests per 10 seconds max
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return isExcludedPath(c)
		},
		Expiration: 10 * time.Second,
		Max:      20,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	return &Router{
		app: app,
		UserController: user_controller,
	}
}

// Utility function to convert a route path with parameters to a regex pattern
func RoutePathToRegex(path string) *regexp.Regexp {
	re := regexp.MustCompile(`:[^/]+`)
	regexPath := "^" + re.ReplaceAllString(path, `[^/]+`) + "$"
	return regexp.MustCompile(regexPath)
}

func (router *Router) RegisterUserRoutes() {
	router.app.Post("/user", router.UserController.CreateNewUser)
}

func (router *Router) RegisterRoutes() {
	router.RegisterUserRoutes();
	fmt.Println("Routes registered")
}

func (router *Router) StartServer() {
	router.RegisterUserRoutes()
	router.app.Listen((":5000"))
	fmt.Println("API Server running on :5000")
}