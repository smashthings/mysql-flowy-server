# This is the docker image that can be used in flowy, simply copy this to the base of the repo

FROM node:14
COPY . /flowy
WORKDIR /flowy
RUN yarn install && yarn build
CMD ["yarn", "start"]