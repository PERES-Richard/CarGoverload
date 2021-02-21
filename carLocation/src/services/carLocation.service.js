const repo = require('../repositories/neo4j_repository')
const axios = require('axios');

const DISTANCE_MARGIN = 250;

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

    console.log("######## searchId : ", searchId)

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
                trackedCars.push({
                    node: nodes[i],
                    destinationNode: destNode,
                    car: car,
                    distance: getDistanceFromLatLonInKm(node.latitude, node.longitude, destNode.latitude, destNode.longitude)
                })
            })
        } else {
            const closeNodes = await repo.getNodesCloserThan(destNode.id, DISTANCE_MARGIN)
            console.log("########## Close nodes fetched : ", closeNodes)
            closeNodes.filter(a => a.types.includes(carTypeId)).forEach(n => {
                cars.forEach(car => {
                    trackedCars.push({
                        node: nodes[i],
                        destinationNode: n,
                        car: car,
                        distance: getDistanceFromLatLonInKm(node.latitude, node.longitude, n.latitude, n.longitude)
                    })
                })
            })
        }
    }

    console.log("{ \"searchId\": \"" + searchId + "\", \"results\":" + JSON.stringify(trackedCars) + " }");
    callback("car-location-result", "{ \"searchId\": \"" + searchId + "\", \"results\":" + JSON.stringify(trackedCars) + " }").catch(err => console.log("Error: " + err));
}

function getDistanceFromLatLonInKm(lat1, lon1, lat2, lon2) {
    const R = 6371; // Radius of the earth in km
    const dLat = deg2rad(lat2 - lat1);  // deg2rad below
    const dLon = deg2rad(lon2-lon1);
    const a =
        Math.sin(dLat/2) * Math.sin(dLat/2) +
        Math.cos(deg2rad(lat1)) * Math.cos(deg2rad(lat2)) *
        Math.sin(dLon/2) * Math.sin(dLon/2)
    ;
    const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a));
     // Distance in km
    return R * c;
}

function deg2rad(deg) {
    return deg * (Math.PI/180)
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
