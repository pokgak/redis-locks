# TODO

- limit to only `MAX_ITEMS` successful orders
- after that, return 200 to users with message "No more items available"
- send all orders to a queue/kafka for postprocessing
- setup frontend that shows live counter of orders available


## Testing

TODO: setup k6s testing script that will send POST request with a unique client-id
