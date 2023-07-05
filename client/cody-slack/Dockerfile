# Choose the latest LTS (Long Term Support) version of node as our base image
FROM node:lts@sha256:f4698d49371c8a9fa7dd78b97fb2a532213903066e47966542b3b1d403449da4

ARG PROJECT_ROOT="."

# Create a directory for our application
WORKDIR /app

# Copy only the build artifact and package.json into our app directory
COPY ${PROJECT_ROOT}/dist .
COPY ${PROJECT_ROOT}/package .

# Install only production dependencies
RUN npm install --omit=dev --silent

# The application's default port is picked from the environment
EXPOSE $PORT

# Start the application
CMD [ "node", "app.js" ]
