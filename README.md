# Overload

## Technology

- [Go](https://golang.org/) for most most services
- [Node.js](https://nodejs.org/) with [ExpressJS](https://expressjs.com/)
- [PostgreSQL](https://www.postgresql.org/) database the secure model (store bookings)
- [Neo4j](https://neo4j.com/) database to represent each node and powerful distance calculation

## How to use

You can launch `./prepare.sh` to build the docker containers and check if each service is Ok

Get to `frontend/index.html` to use the UI

You can check the neo4j database nodes by going to `http://localhost:7474/browser/`
- Username : neo4j
- Password : superpassword
