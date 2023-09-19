module github.com/Eden/go-gin-example

go 1.20

require (
	github.com/astaxie/beego v1.10.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fvbock/endless v0.0.0-20170109170031-447134032cb6
	github.com/gin-gonic/gin v1.9.1
	github.com/go-ini/ini v1.67.0
	github.com/gomodule/redigo v1.8.9
	github.com/jinzhu/gorm v1.9.16
	github.com/robfig/cron v1.2.0
	github.com/unknwon/com v1.0.1
)

require (
	github.com/bytedance/sonic v1.10.0 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
	github.com/chenzhuoyu/iasm v0.9.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.15.1 // indirect
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.9 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	golang.org/x/arch v0.4.0 // indirect
	golang.org/x/crypto v0.12.0 // indirect
	golang.org/x/net v0.14.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/Eden/go-gin-example/conf => ../go-gin-example/pkg/conf
	github.com/Eden/go-gin-example/middleware => ../go-gin-example/middleware
	github.com/Eden/go-gin-example/models => ../go-gin-example/models
	github.com/Eden/go-gin-example/pkg/e => ../go-gin-example/pkg/e
	github.com/Eden/go-gin-example/pkg/setting => ../go-gin-example/pkg/setting
	github.com/Eden/go-gin-example/pkg/util => ../go-gin-example/pkg/util
	github.com/Eden/go-gin-example/routers => ../go-gin-example/routers
	github.com/Eden/go-gin-example/service => ../go-gin-example/service
)
