FROM node:bullseye-slim

# Set Node.js app folder
RUN mkdir -p /home/node/graphql-api/node_modules
WORKDIR /home/node/graphql-api

# Copy dependencies
COPY ./package*.json .
RUN chown -R node:node /home/node/graphql-api

# Install dependencies
USER node
RUN npm install

# Copy application files
COPY --chown=node:node ./src .

# build / transpile
RUN npm run build

EXPOSE 5000
CMD [ "node", "./dist/index.js" ]