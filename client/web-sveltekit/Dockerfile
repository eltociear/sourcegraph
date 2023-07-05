# This image creates a standalone web server for demoing the prototype
FROM node:16@sha256:466d0a05ecb1e5b9890960592311fa10c2bc6012fc27dbfdcc74abf10fc324fc

ARG PROJECT_ROOT="."

WORKDIR /usr/src/app
RUN echo '{"type": "module"}' > package.json
RUN npm install express@~4.18 http-proxy-middleware@~2.0
COPY  ${PROJECT_ROOT}/server.js ./
COPY ${PROJECT_ROOT}/build ./build

ENV SOURCEGRAPH_API_HOST=https://sourcegraph.com
EXPOSE 4173
CMD ["node", "server.js"]
