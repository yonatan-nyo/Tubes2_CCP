# Stage 1: Build Vite frontend
FROM node:22.15.0 AS frontend-builder
WORKDIR /app/frontend
COPY src/frontend/package*.json ./
RUN npm install
COPY src/frontend/ ./
RUN npm run build

# Stage 2: Build Go backend
FROM golang:1.24.2 AS backend-builder
WORKDIR /app/backend
COPY src/backend/go.* ./
RUN go mod download
COPY src/backend/ ./
RUN go build -o backend-main .

# Stage 3: Build Go scraper
FROM golang:1.24.2 AS scraper-builder
WORKDIR /app/scraper
COPY src/scraper/go.* ./
RUN go mod download
COPY src/scraper/ ./
RUN go build -o scraper-main .

# Stage 4: Final image
FROM node:22.15.0 AS final
WORKDIR /app

# Copy required runtime data for the backend
COPY src/backend/data ./data
COPY src/backend/public ./public
# Copy Go binaries
COPY --from=backend-builder /app/backend/backend-main ./backend
COPY --from=scraper-builder /app/scraper/scraper-main ./scraper

# Copy supervisor config
RUN npm install -g concurrently

# Script to run both Go binaries and frontend server
COPY start.sh .

RUN chmod +x ./start.sh

EXPOSE 4000

CMD ["./start.sh"]
