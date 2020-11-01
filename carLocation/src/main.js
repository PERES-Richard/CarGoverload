var repo = require('./repositories/neo4j_repository')
var ctrl = require('./controllers/carLocation.controller')

function main() {
    repo.initNeo4jSession()
    repo.populateDatabase()
    const PORT = process.env.PORT ? process.env.PORT : "3005"
    ctrl.initAndListen(PORT)
}

main()