# src/frontend/Dockerfile
FROM node:22.15.0

WORKDIR /app

COPY . .

RUN npm install

# If you're using Next.js
RUN npm run build

EXPOSE 4001

CMD ["npm", "start"]
