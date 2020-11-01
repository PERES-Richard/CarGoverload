var neo4j = require('neo4j-driver')
var driver

function initNeo4jSession() {
   const HOST = process.env.NEO4J_HOST ? process.env.NEO4J_HOST : "localhost";
    driver = neo4j.driver(
        'neo4j://'+HOST,
        neo4j.auth.basic('neo4j', 'superpassword')
    )
}

function addNode(id, name, types, lat, lgt) {
    var session = driver.session()
    session.run('CREATE (a:Node {id: $id, name: $name, types: $types, latitude: $lat, longitude: $lgt}) RETURN a', {
        name: name,
        id: neo4j.int(id),
        types: types,
        lat: lat,
        lgt: lgt
    }).subscribe({
        onNext: record => {
            console.log('Added node: ' + JSON.stringify(record))
        },

        onCompleted: () => {
            session.close() // returns a Promise
        },
    })
}

async function deleteAllNodes() {
    var session = driver.session()
    await session.run('MATCH (a:Node) DELETE a', {})
    await session.close()
}

async function populateDatabase() {
    await deleteAllNodes()
    addNode(0, 'Marseille', ['liquid','solid'], 43.9415356, 4.7632126)
    addNode(1, 'Avignon', ['liquid'], 43.9415356, 4.7632126)
    addNode(2, 'Nice', ['solid'], 43.7031691, 7.1827772)
}

module.exports =  {
    initNeo4jSession,
    populateDatabase
}