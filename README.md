# melancholy

🐳Script generation tool for multi-stage build using cache with CI which is docker in docker

🚨Assuming use of BuildKit

## Installation

```sh
$ go get -u github.com/kotaroooo0/melancholy
```

## How to use

```sh
$ cat Dockerfile
FROM alpine:latest as stage1
WORKDIR /workdir
RUN touch stage1.txt
FROM alpine:latest as stage2
WORKDIR /workdir
RUN touch stage2.txt
FROM alpine:latest
WORKDIR /workdir
COPY --from=stage1 /workdir/stage1.txt /workdir/stage1.txt
COPY --from=stage2 /workdir/stage2.txt /workdir/stage2.txt

$ melancholy -i image_name
# ----- Build image -----
docker build -t image_name:latest --cache-from=image_name:stage1,image_name:stage2,image_name:latest --build-arg BUILDKIT_INLINE_CACHE=1 .
# ----- Attach tags -----
docker build -t image_name:stage1 --target=stage1 --build-arg BUILDKIT_INLINE_CACHE=1 . &
docker build -t image_name:stage2 --target=stage2 --build-arg BUILDKIT_INLINE_CACHE=1 . &
wait
# ----- Push images -----
docker push image_name:stage1 &
docker push image_name:stage2 &
docker push image_name:latest &
wait
```

## Author

Kotaro Adachi (@kotaroooo0)

## License

MIT
