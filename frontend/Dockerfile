# Build Stage
FROM node:18-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

COPY . .
RUN npm run build

# Production Stage
FROM nginx:alpine

# Copy built assets from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Optional: Add custom nginx.conf if needed for SPA fallback
# COPY nginx.conf /etc/nginx/conf.d/default.conf

# Default Nginx config usually works for static files, 
# but for SPA routing (React Router) usually need try_files $uri /index.html.
# We'll create a simple config inline or assume default for now.
# To support standard SPA routing without 404s on refresh:

RUN echo 'server { \
    listen 80; \
    location / { \
    root /usr/share/nginx/html; \
    index index.html index.htm; \
    try_files $uri $uri/ /index.html; \
    } \
    }' > /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
