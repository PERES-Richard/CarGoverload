const neo4j = require('neo4j-driver')
let driver

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

async function initNeo4jSession() {
    if (!getDriver()) {
        throw 'Could not create driver'
    }
    while (!(await databaseTest())){
        console.log('Neo4j database not ready, retrying in 10 seconds...')
        await sleep(10000)
    }
    await populateDatabase()
}

function getDriver() {
    const HOST = process.env.NEO4J_HOST ? process.env.NEO4J_HOST : "localhost";
    try {
        driver = neo4j.driver(
            'neo4j://'+HOST,
            neo4j.auth.basic('neo4j', 'superpassword')
        )
    } catch (e) {
        console.log('Database connection error: ', e)
        return false
    }
    return true
}

async function databaseTest() {
    const session = driver.session()
    try {
        console.log('Launching Neo4j database test...')
        await session.run('MATCH (e: Node) RETURN e')
        console.log('Neo4j database ok')
        return true
    } catch(e) {
        console.log('Database test error: ', e)
        return false
    }
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
    const t = []
    types.forEach(type => {
        t.push(neo4j.int(type))
    })
    await session.run('CREATE (a:Node {id: $id, name: $name, types: $types, latitude: $lat, longitude: $lgt}) RETURN a', {
        name: name,
        id: neo4j.int(id),
        types: t,
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

    const nodes = [];
    records.records.forEach(function(record){
       const recordProperties = record["_fields"][0].properties
       nodes.push({
           name: recordProperties.name,
           id: recordProperties.id.low,
           types: recordProperties.types,
           latitude: recordProperties.latitude,
           longitude: recordProperties.longitude
       })
    });

    res.send(JSON.stringify(nodes));
}

async function getNode(id) {
    const session = driver.session();
    const records = await session.run('MATCH (a: Node) WHERE a.id = $id RETURN DISTINCT a LIMIT 1', {
        id: neo4j.int(id)
    });
    await session.close();

    if (records.records[0] === undefined)
        return undefined
    const neoNode = records.records[0]["_fields"][0].properties
    neoNode.id = id
    const intTypes = []
    neoNode.types.forEach(t => {
        intTypes.push(t.low)
    })
    neoNode.types = intTypes
    return neoNode
}

async function getNodesCloserThan(nodeId, distance) {
    const session = driver.session();
    const records = await session.run('MATCH (a: Node)-[b: Distance]-(c: Node) WHERE a.id = $nodeId AND b.value < $distance RETURN c', {
        nodeId: neo4j.int(nodeId),
        distance: neo4j.int(distance)
    });
    await session.close();
    if (records.records[0] === undefined)
        return []
    const res = []
    for (let i = 0; i < records.records.length; i++) {
        const neoNode = records.records[i]["_fields"][0].properties
        neoNode.id = neoNode.id.low
        const intTypes = []
        neoNode.types.forEach(t => {
            intTypes.push(t.low)
        })
        neoNode.types = intTypes
        res.push(neoNode)
    }
    console.log(res)
    return res
}

async function getCarType(id) {
    const session = driver.session();
    const records = await session.run('MATCH (a: CarType) WHERE a.id = $id RETURN DISTINCT a LIMIT 1', {
        id: neo4j.int(id)
    });
    await session.close();
    if (records.records[0] === undefined)
        return undefined
    const neoType = records.records[0]["_fields"][0].properties
    neoType.id = id
    return neoType
}

async function deleteAllNodes() {
    const session = driver.session();
    await session.run('MATCH p=()-[r:Distance]->() DELETE p', {});
    await session.close();
    console.log("Deleting all nodes");
}

async function getAllCarTypes(res) {
    const session = driver.session();
    const records = await session.run('MATCH (a: CarType) RETURN a', {});
    await session.close();

    const carTypes = [];
    records.records.forEach(function(record){
        const recordProperties = record["_fields"][0].properties
        carTypes.push({
            name: recordProperties.name,
            id: recordProperties.id.low
        });
    });
    res.send(JSON.stringify(carTypes));
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


    await addCarType(1, "Solid");
    await addCarType(2, "Liquid");


    await addNode(1, 'Marseille', [1, 2], 43.9415356, 4.7632126);
    await addNode(2, 'Avignon-liquid', [2], 43.9415356, 4.7632126);
    await addNode(3, 'Nice', [1], 43.7031691, 7.1827772);
    await addNode(4, 'Paris', [1, 2], 48.8588377, 2.2770202);
    await addNode(5, 'Avignon-solid', [1], 43.9415387, 4.7632200);

    await addDistance(1, 2, 85);
    await addDistance(1, 3, 215);
    await addDistance(2, 3, 199);
    await addDistance(1, 4, 660);
    await addDistance(2, 5, 20);
}

module.exports =  {
    initNeo4jSession,
    databaseTest,
    getAllNodes,
    getAllCarTypes,
    getNodesCloserThan,
    getNode
}
