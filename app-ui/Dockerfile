FROM node:bookworm

# Set Node.js app folder
RUN mkdir -p /home/node/app-ui/
WORKDIR /home/node/app-ui

# Copy dependencies
COPY . .

RUN chown -R node:node /home/node/app-ui
RUN chmod + x get-token.sh

# Install dependencies
USER node

# Incl package.json devDeps eg tsc
RUN NODE_ENV=development npm ci


EXPOSE 3000

# build / transpile and host for docker env.
# TODO: proper build.
CMD [ "./node_modules/.bin/vite", "--host", "--port", "3000"]
