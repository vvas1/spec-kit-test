# Frontend dev server for local development (004-containerize-frontend)
# Dev server listens on 5137 inside container; host maps 3000:5137
FROM node:20-alpine

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

# Vite default is 5173; spec requires 5137 for containerized run
ENV VITE_API_URL=http://backend:8080
EXPOSE 5137

CMD ["npm", "run", "dev", "--", "--port", "5137", "--host", "0.0.0.0"]
