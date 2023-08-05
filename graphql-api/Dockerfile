FROM node:bookworm

# Set Node.js app folder
RUN mkdir -p /home/node/graphql-api/
WORKDIR /home/node/graphql-api

# Copy dependencies
COPY . .
RUN chown -R node:node /home/node/graphql-api

# Install dependencies
USER node

# Incl package.json devDeps eg tsc
RUN NODE_ENV=development npm ci

# build / transpile
RUN npm run build

EXPOSE 5000

HEALTHCHECK --interval=15s CMD curl -f http://localhost:5000/health || exit 1

CMD [ "node", "./dist/index.js" ]
