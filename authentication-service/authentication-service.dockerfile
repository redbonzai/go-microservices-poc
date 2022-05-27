FROM alpine:latest
RUN mkdir /app
COPY authApp /app
#COPY --from=builder /app/brokerApp /app

CMD ["/app/authApp"]
