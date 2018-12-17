// Test transaction

Head();
// Transact("0xe4CbBbeDF9c6a8143023953A8eE4c95e600144A8", 0.02);
// Transact(8);

//RunTransactionByOne();
//GetBalaceByAccount(1);
//GetBalaceByAccount(2);
//GetBalaceByAccount(8);
// Contractdeploy();
GetInfobloc();

// ***************************************
// Launch transaction in loop
// from - account from summ send to other accounts
// ***************************************
function GetBalaceByAccount(from){
	ac=eth.accounts[from];
    fb=eth.getBalance(eth.accounts[from]);
	console.log("Account : ", ac, "=",fb);
}	

// ***************************************
// Launch transaction in loop
// from - account from summ send to other accounts
// ***************************************
function RunTransactionByOne(from){
	
	console.log("Start transactions.")
	trz(from,9, 0.000123);
	trz(from,1, 0.0002);
	trz(from,2, 0.0003);
	trz(from,3, 0.0004);
	trz(from,4, 0.0005);
	trz(from,5, 0.0006);
	trz(from,6, 0.0007);
	trz(from,7, 0.00012);
	trz(from,9, 0.0001);
	trz(from,10,0.000123);

	GetBalaceByAccount(from);
}


// ***************************************
// Launch transaction in loop
// ***************************************
function RunTransaction(){

    a=[1,2,3,4,6,7,9];

    for (var pr in a){
    	trz(8,a[pr], 0.0005)
        console.log(a[pr]);    	
    }
}    

// ***************************************
// Get my accounts
// ***************************************
function GetAllAccounts(){

    a=eth.accounts;

    for (var pr in a){
        console.log(a[pr]);    	
    }
}    

// ***************************************
// Test
// ***************************************
function Head(){
    t="Testing main procedure\n";
    t=t+"Version 0.01";
    console.log(t);
}

// ***************************************
// Transaction
// ***************************************
function trz(in_from, in_to, vall) {

	console.log("");
	console.log("**************************");
	console.log("Send to adress ->>", in_to);

	dt = "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675";
	// dt2= web3.utf8ToHex("Hello world"); 
	// console.log(dt2);

	fb=eth.getBalance(eth.accounts[in_from]);
	tb=eth.getBalance(eth.accounts[in_to]);
	console.log("Start balance fron", fb);
	console.log("Start balance to",   tb);

	from = eth.accounts[in_from];
	to   = eth.accounts[in_to];

    // "0x9123e70"
	vallue = web3.toWei(vall, "ether");
    console.log("Value       : ", vallue);	

	// from ="0x50965E6984A63998381Fa4D8d0aE70E94A619b7b";


	params={"from":from,"to":to,"gas":"0x76c0", "gasPrice":"0x9184e720","value":vallue,"data":dt};

	personal.unlockAccount(from,"Password");
	tr=eth.sendTransaction(params);
	personal.lockAccount(from);

	console.log("Transaction : ", tr);
	console.log("FROM        : ", from);
	console.log("TO          : ", to);

	fbl=eth.getBalance(eth.accounts[in_from]);
	tbl=eth.getBalance(eth.accounts[in_to]);

	console.log("Finish balance from ......", fbl);
	console.log("Finish balance to   ......", tbl);
	console.log("---------------------------------------------");
}

// ******************************************
// Транзакция на счет 
// ******************************************
function Transact(in_from){
      
     dt   = "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675";
	 to   = "0xe4CbBbeDF9c6a8143023953A8eE4c95e600144A8";
	 
     vallue = web3.toWei(0.00222, "ether");
	 from = eth.accounts[in_from];
     personal.unlockAccount(from,"Password");

     params={"from":from,"to":to, "gas":"0x76c0", "gasPrice":"0x9184e720","value":vallue,"data":dt};

 	 tr = eth.sendTransaction(params);
	 personal.lockAccount(from);

    console.log("Transaction : ", tr);
	console.log("FROM        : ", from);
	console.log("TO          : ", to);
}

// ***************************************
// Blocks info
// ***************************************
function GetInfobloc(){
	console.log("New Account");
	bl=eth.getBlock("pending",true);
	console.log("NEW BLOCK:",bl.number);

	// txpool
	// tx=txpool.ValueOf();
	// console.log(tx);

	// New personal Account
	// pr=personal.newAccount("ss");
	// console.log("Acccount :", pr);

	// acc=eth.accounts;
	// console.log(acc);

	// Cинхронизированна нода или нет
	tinc=eth.syncing;
	console.log("Thincing :",tinc);
	if (tinc==false){
	    console.log("Node is syncing");
	}else{
	    console.log("Node is NOT syncing!");
	}

}

// ************************************************
// Deploy contract
// ************************************************
function Contractdeploy(){
   var source   = 'contract Multiply7 { event Print(uint); function multiply(uint input) returns (uint) { Print(input * 7); return input * 7; } }';
   var compiled = web3.eth.compile.solidity(source);
   var code     = compiled.Multiply7.code;
   var abi      = compiled.Multiply7.info.abiDefinition;
   var fromin   = eth.accounts[8];

   web3.eth.contract(abi).new({from: fromin, data: code}, function (err, contract) {
   if (!err && contract.address)
         console.log("deployed on:", contract.address);
      }
);
 
   multi.Print(function(err, data) { console.log(JSON.stringify(data)) })
}
