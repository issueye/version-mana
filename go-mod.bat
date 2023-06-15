rmdir /s/q vendor

:: 更新依赖
go mod tidy
:: 更新本地依赖
go mod vendor