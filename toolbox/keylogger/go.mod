module github.com/kuangcp/gobase/toolbox/keylogger

go 1.16

require (
	github.com/gin-gonic/gin v1.7.2
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/gvalkov/golang-evdev v0.0.0-20191114124502-287e62b94bcb
	github.com/kuangcp/gobase/pkg/ctk v1.0.9
	github.com/kuangcp/gobase/pkg/ghelp v1.0.0
	github.com/kuangcp/gobase/pkg/stopwatch v1.0.1
	github.com/kuangcp/logger v1.0.8
	github.com/manifoldco/promptui v0.8.0
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.19.0 // indirect
	github.com/webview/webview v0.0.0-20210330151455-f540d88dde4e
)

replace github.com/kuangcp/gobase/pkg/ghelp => ./../../pkg/ghelp
