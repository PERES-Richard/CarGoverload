const repo = require('./repositories/neo4j_repository')
const ctrl = require('./controllers/carLocation.controller');

function main() {
    repo.initNeo4jSession().catch(e => {
        console.log(e)
    })
    const PORT = process.env.PORT ? process.env.PORT : "3005"
    ctrl.initAndListen(PORT)
}

main()
