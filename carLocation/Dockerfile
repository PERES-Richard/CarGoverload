FROM node:14.14
WORKDIR /app
COPY ./package.json ./
RUN npm install
COPY . .
#RUN npm run build

EXPOSE 3005
CMD ["node", "src/main.js"]
