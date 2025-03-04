FROM golang:1.24

# Build the app
COPY . /build
WORKDIR /build
RUN go install
RUN go build -o login-service

# Create the /app directory
RUN mkdir -p /app

# Copy the built application to the /app directory
RUN cp /build/login-service /app/login-service

# Remove the build directory
RUN rm -rf /build

WORKDIR /app

# Generate the RSA keys
RUN openssl genrsa -out app.rsa 3072
RUN openssl rsa -in app.rsa -pubout > app.rsa.pub

USER 1000:1000

CMD ["./login-service"]