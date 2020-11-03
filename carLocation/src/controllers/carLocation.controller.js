const express = require('express')
const app = express();
const cors = require('cors');
const repository = require('../repositories/neo4j_repository');

function initRoutes() {
    app.get('/car-location/ok', (req, res) => {
        res.send('ok')
    })
    app.get('/car-location/findAllNodes', async (req, res) => {
        await repository.getAllNodes(res);
    });
    app.get('/car-location/findAllCarTypes', async (req, res) => {
        await repository.getAllCarTypes(res);
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
