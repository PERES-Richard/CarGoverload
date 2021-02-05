const repo = require('../repositories/neo4j_repository')
const kafka = require('../controllers/carLocation.kafka')
const axios = require('axios');

async function newSearch(value) {
    const searchParameters = JSON.parse(value)
    if (searchParameters.nodeId !== undefined &&
        searchParameters.carTypeId !== undefined &&
        searchParameters.distance !== undefined &&
        searchParameters.searchId !== undefined) {
            searchTrackedCars(
                searchParameters.nodeId,
                searchParameters.carTypeId,
                searchParameters.distance,
                searchParameters.searchId
            )
    } else {
        throw 'Error in message parameters'
    }
}

async function validateSearch(value) {
    newSearch(value)
}

async function searchTrackedCars(nodeId, carTypeId, distance, searchId, destinationNodeId) {
    const nodes = []
    const node = await repo.getNode(nodeId)
    let includes = false
    node.types.forEach(t => {
        if (t == carTypeId) {
            includes = true
        }
    })
    if (includes) {
        nodes.push(node)
    } else {
        const tmpNodes = await repo.getNodesCloserThan(nodeId, distance)
        tmpNodes.forEach(node => {
            let inc = false
            node.types.forEach(t => {
                if (t == carTypeId) {
                    inc = true
                }
            })
            if (inc) {
                nodes.push(node)
            }
        })
    }

    const trackedCars = []
    for (let i = 0; i < nodes.length; i++) {
        const cars = await getCloseCars(node.latitude, node.longitude, carTypeId)
        const destNode = await repo.getNode(destinationNodeId)
        if (destNode.types.includes(carTypeId)) {
            cars.forEach(car => {
                trackedCars.push({node: nodes[i], destinationNode: destNode, car: car})
            })
        } else {
            const closeNodes = repo.getNodesCloserThan(destinationNodeId, distance)
            closeNodes.filter(a => a.types.includes(carTypeId)).forEach(n => {
                cars.forEach(car => {
                    trackedCars.push({node: nodes[i], destinationNode: n, car: car})
                })
            })
        }
    }

    kafka.sendMessage("car-location-result", "{ searchId: " + searchId + ", results:" + JSON.stringify(trackedCars) + " }")
}

async function getCloseCars(latitude, longitude, carTypeId) {
    const CAR_TRACKING_PORT = process.env.CAR_TRACKING_PORT ? process.env.CAR_TRACKING_PORT : "3005"
    const CAR_TRACKING_HOST = process.env.CAR_TRACKING_HOST ? process.env.CAR_TRACKING_HOST : "localhost"
    const result = await axios.get('http://'+CAR_TRACKING_HOST+':'+CAR_TRACKING_PORT+'/car-tracking/get-cars?latitude='+latitude+'&longitude='+longitude+'&type='+carTypeId)
        .catch(e => {
            console.log('Error: car tracking service unreachable: ', e)
        })
    return result.data
}

module.exports = {
    newSearch,
    validateSearch
}