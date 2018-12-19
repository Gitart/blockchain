// (c) Essentia - 2018
// Savchenko Arthur
// 18-12-2018 

// RPC docs
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#list-of-management-apis

package main

import (	
    "fmt"
    "io/ioutil"
   	"time"
    "log"
    "flag"
    "net/http"
    "strings"
    "encoding/json"
    // "math/big"
)

// 
// Main procedure
// 
func main() {

	 // // Command-Line Flags
	 ip    := flag.String ("ip",      "127.0.0.1",   "Ip")
     port  := flag.String ("port",    "8545",        "Port")
     pass  := flag.String ("pass",    "Secret00$$$", "Password")
     from  := flag.String ("from",    "",            "From")
     to    := flag.String ("to",      "",            "To")
     count := flag.Int    ("cnt",     5,             "Count transaction")
     money := flag.Int64  ("money",   1000,          "Money for transaction")
     gas   := flag.Int64  ("gas",     50000,         "Gas")
     unlck := flag.Bool   ("unlock",  true,          "Unlock for accounts")

     flag.Parse()
     
     Money := ""
     Gas   := ""
     // value := big.NewInt(*money)
     // fmt.Println(value)

    var M   int64 = 107
    // var Gas int64 = 50000

    if *count==0{
       fmt.Println("INFO: Replace count transaction to 1 !") 
       *count=1 
    }

    // Check ip
    if *ip==""{
       fmt.Println("WARNING: Ip address requared !")
       return
    }
    
    // Check port
    if *port==""{
       fmt.Println("WARNING: port is empty !")
       return
    }

    // Check form address
    if *from=="" {
       fmt.Println("WARNING: From address is empty!")
       return
    }else{
       fmt.Println("From ..............: ", *from) 
    }

    // Check form address
    if *to=="" {
       fmt.Println("WARNING: To address is empty!")
       return
    }else{
       fmt.Println("To ................: ", *to) 
    }

    // Check money
    if *money==0 {
        fmt.Println("WARNING! Money = 0 ! Please input money.")
        Money=int642hex(M)  
        fmt.Println("Money = "+Money)
    }else{
        // Money=int642hex(*money*1000000000000000000)      
        Money=int642hex(*money)      
        fmt.Println("Money .............: ", *money, "Hex :", Money ) 
    }

    // Check gas
    if *gas==0 {
       fmt.Println("WARNING: Gass is empty !")
       return
    }else{
       Gas=int642hex(*gas)  
       fmt.Println("Gas ...............: ", *gas, "Hex : ", Gas) 
    }

    // Check password
    if *pass==""{
       fmt.Println("WARNING: Password is empty !")
       return
    }else{
       GlobalPass = *pass
       fmt.Println("Password ..........: ", *pass) 
    }

    // Create url path
    url   := "http://" + *ip + ":" + *port + "/"
    fmt.Println("Url ...............: ", url,"\n") 
    
    // Accounts
    Accs :=[]string{*from, *to}
    
    if *unlck {
        // Unlock all accounts
        for _, Accr:=range(Accs){
            UnlckAccount(url, Accr, *pass)
            fmt.Println(url, Accr, *pass)
        }
    }

    // Gas Price
    // 0x9184         - old value
    // 0x3b9aca00     - 
    GasPrice:=Get_gas_price(url)
    fmt.Println("Gas price ........ : ", GasPrice)
    fmt.Println("")


    fmt.Println("**********************************  T  R  A  N  S  A  C  T  I  O  N  S  ***************************************** ", "\n")
    // Transaction 
    for i := 0; i < *count; i++ {
        func (){
            gettx := SetTransaction(url, *from, *to, Gas, GasPrice, Money, GlobalData)
            log.Println("Tx: ", gettx)
            time.Sleep(time.Millisecond*600)
        }()
        

    
        fmt.Printf("Tx id %v  %s => %s eth[%v] \n", i,  *from, *to, *money ) 
    }


    fmt.Println("\n\n THE END!")
}    

