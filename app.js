const express = require('express');
const redis = require('redis');
const app = express();
const port = 8080;

app.use(express.json());

const MAX_ITEMS = 80_000;

const rc = redis.createClient();

rc.on('connect', () => {
    rc.set('items_available', MAX_ITEMS, (err, reply) => {
        if (err) {
            console.error(err);
        } else if (reply === 1) {
            console.log('Initialized "items_available" with value:', MAX_ITEMS);
        } else {
            console.log('"items_available" already exists');
        }
        rc.quit();
    });
});

rc.connect();

app.post('/order', (req, res) => {
    const order = req.body;
    // Process the order here
    
});

app.listen(port, () => {
    console.log(`Server listening at http://localhost:${port}`);
});
