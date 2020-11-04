const express = require('express')
const app = express();
const cors = require('cors');
const repository = require('../repositories/neo4j_repository');
const service = require('../services/carLocation.service');

function initRoutes() {
    app.get('/car-location/ok', async (req, res) => {
        if(await repository.databaseTest()){
            res.send('ok')
        } else {
            res.send('Waiting for neo4j')
        }
    })
    app.get('/car-location/findAllNodes', async (req, res) => {
        await repository.getAllNodes(res);
    });
    app.get('/car-location/findAllCarTypes', async (req, res) => {
        await repository.getAllCarTypes(res);
    });
    app.get('/car-location/searchTrackedCars', async (req, res) => {
        console.log('Requested node:', req.query.node, ' carTypeId', req.query.carType);
        await service.searchTrackedCars(req.query.node, req.query.carTypeId, res);
    });
}

function listen(port) {
    app.listen(port, () => {
        console.log(`Server listening on port ${port}`)
    })
}

function initAndListen(port) {
    app.use(cors());
    initRoutes()
    listen(port)
}

module.exports = {
    initAndListen
}
