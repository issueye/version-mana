export CGO_ENABLED=0
export GOARCH=amd64
export GOOS=linux

go mod tidy

go mod vendor

TAG="Beta"
BRANCH=$(git symbolic-ref --short -q HEAD)
COMMIT=$(git rev-parse --verify HEAD)
NOW=$(date '+%FT%T%z')

VERSION="v0.1.1-${TAG}"
APPNAME="VersionMana-${VERSION}"
DESCRIPTION="版本管理服务"

go build -o bin/${APPNAME} -tags=ui -ldflags "-X demo/build.AppName=Demo \
-X github.com/issueye/version-mana/internal/initialize.Branch=${BRANCH} \
-X github.com/issueye/version-mana/internal/initialize.Commit=${COMMIT} \
-X github.com/issueye/version-mana/internal/initialize.Date=${NOW} \
-X github.com/issueye/version-mana/internal/initialize.AppName=${DESCRIPTION} \
-X github.com/issueye/version-mana/internal/initialize.Version=${VERSION}" .