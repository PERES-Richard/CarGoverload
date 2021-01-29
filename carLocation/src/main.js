const repo = require('./repositories/neo4j_repository')
const kafka = require('./controllers/carLocation.kafka')


function main() {
    repo.initNeo4jSession().catch(e => {
        console.log(e)
    })
    const PORT = process.env.PORT ? process.env.PORT : "3005"

    kafka.initConnection().catch(e => {
        console.log(e)
    })
}

main()
