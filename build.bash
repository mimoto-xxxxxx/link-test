PLATFORMS="darwin/386 darwin/amd64 freebsd/386 freebsd/amd64 freebsd/arm linux/386 linux/amd64 linux/arm windows/386 windows/amd64"
eval "$(go env)"

for PLATFORM in $PLATFORMS; do
	GOOS=${PLATFORM%/*}
	GOARCH=${PLATFORM#*/}
	OUTPUT=`echo $@ | sed 's/\.go//'` 

	if [ "$GOOS" = "windows" ]; then
    EXT=".exe"
  else
    EXT=""
	fi

	CMD="GOOS=${GOOS} GOARCH=${GOARCH} go build -o build/link-test-${GOOS}-${GOARCH}${EXT} $@"
  echo "$CMD"
  eval "$CMD"
done
