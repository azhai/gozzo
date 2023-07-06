module github.com/azhai/gozzo

go 1.20

require (
	github.com/go-playground/form/v4 v4.2.0
	github.com/gobwas/glob v0.2.3
	github.com/gofiber/fiber/v2 v2.47.0
	github.com/klauspost/cpuid/v2 v2.2.5
	github.com/pkg/errors v0.9.1
	go.uber.org/zap v1.24.0
	golang.org/x/tools v0.11.0
	gorm.io/gorm v1.25.2
	xorm.io/xorm v1.3.2
)

require (
	github.com/stretchr/testify v1.8.2 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/mod v0.12.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
)

replace github.com/gofiber/fiber/v2 => ../fiber
