# Build from the repo root: docker build -f api/Dockerfile .
# The Go compile in the builder stage is the expensive layer; on Blacksmith
# it persists in the Docker layer cache between runs.
FROM golang:1.26 AS builder
WORKDIR /src
COPY api/go.mod api/go.sum ./
RUN go mod download
COPY api/ ./
# -a forces a full rebuild of all packages so the compile cost is honest
# and consistent across runs (no host toolchain cache).
RUN CGO_ENABLED=0 go build -a -o /blacksmith-demo-api .

FROM gcr.io/distroless/static-debian12
COPY --from=builder /blacksmith-demo-api /blacksmith-demo-api
EXPOSE 3000
ENTRYPOINT ["/blacksmith-demo-api"]
