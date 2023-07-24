FROM node:16.14.2-alpine3.15

WORKDIR /app

COPY package.json .

COPY . .

# For Chiado Testnet
 COPY .env.chiado .env

# Installs all node packages
RUN npm install

# Build for production.
RUN npm run build --production

# Install `serve` to run the application.
RUN npm install -g serve

# Uses port which is used by the actual application
EXPOSE 3000

# Run application

# Use for development server
# CMD [ "npm", "start" ]

# Use for production server
CMD serve -s build