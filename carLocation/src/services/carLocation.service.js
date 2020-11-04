const repo = require('../repositories/neo4j_repository')

async function searchTrackedCars(nodeId, carTypeId, res) {
    const nodes = []
    const node = await repo.getNode(nodeId)
    if (node.types.contains(carTypeId)) {
        nodes.push(node)
    } else {

    }
}

module.exports = {
    searchTrackedCars
}