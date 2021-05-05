FROM alpine

RUN mkdir -p /app/bin
ADD ./bin/notifier /app/bin/
RUN chmod +x /app/bin/notifier
WORKDIR /app

ENTRYPOINT ["/app/bin/notifier"]