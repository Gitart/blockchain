
console.log("Тестовый запрос");
myacc   = eth.accounts[8];
balance = web3.fromWei(eth.getBalance(myacc), "ether");
mybal   = eth.getBalance(myacc);

console.log("Account",    myacc);
console.log("Balance",    balance);
console.log("My Balance", mybal);


// New account
blk = eth.getBlock("latest");
console.log("Blok", blk.number);
console.log("---------------------------------");
//console.log("Blok all transactions", blk.transactions);
console.log("Blok first transact",   blk.transactions[0]);
console.log("Blok timestamp",        blk.timestamp);
console.log("Blok size", blk.size);


blkn=eth.getBlock(blk.number);

console.log("Blok hash", blkn.hash);





