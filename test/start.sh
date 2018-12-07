
echo Start Node test
./geth --rinkeby --rpc --rpcapi db,eth,net,web3,personal --cache=2048 --datadir db --verbosity 0  -syncmode light --cache=1024 --rpcport 8545 --rpcaddr 127.0.0.1  console

## ./geth --rinkeby --syncmode "fast" --rpc --rpcapi db,eth,net,web3,personal --cache=1024 --rpcport 8545 --rpcaddr 127.0.0.1 --rpccorsdomain "*"
## ./geth --testnet --cache=2048 --datadir db --verbosity 0 console
## ./geth --rinkeby --rpc --rpcapi db,eth,net,web3,personal --cache=2048 --datadir db --verbosity 0  --mine --cache=1024 --rpcport 8545 --rpcaddr 127.0.0.1  console

