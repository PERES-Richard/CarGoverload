const express = require('express')
const app = express()
const repository = require('../repositories/neo4j_repository');

function initRoutes() {
    app.get('/car-location/ok', (req, res) => {
        res.send('ok')
    })
    app.get('/car-location/findAllNodes', async (req, res) => {
        await repository.getAllNodes(res);
    });
}

function listen(port) {
    app.listen(port, () => {
        console.log(`Server listening on port ${port}`)
    })
}

function initAndListen(port) {
    initRoutes()
    listen(port)
}

module.exports = {
    initAndListen
}
