# Build from the repo root: docker build -f api/Dockerfile .
# The Go compile in the builder stage is the workload that hurts under QEMU
# emulation on the arm64 leg of a multi-platform build.
FROM golang:1.25 AS builder
WORKDIR /src
COPY api/go.mod api/go.sum ./
RUN go mod download
COPY api/ ./
# -a forces a full rebuild of all packages so the compile cost is honest
# and consistent across runs (no host toolchain cache).
RUN CGO_ENABLED=0 go build -a -o /confab-api .

FROM gcr.io/distroless/static-debian12
COPY --from=builder /confab-api /confab-api
EXPOSE 3000
ENTRYPOINT ["/confab-api"]
