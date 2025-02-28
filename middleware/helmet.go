package middleware

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber"
)

// HelmetConfig ...
type HelmetConfig struct {
	// Skip defines a function to skip middleware.
	// Optional. Default: nil
	Skip func(*fiber.Ctx) bool
	// XSSProtection
	// Optional. Default value "1; mode=block".
	XSSProtection string
	// ContentTypeNosniff
	// Optional. Default value "nosniff".
	ContentTypeNosniff string
	// XFrameOptions
	// Optional. Default value "SAMEORIGIN".
	// Possible values: "SAMEORIGIN", "DENY", "ALLOW-FROM uri"
	XFrameOptions string
	// HSTSMaxAge
	// Optional. Default value 0.
	HSTSMaxAge int
	// HSTSExcludeSubdomains
	// Optional. Default value false.
	HSTSExcludeSubdomains bool
	// ContentSecurityPolicy
	// Optional. Default value "".
	ContentSecurityPolicy string
	// CSPReportOnly
	// Optional. Default value false.
	CSPReportOnly bool
	// HSTSPreloadEnabled
	// Optional.  Default value false.
	HSTSPreloadEnabled bool
	// ReferrerPolicy
	// Optional. Default value "".
	ReferrerPolicy string
}

// HelmetConfigDefault is the defaul Helmet middleware config.
var HelmetConfigDefault = HelmetConfig{
	Skip:               nil,
	XSSProtection:      "1; mode=block",
	ContentTypeNosniff: "nosniff",
	XFrameOptions:      "SAMEORIGIN",
}

// Helmet ...
func Helmet(config ...HelmetConfig) func(*fiber.Ctx) {
	log.Println("Warning: middleware.Helmet() is deprecated since v1.8.3, please use github.com/gofiber/helmet")
	// Init config
	var cfg HelmetConfig
	if len(config) > 0 {
		cfg = config[0]
	}
	// Set config default values
	if cfg.XSSProtection == "" {
		cfg.XSSProtection = HelmetConfigDefault.XSSProtection
	}
	if cfg.ContentTypeNosniff == "" {
		cfg.ContentTypeNosniff = HelmetConfigDefault.ContentTypeNosniff
	}
	if cfg.XFrameOptions == "" {
		cfg.XFrameOptions = HelmetConfigDefault.XFrameOptions
	}
	// Return middleware handler
	return func(c *fiber.Ctx) {
		// Skip middleware if Skip returns true
		if cfg.Skip != nil && cfg.Skip(c) {
			c.Next()
			return
		}
		if cfg.XSSProtection != "" {
			c.Set(fiber.HeaderXXSSProtection, cfg.XSSProtection)
		}
		if cfg.ContentTypeNosniff != "" {
			c.Set(fiber.HeaderXContentTypeOptions, cfg.ContentTypeNosniff)
		}
		if cfg.XFrameOptions != "" {
			c.Set(fiber.HeaderXFrameOptions, cfg.XFrameOptions)
		}
		if (c.Secure() || (c.Get(fiber.HeaderXForwardedProto) == "https")) && cfg.HSTSMaxAge != 0 {
			subdomains := ""
			if !cfg.HSTSExcludeSubdomains {
				subdomains = "; includeSubdomains"
			}
			if cfg.HSTSPreloadEnabled {
				subdomains = fmt.Sprintf("%s; preload", subdomains)
			}
			c.Set(fiber.HeaderStrictTransportSecurity, fmt.Sprintf("max-age=%d%s", cfg.HSTSMaxAge, subdomains))
		}
		if cfg.ContentSecurityPolicy != "" {
			if cfg.CSPReportOnly {
				c.Set(fiber.HeaderContentSecurityPolicyReportOnly, cfg.ContentSecurityPolicy)
			} else {
				c.Set(fiber.HeaderContentSecurityPolicy, cfg.ContentSecurityPolicy)
			}
		}
		if cfg.ReferrerPolicy != "" {
			c.Set(fiber.HeaderReferrerPolicy, cfg.ReferrerPolicy)
		}
		c.Next()
	}
}
