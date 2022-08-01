FROM golang:1.18.1-buster as builder


WORKDIR /app
COPY . ./

# running linters
RUN make deps

# build artifact
RUN make build

FROM debian:buster-slim

# Copy our static executable
COPY --from=builder /app/artifacts/svc /svc

# Port on which the core will be exposed.
EXPOSE 8081

#Run Container as nonroot
USER nobody

# Run the svc binary.
CMD ["./svc"]
