module github.com/kuangcp/gobase/mybook

go 1.13

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/jinzhu/gorm v1.9.12
	github.com/jroimartin/gocui v0.4.0
	github.com/kuangcp/gobase/cuibase v0.0.0-20200523175218-066ce23005b2
	github.com/mattn/go-runewidth v0.0.8 // indirect
	github.com/mattn/go-sqlite3 v2.0.2+incompatible
	github.com/nsf/termbox-go v0.0.0-20200204031403-4d2b513ad8be // indirect
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.4.0
	github.com/wonderivan/logger v1.0.0
)

replace github.com/kuangcp/gobase/cuibase => ./../cuibase
