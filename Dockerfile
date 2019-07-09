## == builder image ==

FROM cortezaproject/corteza-server-builder:latest AS builder

WORKDIR /go/src/github.com/crusttech/crust-bundle

COPY . .

RUN scripts/builder-make-bin.sh

## == webapp image ==

# @todo using corteza for now
# FROM crusttech/webapp:latest as webapp
FROM cortezaproject/corteza-webapp:latest as webapp

## == target image ==

FROM alpine:3.7

RUN apk add --no-cache ca-certificates

COPY --from=builder /bin/crust-bundle /bin

# @todo using corteza for now
#COPY --from=webapp /crust/webapp /webapp
COPY --from=webapp /webapp /webapp

EXPOSE 80
ENTRYPOINT ["/bin/crust-bundle"]
CMD ["serve"]

