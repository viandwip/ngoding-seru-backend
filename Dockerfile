FROM golang:1.22.2-alpine AS build

#
WORKDIR /goapp

# copy semua file di lokasi saat ini(.) lalu paste ke workdir
COPY . .

# download depedency
RUN go mod download
# build file golang
RUN  go build -v -o /goapp/musalabel ./cmd/main.go

FROM alpine:3.14

WORKDIR /app
# copy isi dari goapp lalu paste ke app
COPY --from=build /goapp /app/

#karna di eksekusi di linux(alpine) harus daftarin path variable dahulu
ENV PATH="/app:${PATH}"

EXPOSE 9090

ENTRYPOINT [ "musalabel" ]
