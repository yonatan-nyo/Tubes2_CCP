FROM node:22.15.0

WORKDIR /app

COPY . .

RUN npm install

# Build the Vite app
RUN npm run build

EXPOSE 4001

# Use serve to host the static build
RUN npm install -g serve

# Start the app on port 4001
CMD ["serve", "-s", "dist", "-l", "4001"]