// *****************************************************************************************************
// Перечисление Со счета А-Б
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func SetTransaction(Url, From, To, Gas, Gasprise, Value, Data string) string {
   
    
    // Transaction
    Params  := `"params":[{"from":"`+From+`","to":"`+To+`","gas":"`+Gas+`","gasPrice":"`+Gasprise+`","value": "`+Value+`","data": "`+Data+`"}]`
    payload :=  strings.NewReader(`{"id":104, "jsonrpc":"2.0", "method":"eth_sendTransaction", `+Params+`}`)
    req, _  :=  http.NewRequest("POST", Url, payload)
    req.Header.Add("content-type",  "application/json")
    req.Header.Add("cache-control", "no-cache")

     res, err := http.DefaultClient.Do(req)
     if err!=nil{
        log.Println("Error client ", err.Error())
        return ""
     }

     defer res.Body.Close()
     body, errb := ioutil.ReadAll(res.Body)
     if errb!=nil{
        log.Println("Error read body", errb.Error())
        return ""
     }
     return string(body)
}


// ************************************************************
// Unlock account
// ************************************************************
func UnlckAccount(Url, Acc, Pas string){
     defer func(){
          recover()
     }()


        unlocked       := strings.NewReader(`{"id":100, "jsonrpc":"2.0", "method": "personal_unlockAccount", "params":["` + Acc + `","` + Pas + `", 15000]}`)
        unlreq, errt   := http.NewRequest("POST", Url, unlocked)

        if errt!=nil{
           // w.Write([]byte(fmt.Sprintf("Error")))
           log.Println("Unlock: Error client : ", errt.Error())
           return
        }

        unlreq.Header.Add("content-type",  "application/json")
        unlreq.Header.Add("cache-control", "no-cache")

        resu, err_unl := http.DefaultClient.Do(unlreq)
        if err_unl!=nil{
           log.Println("Error client : ", err_unl.Error())
        }
               
        defer resu.Body.Close()

        bodys, errbs := ioutil.ReadAll(resu.Body)
        if errbs!=nil{
           log.Println("Error read body", errbs.Error())
        }

        
        ret:=Runlocked(Acc,string(bodys))
        fmt.Println("UNLOCK : ",ret, "\n")
}



// ******************************************************
// Get Gas Price 
// ******************************************************
func Get_gas_price(Url string) string {

            var Rs Mst

             // Gas Price 
            GetGasPrice    :=  strings.NewReader(`{"jsonrpc":"2.0","method":"eth_gasPrice","params":[],"id":73}`)
            GasReq, errt   :=  http.NewRequest("POST", Url, GetGasPrice)

            if errt!=nil{
               log.Println("Error gasprice get : ", errt.Error())
               return ""
            }
            GasReq.Header.Add("content-type",  "application/json")
            GasReq.Header.Add("cache-control", "no-cache")

            resGas, err_gasprice := http.DefaultClient.Do(GasReq)
            if err_gasprice!=nil{
               log.Println("Error gas price body return : ", err_gasprice.Error())
               return ""
            }
           
            defer resGas.Body.Close()
            
            // Read
            GasPrice, errgasprice := ioutil.ReadAll(resGas.Body)
            if errgasprice!=nil{
               log.Println("Error get gas price ", errgasprice.Error())
               return ""
            }
            
            // Unmarshal body
            errj:=json.Unmarshal([]byte(GasPrice), &Rs)
            
            if errj!=nil{
               log.Println("Error ummarshaling.", errj.Error())
               return "0"
            }

            Fs:=Rs["result"].(string)
            return Fs
}



// *********************************************************************
// Convert int -> Str(0x.....) -> 
// *********************************************************************
func int642hex (i int64) string {
    d:= fmt.Sprintf("0x%02x", i)
    return d
}


// *********************************************************************
// Analyze 
// *********************************************************************
func Mreturn(T string) string{
     if strings.Contains(T,"The method personal_newAccount does not exist/is not available") {
        return "Не включена опция Personal"
     }

     if strings.Contains(T,"Error client connect") {
        return "Нoда не доступна."
     }
     
     return T
}

// *********************************************************************
// Analyze 
// *********************************************************************
func Runlocked(A,T string) string{
     if strings.Contains(T,"true") {
        return "Address : " + A +" Unlocked !"
     }

    return "WARNING: Address : " + A +"  NO UNLOCKED !"
}
