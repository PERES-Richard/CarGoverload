const repo = require('../repositories/neo4j_repository')
const axios = require('axios');

const DISTANCE_MARGIN = 50;

async function newSearch(value, callback) {
    const searchParameters = JSON.parse(value)
    if (searchParameters.departureNode !== undefined &&
        searchParameters.arrivalNode !== undefined &&
        searchParameters.carType !== undefined &&
        searchParameters.searchId !== undefined) {
            searchTrackedCars(
                searchParameters.departureNode,
                searchParameters.arrivalNode,
                searchParameters.carType,
                searchParameters.searchId,
                callback
            )
    } else {
        throw 'Error in message parameters'
    }
}

async function validateSearch(value, callback) {
    await newSearch(value, callback)
}

async function searchTrackedCars(departureNode, arrivalNode, carType, searchId, callback) {
    const carTypeId = (await repo.getCarType(carType)).id;
    console.log("######## car type id : ", carTypeId)

    const nodes = []
    const node = await repo.getNode(departureNode);
    console.log("############# Node departure : ", node);
    let includes = false
    node.types.forEach(t => {
        if (t == carTypeId) {
            includes = true
        }
    })
    if (includes) {
        nodes.push(node)
    } else {
        const tmpNodes = await repo.getNodesCloserThan(node.id, DISTANCE_MARGIN)
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
        const destNode = await repo.getNode(arrivalNode);
        console.log("############# Node arrival : ", destNode);
        if (destNode.types.includes(carTypeId)) {
            cars.forEach(car => {
                trackedCars.push({node: nodes[i], destinationNode: destNode, car: car})
            })
        } else {
            const closeNodes = repo.getNodesCloserThan(destNode.id, DISTANCE_MARGIN)
            closeNodes.filter(a => a.types.includes(carTypeId)).forEach(n => {
                cars.forEach(car => {
                    trackedCars.push({node: nodes[i], destinationNode: n, car: car})
                })
            })
        }
    }

    console.log(JSON.stringify(trackedCars));
    callback("car-location-result", "{ \"searchId\": \"" + searchId + "\", \"results\":" + JSON.stringify(trackedCars) + " }").catch(err => console.log("Error: " + err));
}

async function getCloseCars(latitude, longitude, carTypeId) {
    const CAR_TRACKING_PORT = process.env.CAR_TRACKING_PORT ? process.env.CAR_TRACKING_PORT : "3005"
    const CAR_TRACKING_HOST = process.env.CAR_TRACKING_HOST ? process.env.CAR_TRACKING_HOST : "localhost"
    const result = await axios.get('http://'+CAR_TRACKING_HOST+':'+CAR_TRACKING_PORT+'/car-tracking/get-cars?latitude='+latitude+'&longitude='+longitude+'&type='+carTypeId);
    console.log("###############Fetched cars : ", result.data);
    return result.data;
}

module.exports = {
    newSearch,
    validateSearch
}
