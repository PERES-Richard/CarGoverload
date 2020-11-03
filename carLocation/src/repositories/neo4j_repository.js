const neo4j = require('neo4j-driver')
let driver

function initNeo4jSession() {
   const HOST = process.env.NEO4J_HOST ? process.env.NEO4J_HOST : "localhost";
    driver = neo4j.driver(
        'neo4j://'+HOST,
        neo4j.auth.basic('neo4j', 'superpassword')
    )
}

async function addCarType(id, name) {
    const session = driver.session();
    await session.run('CREATE (c:CarType {id: $id, name: $name}) RETURN c', {
        name: name,
        id: neo4j.int(id),
    });
    await session.close();
    console.log("CarType added : " + id + " - " + name);
}

async function addNode(id, name, types, lat, lgt) {
    const session = driver.session();
    await session.run('CREATE (a:Node {id: $id, name: $name, types: $types, latitude: $lat, longitude: $lgt}) RETURN a', {
        name: name,
        id: neo4j.int(id),
        types: types,
        lat: lat,
        lgt: lgt
    });
    await session.close();
    console.log("Node added : " + id + " - " + name);
}

async function addDistance(idNode1, idNode2, value) {
    const session = driver.session();
    await session.run('MATCH (a: Node), (b: Node) WHERE a.id = $idNode1 AND b.id = $idNode2 CREATE (a)-[d: Distance {value: $value}]->(b) RETURN a, b, d', {
        idNode1: neo4j.int(idNode1),
        idNode2: neo4j.int(idNode2),
        value: neo4j.int(value)
    });
    await session.close();
    console.log("Distance added : " + idNode1 + " -> " + value + " -> " + idNode2);
}

async function getAllNodes(res) {
    const session = driver.session();
    const records = await session.run('MATCH (a: Node) RETURN DISTINCT a', {})
    await session.close();
    res.send(JSON.stringify(records.records));
}

async function deleteAllNodes() {
    const session = driver.session();
    await session.run('MATCH p=()-[r:Distance]->() DELETE p', {});
    await session.close();
    console.log("Deleting all nodes");
}

async function getAllCarTypes(res) {
    const session = driver.session();
    const nodes = [];
    await session.run('MATCH (a: CarType) RETURN a', {}).subscribe({
        onNext: record => {
            console.log('Fetched cartype: ' + JSON.stringify(record))
            nodes.push(record);
        },

        onCompleted: () => {
            res.send(JSON.stringify(nodes));
            session.close() // returns a Promise
        },
    })
}

async function deleteAllCarTypes() {
    const session = driver.session();
    await session.run('MATCH (c: CarType) DELETE c', {});
    await session.close();
    console.log("Deleting all car types");
}

async function populateDatabase() {
    await deleteAllCarTypes();
    await deleteAllNodes();


    await addCarType(0, "Solid");
    await addCarType(1, "Liquid");


    await addNode(0, 'Marseille', [0, 1], 43.9415356, 4.7632126);
    await addNode(1, 'Avignon-liquid', [1], 43.9415356, 4.7632126);
    await addNode(2, 'Nice', [0], 43.7031691, 7.1827772);
    await addNode(3, 'Paris', [0, 1], 48.8588377, 2.2770202);
    await addNode(4, 'Avignon-solid', [0], 43.9415387, 4.7632200);

    await addDistance(0,1, 85);
    await addDistance(0,2, 215);
    await addDistance(1, 2, 199);
    await addDistance(0, 3, 660);
    await addDistance(1, 4, 20);
}

module.exports =  {
    initNeo4jSession,
    populateDatabase,
    getAllNodes,
    getAllCarTypes
}
