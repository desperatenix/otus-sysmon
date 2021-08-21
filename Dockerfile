FROM golang:1.16.7-stretch as builder
WORKDIR /sysmon
COPY . .

RUN go build -o /sysmon/sysmon /sysmon/cmd/sysmon2

FROM ubuntu:20.04 as runner
RUN useradd -m sysmon
USER sysmon
WORKDIR /home/sysmon
COPY --from=builder  /sysmon/sysmon /usr/bin/sysmon
EXPOSE 8080
ENTRYPOINT ["/usr/bin/sysmon"]