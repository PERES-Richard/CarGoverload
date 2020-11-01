const express = require('express')
const app = express()

function initRoutes() {
    app.get('/car-location/ok', (req, res) => {
        res.send('ok')
    })
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