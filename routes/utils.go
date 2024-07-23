package routes

import (
	"regexp"

	"github.com/gofiber/fiber/v2"
)

// Utility function to convert a route path with parameters to a regex pattern
func RoutePathToRegex(path string) *regexp.Regexp {
	re := regexp.MustCompile(`:[^/]+`)
	regexPath := "^" + re.ReplaceAllString(path, `[^/]+`) + "$"
	return regexp.MustCompile(regexPath)
}

func IsExcludedPath(c *fiber.Ctx, excludedPaths map[string]bool) bool {
		path := c.Path()
		for k, _ := range excludedPaths {
			regex := RoutePathToRegex(k)
			if regex.MatchString(path) {
				return true
			}
		}
		return false
	}