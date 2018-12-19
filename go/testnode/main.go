// 
// Copyright 2018
// Test module for hight load 
// Essentia
// Savchenko Arthur
// 14-02-2018 

// Materials
// https://github.com/ovcharovvladimir/essentiaHybrid/tree/master/waletproxy

// RPC docs
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#list-of-management-apis
package main

import (	
    "fmt"
    "io/ioutil"
   	"time"
    "log"
    "net/http"
    "encoding/json"
    "net/http/httputil"
    "net/url"
    "strconv"
    "bytes"
    "html/template"
    "path"
    "strings"
    "context"
    "io"
    "os"
    // "math"
    "github.com/fatih/color"
    "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
)


// Address from -> To
const (

    MPort      = "7818"

	// ACC     = "0x97582e614f2f8e38f72a9d8d757a18057e68d337"
	// ACC     = "0x5c81a844fa4dbe3f3cba50f87821cb572230da7e"        // From
	// ACC2    = "0xbec03d897214d1016b860b0b3ebbd4290ec8deda"        // To
    // ACC     = "0x97d794d474e15715f607400a5bd325447977d0e1"        // Money
    // ACC     = "0x5098eac4596193c949f89ffc277bc20d01f6342b"        // from

 // ACC        = "0xb7087caab8ebe5dc230fc7d04e02ca5f7a038f72"        // 123
    ACC        = "0x7cfe90b21d7e8bb46cf52e367ade8674891e0658"              
    ACC2       = "0x38f50f4aa11b003893c93a4f70c088611502fc82"         


 // ACC2       = "0x349de42e58a10dc5cea7a8020e6d2658db7a0118"
 // ACC2       = "0x403b52bd6a52830449865ee77df673a4ba5ab4af"         // to  
 // ACC2       = "0x1b4540f309ca52b317dbc963c226ffa469c29280"

 // Money  
 // ACC2       = "0x9108f1bd24823337b460603c93b171e61f66a18f"             
    ACC4       = "0x11ac8b9dc5f7843e1b1f91d74611ab169448e4c4"

    Nodelocal  = "127.0.0.1" 
    Gess01     = "18.188.111.198"                                    // Gess 01
    Node00     = "18.222.125.100"                                    // Node 00
    Node02     = "188.190.240.195"                                   // Gess 02
    Node03     = "18.217.164.134"                                    // Gess 03
    Node04     = "18.224.11.186"                                     // Gess 04
    Node05     = "18.224.106.72"                                     // Gess 05

    Miner02    = "18.224.159.84"                                     // Amazon miner 02
    Miner03    = "18.224.168.178"                                    // Amazon miner 03
    Miner04    = "18.224.198.158"                                    // Amazon miner 04

    Acc_1      = "0x"
    Acc_2      = "0x"
    Acc_3      = "0x"
    Acc_4      = "0x"
    Acc_5      = "0x"
    Acc_6      = "0x"
    Acc_7      = "0x"
    Acc_8      = "0x"
    Acc_9      = "0x"
    Acc_10     = "0x"

    SecPass    = "JxJkJo0klM0012S"
    Portnode   = "8545"
 
    // GlobalPass = "Gerda1000"
    GlobalPass = "6apa6aHs"
    GlobalData = "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"
)

//**********************************************************
// Main
//**********************************************************
func main() {

    // Set color
    Whites  := color.New(color.FgWhite).PrintlnFunc()
    FgMag   := color.New(color.FgHiYellow)
    Cyan    := color.New(color.Bold, color.FgCyan)
    FgRed   := color.New(color.Bold, color.FgRed).PrintlnFunc()
    FgGreen := color.New(color.Bold, color.FgGreen).PrintlnFunc()

    FgGreen("Proxy server")
    Whites("Version: 1.01 (Testing)")

    // Active node
    Host  := ActiveNode()
    Sett  := ReadSettingFile()
    Port  := ":"+ Sett.Mainport

    FgMag.Print("Active node : ")
    FgRed(Host)
    Cyan.Printf("Listen proxy port %s \n", Port)

    // Check active node
    if Host== "" {
       Host="Error: All nodes disabled."
    }

    // For test with other resources
    // Host="localhost:5555"
    // Host="www.youtube.com"
    proxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme:"http", Host:Host})
    proxy.Transport = &transport{http.DefaultTransport}

    // Proxy
    proxy.Director = func(req *http.Request) {
    
    // Allows
    req.Header.Set("Content-Type","application/json;charset=utf-8")
    req.Header.Set("Access-Control-Allow-Origin","*")
    req.Header.Set("Access-Control-Allow-Headers","X-Requested-With")
    req.Header.Set("X-Forwarded-For",Host)

    req.Host       = Host
    req.URL.Host   = Host
    req.URL.Scheme = "http"

    }

    // Proxy
    http.Handle("/",                  proxy)                  // Proxy

    // Routing
    http.HandleFunc("/nodes/",        ShowNodes)              // Show all nodes
    http.HandleFunc("/active/",       Active_node)            // Show active first node for work
    http.HandleFunc("/actives/",      Active_nodes)           // Show active nodes
    http.HandleFunc("/down/",         Down_nodes)             // Show down nodes
    http.HandleFunc("/test/",         Api_test)               // Test service response
    
    // Admin panel
    http.HandleFunc("/admin/",        Api_admin)              // Admin panel
    http.HandleFunc("/admin/test/",   Api_admin_test)         // Admin test panel


    http.HandleFunc("/report/node/",  Nodes_report)           // Node report

    http.HandleFunc("/workid/",       NetworkIdTest)          // Get workid
    http.HandleFunc("/wrkid/",        NetworkIds)             // Get workid
    http.HandleFunc("/balances/",     PreviewBalanace)        // Просмoтр всех счетов
    http.HandleFunc("/balance/",      GetMyBalance)           // Получение одного баланса

    // Transactions
    http.HandleFunc("/transact/",     SetTransact)            // Транзакция
    

    // Test ноде
    http.HandleFunc("/htransact/",    SetTransact_hight)      // Нагруженный тест!!!!!!!!
    http.HandleFunc("/htransact2/",   SetTransact_hight)      // Нагруженный тест
    http.HandleFunc("/htransact3/",   SetTransact_hight)      // Нагруженный тест
    http.HandleFunc("/htransact4/",   SetTransact_hight)      // Нагруженный тест

    http.HandleFunc("/infinity1/",    SetTransact_infinity)   // Нагруженный тест
    http.HandleFunc("/infinity2/",    SetTransact_infinity)   // Нагруженный тест
    http.HandleFunc("/testpanel/",    Testpanel)              // Управление тестом

    // Accounts 
    // Работа со счетами
    http.HandleFunc("/account/",      SetAccount)             // New Account
    http.HandleFunc("/ac/new/",       NewAccount_wrk)         // Adding new aacount to tets net
    http.HandleFunc("/ac/one/",       NewAccount_one)         // Adding new aacount one
    http.HandleFunc("/ac/newopen/",   NewAccount_Open)        // Adding new aacount one from form (Panel)

    http.HandleFunc("/tx/pool/",      TxtPool)                // Adding new aacount to tets net
    http.HandleFunc("/tx/pl/",        ViewT3es)               // Adding new aacount to tets net

    // Database
    http.HandleFunc("/db/start/",     Db_start)               // Adding new aacount to tets net

    // Test
    http.HandleFunc("/test/account/", D_scan_accounts_test)   // Adding new aacount to tets net
    http.HandleFunc("/test/rend/",    Randoms)                // Test random function

    // Block
    http.HandleFunc("/test/block/",   Block_info)                // Test random function

    // Two node
    // http.HandleFunc("/nodes/",        Create_account)      // Create Account
    err := http.ListenAndServe(Port, nil)

    // Error
    if err != nil {
       log.Println("Error start service.",err.Error())
    }
}



//************************************************************
//  Transport
//************************************************************
func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {

  resp, err = t.RoundTripper.RoundTrip(req)
  if err != nil {
     return nil, err
  }

  b, err := ioutil.ReadAll(resp.Body)
  if err != nil {
     return nil, err
  }

  err = resp.Body.Close()
  if err != nil {
     return nil, err
  }

  b                  = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)
  body              := ioutil.NopCloser(bytes.NewReader(b))
  resp.Body          = body
  resp.ContentLength = int64(len(b))
  resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
  resp.Header.Set("Content-Type", "application/json; charset=utf-8")
  return resp, nil
}

//************************************************************
// Show all nodes
//************************************************************
func ShowNodes(w http.ResponseWriter, req *http.Request) {

    Nd:=ReadSettingFile()

    // Looop for nodes
    for _,nd:=range Nd.Nodes{

         if ChekNodeWork(nd.Ip, nd.Port){
            nd.Status="Active"
         }else{
            nd.Status="Not available"
         }
         WPr(nd,w)
    }
}

//************************************************************
// Show first active node
//************************************************************
func Active_node(w http.ResponseWriter, req *http.Request) {
     An:=""

     if ActiveNode()!="" {
        An = ActiveNode()
     }else{
        An = "No active node."
     }

     ft:="Active note :" + An
     Wprn(ft,w)
}

//************************************************************
// Show all active nodes
//************************************************************
func Active_nodes(w http.ResponseWriter, req *http.Request) {
    Nd := ReadSettingFile()

    // Looop for nodes
    for _,nd:=range Nd.Nodes{
         if ChekNodeWork(nd.Ip, nd.Port){
            nd.Status="Active"
            WPr(nd, w)
         }
    }
}

//************************************************************
// Text  output rep
//************************************************************
func Wprn(Txt string, w http.ResponseWriter){
     w.Write([]byte(Txt))
}

//************************************************************
// Format  output rep
//************************************************************
func WPr(nd Node, w http.ResponseWriter){
     ft:=fmt.Sprintf("Node :%-15s Ip : %-15s Port : %-6s Status : %-20s \n<br>",nd.Note,nd.Ip,nd.Port,nd.Status)
     w.Write([]byte(ft))
}

//************************************************************
// Show all disabled node
//************************************************************
func Down_nodes(w http.ResponseWriter, req *http.Request) {
    Nd:=ReadSettingFile()
    // Looop for nodes
    for _,nd:=range Nd.Nodes{

         if !ChekNodeWork(nd.Ip, nd.Port){
             nd.Status ="Is Down"
             WPr(nd,w)
          }
    }
}

//************************************************************
// Test
//************************************************************
func Api_test (wr http.ResponseWriter, rq *http.Request) {
     t:="Test service " + time.Now().Format("02/01/2006 15:45")
     fmt.Println(t)
     Wprn(t,wr)
}

//************************************************************
// Set the proxied request's host to the destination host (instead of the
// source host).  e.g. http://foo.com proxying to http://bar.com will ensure
// that the proxied requests appear to be coming from http://bar.com
//
// For both this function and queryCombiner (below), we'll be wrapping a
// Handler with our own HandlerFunc so that we can do some intermediate work
//************************************************************
func sameHost(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
           r.Host = r.URL.Host
           handler.ServeHTTP(w, r)
    })
}

//************************************************************
// Allow cross origin resource sharing
//************************************************************
func addCORS(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json;charset=utf-8")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
        handler.ServeHTTP(w, r)
    })
}

//************************************************************
// Get info about first active note
//************************************************************
func ActiveNode() string {
    Rip := ""
    Nd  := ReadSettingFile()

    // Looop for nodes
    for _,nd:=range Nd.Nodes{
        if ChekNodeWork(nd.Ip, nd.Port){
           Rip=nd.Ip+":"+ nd.Port
           break
          }
    }

    return Rip
}

//************************************************************
// Reading  setting from file json
//************************************************************
func ReadSettingFile() Sett {
    // Setting structure
    var m Sett

    // Open file with settings
    file, err := ioutil.ReadFile("./setting.json")
    if err != nil {
       log.Println("Error reading setting file.")
       return Sett{}
    }

    errj:=json.Unmarshal([]byte(file), &m)
    if errj!=nil{
       log.Println("Error ummarshaling.")
       return Sett{}
    }
    return m
}

//************************************************************
// Check node
//************************************************************
func ChekNodeWork(Ip, Port string) bool {
     timeout   := time.Duration(500 * time.Millisecond )
     client    := http.Client{Timeout: timeout}
     resp, err := client.Get("http://"+Ip+":"+Port)
     // fmt.Println("Resp:", resp, "\n")

    if err!=nil {
       return false
    }

    if resp.StatusCode==200{
       return true
    }else{
       return false
    }
}




// *********************************************************************
// Open new account
// *********************************************************************
func NewAccount_Open(w http.ResponseWriter, req *http.Request){
	Ip := req.URL.Path[len("/acc/newopen"):]

	// T  := AddAccount(Ip, Portnode)
	T  := AddAccount(Ip, MPort)
    fmt.Println("Account Node : ", T)
    r:= Mreturn(T)
    fmt.Fprintf(w, r)  
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
// Admin panel
// *********************************************************************
func Api_admin_test(w http.ResponseWriter, req *http.Request) {

    fp := path.Join("tmp", "in.html")
    tmpl, err := template.ParseFiles(fp, "tmp/main.html")
    Err(err, "Error template execute.")

    errf := tmpl.Execute(w, nil)
    Err(errf, "Error template execute.")
}


// *********************************************************************
// Admin panel
// *********************************************************************
func c(w http.ResponseWriter, req *http.Request) {

    fp := path.Join("tmp", "in.html")
    tmpl, err := template.ParseFiles(fp, "tmp/main.html")
    Err(err, "Error template execute.")

    errf := tmpl.Execute(w, nil)
    Err(errf, "Error template execute.")
}



// *********************************************************************
// Node journal
// /report/node/
// *********************************************************************
func Nodes_report(w http.ResponseWriter, req *http.Request) {
     Nd:=ReadSettingFile()

    // Looop for nodes
    for i, nd:=range Nd.Nodes{
          if ChekNodeWork(nd.Ip, nd.Port){
             Nd.Nodes[i].Status   = "Available"
             Nd.Nodes[i].Disabled = "success"
          }else{
             Nd.Nodes[i].Status   = "Not available"
             Nd.Nodes[i].Disabled = "danger"
          }

          Nd.Nodes[i].Network     = NetworkId(nd.Ip,nd.Port)
          Nd.Nodes[i].Datetime    = time.Now().Format("02.01.2006 15:04:05")
    }

    Dt        := Mst{"Dts": Nd.Nodes, "Title": "Active Nodes ", "Descript": "Serach" }
    fp        := path.Join("tmp", "journal.html")
    tmpl, err := template.ParseFiles(fp, "tmp/main.html")
    Err(err, "Error template execute.")

    errf := tmpl.Execute(w, Dt)
    Err(errf, "Error template execute.")
}



// ********************************************************
// Test
// ********************************************************
func NetworkIdTest(w http.ResponseWriter, req *http.Request) {
     ip:= "localhost"
     ip = "18.217.164.134"

     id := NetworkId(ip,"8545")
     fmt.Println(id)
}


//***************************************************************
// Network Id
//***************************************************************
func NetworkId(Ip,Port string) string {
    var Md Mst
    st  :=NetWorkInfo(Ip,Port)
    errj:=json.Unmarshal([]byte(st), &Md)
    if errj!=nil{
       log.Println("Error ummarshaling.", errj.Error())
       return "0"
    }

    Fs:=FltoStr(Md["id"].(float64))
    return Fs
}

//***************************************************************************
// net_version
// Returns the current network id.
// Returns
// String - The current network id.
// "1": Ethereum Mainnet
// "2": Morden Testnet (deprecated)
// "3": Ropsten Testnet
// "4": Rinkeby Testnet
// "42": Kovan Testnet
// ***************************************************************************
func NetWorkInfo(Ip, Port string) string{

    
     // timeout := time.Duration(500 * time.Millisecond )
     // client  := http.Client{Timeout: timeout}
     // resp, err := client.Get("http://"+Ip+":"+Port)


     url     := "http://" + Ip + ":" + Port + "/"
     payload := strings.NewReader(`{"jsonrpc":"2.0","method":"net_version", "id":1}`)
     req, _  := http.NewRequest("POST", url, payload)

     req.Header.Add("content-type",  "application/json")
     req.Header.Add("cache-control", "no-cache")
     req.Header.Add("token",         "26f916de-812a-f02d-9afa-2b84929097f2")

     res, err := http.DefaultClient.Do(req)
     if err!=nil{
        log.Println("Error client", err.Error())
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

/***************************************************************
  Check Eror
 ****************************************************************/
func Err(Er error, Txt string) {
    if Er != nil {
       log.Println("ERROR : " + Txt)
       return
    }
}

// ***************************************************************
// Preview all accounts balance
// Просомтр остатков на счетах
// 0x97582e614f2f8e38f72a9d8d757a18057e68d337 - yes
// ***************************************************************
func PreviewBalanace(w http.ResponseWriter, req *http.Request) {
    
    // IP:PORT
    Ip   := "18.188.111.198"
    Port := "8545"

    // Accounts list
    Accounts:=[]string{"0x92f81d2cddb753c50822e5d4db41e5c59ddccb82",
                       "0x19ee18e1bb682cd91ed80e47d34c5e38ee11d642",
                       "0x816d3cf337ade88b1d7f6fcae924a5f71a58fe7c",
                       "0x2ea5d7df797eb59abb125af82c7b4a7aa78c07ba",
                       "0xc70a9086b65e92514c423346c3bc018d67c41fb3",
                       "0x7637c25156ed10996a27bb09b482129f3a37ffa9",
                       "0x1f544dd0100e005b9082e0e3b0ecf5ac65debf1b",
                       "0x90df3764818b8c4da1ae9759700fe5a2bd7428cf",
                       "0xe184411f45b84fc1c28d7c41c4c27ef5bcc9ecde",
                       "0x14e0dc2983b7964ecbcd4dd45c6f03bc471f8b3c",
                       "0xa5614fa5d7808289d58c31cb6c8bca61cf7426c9",
                       "0xbc42851e21f66b648bd19d26a44a76bf11f0f7a2",
                       "0x74061b5254318661361526c7fce3cb62197e2ec8",
                       "0x97582e614f2f8e38f72a9d8d757a18057e68d337",
                       "0xdfd221d1d1e3dfe9e563ae312ac247900ae7dfd7",
                       "0x964d55c138f639e8617cfd788b2ea13a90213ff0",
                       "0x90ec0b5f932bf279de9bba675ae7779dde4ff36c",
                       "0x917af71f929d9bd5a92c60c526d4dd618c7db7e8",
                       "0x767199624c5c318043c3bb21ef5015451c57d85c"}

    Accounts=[]string{ACC,ACC2}

    // Preview Accounts
    for _, Acc:=range Accounts {
        bal:=BalanceInfo(Ip,Port, Acc)
        fmt.Println("Счета", Acc, "Balance:", bal)
    }
}

// ***************************************************************
// Preview balance
// Get One My Balance
// ***************************************************************
func GetMyBalance(w http.ResponseWriter, req *http.Request) {
	    
        Ip      := "18.188.111.198"
        Ip       = Node04
        Port    := "8545"
        FgRed   := color.New(color.Bold, color.FgRed).PrintlnFunc()

        fmt.Println("Previe Balance")
        FgRed(ACC)
 
	    bal  := BalanceInfo(Ip, Port, ACC)
        fmt.Println("Счет", ACC, "Balance:", bal)  // 0xde0b6b3a7640000=1000000000000000000

        FgRed(ACC2)
        Ip   = Node05
        bal  = BalanceInfo(Ip, Port, ACC2)
        fmt.Println("Счет", ACC2, "Balance:", bal)

        ACCN:="0x97d794d474e15715f607400a5bd325447977d0e1"

        FgRed(ACC2)
        Ip   = Node05
        bal  = BalanceInfo(Ip, Port, ACCN)
        fmt.Println("Счет", ACCN, "Balance:", bal)


}

// ***************************************************************
// Get Balance by Account
// ***************************************************************
func BalanceInfo(Ip, Port, Accounts string) string{

    // timeout := time.Duration(500 * time.Millisecond )
    // client  := http.Client{Timeout: timeout}
    // resp, err := client.Get("http://"+Ip+":"+Port)

    url     := "http://"+Ip+":"+Port+"/"
    payload := strings.NewReader(`{"id":4, "jsonrpc":"2.0", "method":"eth_getBalance","params":["`+ Accounts +`", "latest"]}`)
    req, _  := http.NewRequest("POST", url, payload)

    req.Header.Add("content-type",  "application/json")
    req.Header.Add("cache-control", "no-cache")

    res, err := http.DefaultClient.Do(req)
    if err!=nil{
       log.Println("Error client", err.Error())
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


/****************************************************************
params: [{
  "from"    :"0xb60e8dd61c5d32be8058bb8eb970870f07233155",
  "to"      :"0xd46e8dd67c5d32be8058bb8eb970870f07244567",
  "gas"     :"0x76c0", // 30400
  "gasPrice":"0x9184e72a000", // 10000000000000
  "value"   :"0x184e72a", // 2441406250   1388-5000
  "data"    :"0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"
}]
****************************************************************/
func SetUnlock(Ip, Port, Acc string) string{
     
    url     := "http://"+Ip+":"+Port+"/"
    payload :=  strings.NewReader(`{"jsonrpc":"2.0","method":"personal_unlockAccount","params":["` + Acc + `", "Gerda1000", 3600],"id":67}`)
    req, _  :=  http.NewRequest("POST", url, payload)

    req.Header.Add("content-type",  "application/json")
    req.Header.Add("cache-control", "no-cache")

    res, err := http.DefaultClient.Do(req)
    if err!=nil{
       log.Println("Error client", err.Error())
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

// ***************************************************************
// Preview all accounts balance
// Просомтр остатков на счетах
// 0x97582e614f2f8e38f72a9d8d757a18057e68d337 - yes
// ***************************************************************
func SetTransact2(w http.ResponseWriter, req *http.Request) {
    
    // IP:PORT
    Ip   := "18.188.111.198"
    Port := "8545"
    From := "0x97582e614f2f8e38f72a9d8d757a18057e68d337"      // Есть деньги
    To   := "0x5c81a844fa4dbe3f3cba50f87821cb572230da7e"      //-My
    T    := SetTransaction(Ip,Port,From,To,"0x76c0", "0x9184e72a", "0x1388", GlobalData)
    fmt.Println("Transaction return : ", T)
}

// ***************************************************************
// Preview all accounts balance
// Просомтр остатков на счетах
// 0x97582e614f2f8e38f72a9d8d757a18057e68d337 - yes
// ***************************************************************
func SetTransact(w http.ResponseWriter, req *http.Request) {
    // IP:PORT
    Ip   := "18.188.111.198"
    Port := "8545"
    Port  = "7818"
    From := ACC      // Есть деньги
    To   := ACC2     //-My

    // U:= SetUnlock(Ip,Port, ACC)
    // fmt.Println("Transaction return : ", U)

    T  := SetTransaction(Ip,Port,From,To,"0x76c0","0x9184e72a","0x1388", GlobalData)
    fmt.Println("Transaction return : ", T)
}


// ************************************************************************ 
// Нагрузочное тестирование
// Usage:  // http://localhost:8545/htransact/1111
// ************************************************************************
func SetTransact_hight(w http.ResponseWriter, req *http.Request) {
     defer func(){
    	  recover()
     }()
    
    // Количество транзакций
    cntr := req.URL.Path[len("/htransact/"):]

    // Accounts
    Accs :=[]string{"0xae97486174ce6abe73813e6eadfaf35d7f7fd09c","0xa3eff9486c10f5bf21b036ae2fb7a562e955194d"}
                 // "0x826e71a4ccb9a4fced80a3f91acf4985dae35402","0xe8342cb20cccd83709d46668b5270f9c88aefe38"}

    // Settings for testing
    Ip   := Gess01        // Ip node  
    Port := "8545"        // Port by default 
    Port  = MPort         // Port resset
    // From := ACC        // Есть деньги
    // To   := ACC2       // My
    Ctr  := 10            // Count transaction 
    // Rnd  := 0

    // Count transaction
    if cntr == "" {
       Ctr  = 10  	
    }else{
       Ctr = Str2int(cntr)	
    }

    fmt.Println("IP node            : ", Ip)
    fmt.Println("Count transaction  : ", Ctr)
  
    // U:= SetUnlock(Ip,Port, ACC)
    // fmt.Println("Transaction return : ", U)

    // sumtx:=fmt.Sprintf("0x%x\n", math.Float32bits(0.00000000100)) 
    // fmt.Println("Summ wei:", sumtx)
    // fmt.Println(hex(1123, true))
    // return

    GlobalPass1 := "Gerda1000" // Gerda1000"
    url         := "http://" + Ip + ":" + Port + "/"

	// Разблокировка (включена - отключена) 
	Unlckk:=true

	if Unlckk {
		// Unlock all accounts
	    for _, Accr:=range(Accs){
	        UnlckAccount(w, url, Accr, GlobalPass1)
	    }
	}


    // Gas Price
	// 0x9184         - old value
	// 0x3b9aca00     - 
	GasPrice:=Get_gas_price(url)
	fmt.Println("Gas price : ", GasPrice)

    // **********************************************
    // Cicle for node by transcation
    // 14-12-2018 16:52:10
    // **********************************************
	for i := 0; i < Ctr; i++ {
         fmt.Println("Tx *****************************************************",i, "\n") 
                  
         go func (){
            // gettx := SetTransaction(Ip, Port, Accs[0], Accs[1], "0x76c0", GasPrice, "0x1308", GlobalData)
            gettx := SetTransaction(Ip, Port, Accs[0], Accs[1], "0x76c0", GasPrice, "0x1308", GlobalData)
            fmt.Println("Tx : ", gettx)
            time.Sleep(time.Millisecond*600)
        }()
	}
    
    fmt.Println("Status: Test End !")
}

// ************************************************************************ 
// Нагрузочное тестирование
// Usage:  // http://localhost:8545/htransact/1111
// ************************************************************************
func SetTransact_hight_old(w http.ResponseWriter, req *http.Request) {
    defer func(){
    	  recover()
    }()
    
    // Количество транзакций
    cntr := req.URL.Path[len("/htransact/"):]

    // Accounts
    Accs :=[]string{"0xb614c2f82b24b96b6241e54a65c0c0e62f62be5d",
                    "0x977058a862cd16ef557b3844b61e432a6a0f0378",
                    "0x42546e93f0f8a64005861a14655e213469fcec9f",
                    "0x4cc910613a6d4af4eb6be633b138a908721114c8",
                    "0x4a4c33751a2841ad880dfca49560ad6a118129ad",
                    "0xb94e5a035108e9356a811dd790c2c9eeb33d73c2",
                    "0xd19d53673976f8e901a817fb7fd1781dc86d141a"}


    // Settings for testing
    Ip   := Gess01        // Ip node  
    Port := "8545"        // Port by default 
    Port  = MPort         // Port resset
    // From := ACC        // Есть деньги
    // To   := ACC2       // My
    Ctr  := 10            // Count transaction 
    Rnd  := 0

    // Count transaction
    if cntr == "" {
       Ctr  = 10  	
    }else{
       Ctr = Str2int(cntr)	
    }

    fmt.Println("IP node            : ", Ip)
    fmt.Println("Count transaction  : ", Ctr)

   
    // U:= SetUnlock(Ip,Port, ACC)
    // fmt.Println("Transaction return : ", U)

    // sumtx:=fmt.Sprintf("0x%x\n", math.Float32bits(0.00000000100)) 
    // fmt.Println("Summ wei:", sumtx)
    // fmt.Println(hex(1123, true))
    // return

      GlobalPass1 := "Gerda1000"
      url         := "http://" + Ip + ":" + Port + "/"
    

 Unlckk:=false

if Unlckk{
	// Unlock all accounts
    for _, Accr:=range(Accs){
        UnlckAccount(w, url, Accr, GlobalPass1)
    }
}


    // Gas Price
	// 0x9184         - old value
	// 0x3b9aca00     - 
	GasPrice:=Get_gas_price(url)
	fmt.Println("Gas price : ", GasPrice)
	            


    // **********************************************
    // Cicle for node by transcation
    // 14-12-2018 16:52:10
    // **********************************************
    // func (){
		    for i := 0; i < Ctr; i++ {


                
                  go func (){

		              // start:=time.Now()
		              // fmt.Println("ID Transaction : ",i)

                      // gettx := SetTransaction(Ip, Port, From, To, "0x76c0", GasPrice, "0x1308", GlobalData)
                  	  gettx := SetTransaction(Ip, Port, Accs[0], Accs[1], "0x76c0", GasPrice, "0x1308", GlobalData)
                  	  fmt.Println("Id 0", i, " Hash Tx:", gettx)
                  	  time.Sleep(time.Millisecond*600)

                      
                      fmt.Println("Tx id: ______________________________________________________",i, "\n") 
                      time.Sleep(time.Second * 1)

                      //Accss:=""

		              // for ss := 0; ss < 10; ss++ {
                       for ss, Accr:=range(Accs){
                       	   


                           ttx := SetTransaction(Ip, Port, Accs[0], Accr, "0x76c0", GasPrice, "0x1308", GlobalData)
                           fmt.Println("Id acc ", ss, "Accounts :", Accr, "Tx:", ttx)
                           time.Sleep(time.Millisecond*500)

                      }    

                      fmt.Println("Разброс денег по всем счетам")
                      fmt.Println("**************************************************")
                      fmt.Println("")



    
                      //  for ss, Accr:=range(Accs){
                      //      ttx := SetTransaction(Ip, Port, Accs[0], Accs[Randoms(100)], "0x76c0", GasPrice, "0x1308", GlobalData)
                      //      fmt.Println("Id acc ", ss, "Accounts :", Accr, "Tx:", ttx)
                      //      time.Sleep(time.Millisecond*500)

                      // }    





                          // Открыты счета с одним паролем !!!
                  	      // важно контролировать счета с одним паролем
                  	      // и вначале открывать счета UNLOCK все
                  	      // потом разбрасывать деньги по остальным 
                          // Изначально на первом счету должны быть монеты


                  	      // TODO:
                  	      // * в произвольном порядке
                  	      // * Проверка на первом счету денег
                  	      // * Запустить все в цикле
                  	      // * Брать из базы данных 
                  	      // * Записать в базу при открытии
                  	      // * Ловить юлок с количеством транзакции
                  	      // * создать функцию рендомно создающая от 0-10 числа


                          ttx := SetTransaction(Ip, Port, Accs[0], Accs[1], "0x76c0", GasPrice, "0x1308", GlobalData)
                          fmt.Println("Id acc 1", i, "Accounts fix:", Accs[0], "Tx:", ttx)
                          time.Sleep(time.Millisecond*500)


                          ttx = SetTransaction(Ip, Port, Accs[0], Accs[3], "0x76c0", GasPrice, "0x1308", GlobalData)
                          fmt.Println("Id acc 2", i, "Accounts fix:", Accs[2] , "Tx:", ttx)
                          time.Sleep(time.Millisecond*500)

                          ttx = SetTransaction(Ip, Port, Accs[0], Accs[4], "0x76c0", GasPrice, "0x1308", GlobalData)
                          fmt.Println("Id acc 3 ", i, "Accounts fix:", Accs[3], "Tx:", ttx)
                          time.Sleep(time.Millisecond*500)

                          ttx = SetTransaction(Ip, Port, Accs[0], Accs[6], "0x76c0", GasPrice, "0x1308", GlobalData)
                          fmt.Println("Id acc 4", i, "Accounts fix:", Accs[5], "Tx:", ttx)
                          time.Sleep(time.Millisecond*500)

                          ttx = SetTransaction(Ip, Port, Accs[0], Accs[7], "0x76c0", GasPrice, "0x1308", GlobalData)
                          fmt.Println("Id acc 5", i, "Accounts fix:", Accs[6], "Tx:", ttx)
                          time.Sleep(time.Millisecond*500)
                          
                          ttx = SetTransaction(Ip, Port, Accs[0], Accs[8], "0x76c0", GasPrice, "0x1308", GlobalData)
                          fmt.Println("Id acc 6", i, "Accounts fix:", Accs[7], "Tx:", ttx)
                          time.Sleep(time.Millisecond*500)
                          
                          ttx = SetTransaction(Ip, Port, Accs[0], Accs[9], "0x76c0", GasPrice, "0x1308", GlobalData)
                          fmt.Println("Id acc 7", i, "Accounts fix:", Accs[8], "Tx:", ttx)
                          time.Sleep(time.Millisecond*500)


                          Rnd = Rendom(10)
                          ttx = SetTransaction(Ip, Port, Accs[Rnd], Accs[2], "0x76c0", GasPrice, "0x1308", GlobalData)
                          fmt.Println("Id acc 2", i, "Accounts fix:", Accs[2] , "Tx:", ttx)
                          time.Sleep(time.Millisecond*500)

                          Rnd = Rendom(10)
                          ttx = SetTransaction(Ip, Port, Accs[Rnd], Accs[3], "0x76c0", GasPrice, "0x1308", GlobalData)
                          fmt.Println("Id acc 2", i, "Accounts fix:", Accs[2] , "Tx:", ttx)
                          time.Sleep(time.Millisecond*500)

                          Rnd = Rendom(10)
                          ttx = SetTransaction(Ip, Port, Accs[Rnd], Accs[2], "0x76c0", GasPrice, "0x1308", GlobalData)
                          fmt.Println("Id acc 2", i, "Accounts fix:", Accs[2] , "Tx:", ttx)
                          time.Sleep(time.Millisecond*500)



                       // fn(Accs[0], Accs[9]) 
                       // ddd:=Get_gas(gettx)
                       // fmt.Println("-----------",ddd)
                   }()


                   /*

                      fn:=function(Acfrom, Accto){
                          ttx = SetTransaction(Ip, Port, Acfrom, Accto, "0x76c0", GasPrice, "0x1308", GlobalData)
                          fmt.Println("Id acc 7", i, "Accounts fix:", Accs[8], "Tx:", ttx)
                          time.Sleep(time.Millisecond*500)
	                  }
                   */

			    	   // T  := SetTransaction(Ip, Port, From, To, "0x76c0", "0x9184e72a", "0x1388", GlobalData)
			           // fmt.Println("Transaction return : ",  T)
			     
			           // finish  := time.Now()
			           // duration:= finish.Sub(start)
			           // fmt.Println("Duration :", duration)
		    }
	// }()	    

  

   //  Ip   = Node05
   //  From = ACC2
   //  To   = ACC3

   //  fmt.Println("***************************************************")
   //  fmt.Println("NODE 5 # ", Ip)
   //  fmt.Println("***************************************************")

   // func (){
	  //   for i := 0; i < 10; i++ {


	  //       // start:=time.Now()
	  //       // fmt.Println("ID Tx: ", i)

	     
	  //           T  := SetTransaction(Ip, Port, From, To, "0x76c0", "0x9184e72a", "0x1388", GlobalData)
	  //           fmt.Println("Transaction return : ",  T)
	  //           bl:=BalanceInfo(Ip, Port, From)
	  //           time.Sleep(time.Millisecond*1000) 
	  //           fmt.Println("Баланс : ", bl)
	        

	  //           finish  := time.Now()
	  //           duration:= finish.Sub(start)
	  //           fmt.Println("Duration :", duration)
	        
	  //   }

   // }()
  
    fmt.Println("Status: Test End !")
}

// ************************************************************
// Unlock account
// ************************************************************
func UnlckAccount(w http.ResponseWriter, url, Acc, Pas string){
	// UNLOCK TRANSACTION 
    // GlobalPass
    // ====================================================


	    // url            := "http://" + Ip + ":" + Port + "/"
		unlocked       := strings.NewReader(`{"id":100, "jsonrpc":"2.0", "method": "personal_unlockAccount", "params":["` + Acc + `","` + Pas + `", 15000]}`)
		unlreq, errt   := http.NewRequest("POST", url, unlocked)

		if errt!=nil{
		   // fmt.Fprintf(w, "Test brocken") 
		   w.Write([]byte(fmt.Sprintf("Error")))
		   log.Println("Error client : ", errt.Error())
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

		fmt.Println("UNLOCK : ", string(bodys))
    // ===================================
    // End unlock
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

// **********************************************
// Get Gas
// На вход hash транзакции
// **********************************************
func Get_gas(Tx string) string {

    url:= "http://" + Gess01 + ":" + MPort + "/"
	r:=Ex(url,"eth_gasPrice",Tx)
	return r
}

// ************************************************************************ 
// Нагрузочное тестирование в 100 итераций
// Нагрузочное тестирование в 1000 
// ************************************************************************
func SetTransact_infinity(w http.ResponseWriter, req *http.Request) {
    // IP:PORT
    Ip   := "18.188.111.198"
    Port := "8545"
    From := ACC      // Есть деньги
    To   := ACC2     //-My

    // U:= SetUnlock(Ip,Port, ACC)
    // fmt.Println("Transaction return : ", U)

    for  {
         start:=time.Now()
    	 func (){
            T  := SetTransaction(Ip, Port, From, To, "0x76c0", "0x9184e72a", "0x1388", GlobalData)
            fmt.Println("Transaction return : ",  T)
        }()

        finish   := time.Now()
        duration := finish.Sub(start)
        fmt.Println("Infinity :", duration)
    }

    fmt.Println("Test Infinity End !")
}


// *****************************************************************************************************
// Перечисление Со счета А-Б
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func SetTransaction(Ip, Port, From, To, Gas, Gasprise, Value, Data string) string {

    // Url with port
    url     := "http://"+Ip+":"+Port+"/"
    unlc    := false
    
    // Unlocked false
    if unlc {
		    unlocked       :=  strings.NewReader(`{"id":100, "jsonrpc":"2.0", "method": "personal_unlockAccount", "params":["`+From+`","`+GlobalPass+`", 15000]}`)
		    unlreq, errt   :=  http.NewRequest("POST", url, unlocked)

		    if errt!=nil{
		       log.Println("Error client : ", errt.Error())
		       return ""
		    }

		    unlreq.Header.Add("content-type",  "application/json")
		    unlreq.Header.Add("cache-control", "no-cache")

		    resu, err_unl := http.DefaultClient.Do(unlreq)
		    if err_unl!=nil{
		       log.Println("Error client : ", err_unl.Error())
		       return ""
		    }
		   
		    defer resu.Body.Close()

		    bodys, errbs := ioutil.ReadAll(resu.Body)
		    if errbs!=nil{
		       log.Println("Error read body", errbs.Error())
		       return ""
		    }

		    fmt.Println("UNLOCK : ", string(bodys))
    }
 


    // Transaction
    Params  := `"params":[{"from":"`+From+`","to":"`+To+`","gas":"`+Gas+`","gasPrice":"`+Gasprise+`","value": "`+Value+`","data": "`+Data+`"}]`
    payload :=  strings.NewReader(`{"id":104, "jsonrpc":"2.0", "method":"eth_sendTransaction", `+Params+`}`)

    req, _  :=  http.NewRequest("POST", url, payload)
    req.Header.Add("content-type",  "application/json")
    req.Header.Add("cache-control", "no-cache")

     res, err := http.DefaultClient.Do(req)
     if err!=nil{
        log.Println("Error client", err.Error())
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


// ***************************************************************
// Preview all accounts balance
// Date |: 28-11-2018 19:51:16
// Просомтр остатков на счетах
// 0x97582e614f2f8e38f72a9d8d757a18057e68d337 - yes
// 0x5c81a844fa4dbe3f3cba50f87821cb572230da7e - My 1 (Token)
// 0xbec03d897214d1016b860b0b3ebbd4290ec8deda - My 2
// https://doc.rust-lang.org/rust-by-example/primitives/tuples.html
// https://tomato-timer.com/
// 
// http://localhost:8545/account/
// ***************************************************************
func SetAccount(w http.ResponseWriter, req *http.Request) {
    
    // IP:PORT
    Ip       := "18.188.111.198"
    Port     := "8545"
    Nodename := Node02

    // 0x1900cd88080ede312a6b7910e66e3c3cea75a479
    // 0x1dccd1ef104eefa133beaf8237eb5b1845dca2e6
    // 0x48f00b326e017af888cedd7883d023587c6de118

    // Open (1) 10 accounts for test 
    for i := 0; i < 100; i++ {

        // func(){    
	        
	        //Ip  = "188.190.240.195"
            T  := AddAccount(Nodename, Port)
            log.Println("ID:", i, "Account 3:", T)
        // }()
    }


    fmt.Println("Ready.", Nodename)
    return

    // 0x1dccd1ef104eefa133beaf8237eb5b1845dca2e6
    // 0x48f00b326e017af888cedd7883d023587c6de118
   
    run5:=true
    if run5 {
       Ip = Node05
       T  := AddAccount(Ip, Port)
       fmt.Println("Account Node 04 : ", T)
    }
}

// ***************************************************************
// Account Create
// ***************************************************************
func AddAccount(Ip, Port string) string {
	 
     p       := `{"id":3, "jsonrpc":"2.0", "method": "personal_newAccount", "params":["` + GlobalPass + `"]}` 
     url     := "http://" + Ip + ":" + Port + "/"
     payload :=  strings.NewReader(p)
     req, _  :=  http.NewRequest("POST", url, payload)

     req.Header.Add("content-type",  "application/json")
     req.Header.Add("cache-control", "no-cache")

     res, err := http.DefaultClient.Do(req)
     if err!=nil{
        log.Println("Error client", err.Error())
        return "Error client connect"
     }
     // defer res.Body.Close()

     body, errb := ioutil.ReadAll(res.Body)
    
     if errb!=nil{
        log.Println("Error read body", errb.Error())
        return "Error read body"
     }

     res.Body.Close()     
     return string(body)
}


// ***************************************************************
// Exec Operation
// ***************************************************************
func getContent(ctx context.Context) {

    // Nodes  
    ipport:= "http://18.217.164.134:8545"  // Gess 03
    ipport = "http://18.188.240.197:8545"  // Gess 02
    ipport = "http://18.224.11.186:8545"   // Gess 04
    ipport = "http://18.224.106.72:8545"   // Gess 05
    ipport = "http://18.188.111.198:8545"  // Gess 01
 // ipport = "http://18.222.125.100:8545"  // Test node 

    // ipport="https://stackoverflow.com/questions/44911726/the-best-way-to-set-the-timeout-on-built-in-http-request"
    // payload := strings.NewReader(`{"jsonrpc":"2.0","method":"net_version", "id":1}`)

    /*
       Открытые кошельки :
       1 {"jsonrpc":"2.0","id":2,"result":"0x32907d682e45593a0a44a34d5c1f8d4aa9f7301e"}
       2 {"jsonrpc":"2.0","id":2,"result":"0x567f7963e42b54b118b76e4ffaaa6af5029b5ada"}
       3 {"jsonrpc":"2.0","id":2,"result":"0x600dec2e3ade5227e364689f39aeee79a8083c8a"}
       4 {"jsonrpc":"2.0","id":2,"result":"0xa4061a85d66b40d0750fa924614dd10742486fe9"}
       5 {"jsonrpc":"2.0","id":2,"result":"0xe1b69477ab8c2da054e9ab6c758454bfa6766bf6"}

       18.188.111.198 
       "0x92f81d2cddb753c50822e5d4db41e5c59ddccb82","0x19ee18e1bb682cd91ed80e47d34c5e38ee11d642","0x816d3cf337ade88b1d7f6fcae924a5f71a58fe7c","0x2ea5d7df797eb59abb125af82c7b4a7aa78c07ba","0xc70a9086b65e92514c423346c3bc018d67c41fb3","0x7637c25156ed10996a27bb09b482129f3a37ffa9","0x1f544dd0100e005b9082e0e3b0ecf5ac65debf1b","0x90df3764818b8c4da1ae9759700fe5a2bd7428cf","0xe184411f45b84fc1c28d7c41c4c27ef5bcc9ecde","0x14e0dc2983b7964ecbcd4dd45c6f03bc471f8b3c","0xa5614fa5d7808289d58c31cb6c8bca61cf7426c9","0xbc42851e21f66b648bd19d26a44a76bf11f0f7a2","0x74061b5254318661361526c7fce3cb62197e2ec8","0x97582e614f2f8e38f72a9d8d757a18057e68d337","0xdfd221d1d1e3dfe9e563ae312ac247900ae7dfd7","0x964d55c138f639e8617cfd788b2ea13a90213ff0","0x90ec0b5f932bf279de9bba675ae7779dde4ff36c","0x917af71f929d9bd5a92c60c526d4dd618c7db7e8","0x767199624c5c318043c3bb21ef5015451c57d85c"
    */

    Par  := make(map[int]string)

    Par[1]=`{"id":1, "jsonrpc":"2.0", "method":"eth_accounts","params":[]}`                                                            // Просомтор всех счетов
    Par[2]=`{"id":2, "jsonrpc":"2.0", "method":"personal_newAccount", "params":["`+GlobalPass+`"]}`                                    // Открытие счетов
    Par[3]=`{"id":3, "jsonrpc":"2.0", "method":"eth_getBlockNumber"}`                                                                  // Номер блока
    Par[4]=`{"id":4, "jsonrpc":"2.0", "method":"eth_getBalance","params":["0x19ee18e1bb682cd91ed80e47d34c5e38ee11d642", "latest"]}`    // Получение баланса

    payload := strings.NewReader(Par[4])

    req, err := http.NewRequest("GET", ipport, payload)
    if err != nil {
       log.Fatal(err)
    }

    req.Header.Add("content-type",  "application/json")
    req.Header.Add("cache-control", "no-cache")

    ctx, cancel := context.WithDeadline(ctx, time.Now().Add(3 * time.Second))
    defer cancel()

    req.WithContext(ctx)

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
       log.Fatal(err)
    }
    defer resp.Body.Close()

    Cnt++
    fmt.Println(Cnt)
    io.Copy(os.Stdout, resp.Body)
}

// ***************************************************************
// Get Network Id
// ***************************************************************
func NetworkIds(w http.ResponseWriter, req *http.Request) {
     fmt.Println("Get process...")

     ctx := context.Background()

     for i := 0; i < 1; i++ {
         fmt.Println("Reguest : ",i)
         getContent(ctx)    
     }
}

// **********************************************
// Общая функция для запросов
// **********************************************
func Ex(url, method, params string) string {

	        var Rs Mst

	        // Get request
            GetGasPrice    :=  strings.NewReader(`{"jsonrpc":"2.0","method":"`+ method +`","params":[` + params + `],"id":100}`)
		    GasReq, errt   :=  http.NewRequest("POST", url, GetGasPrice)

		    if errt!=nil{
		       log.Println("Error gasprice get : ", errt.Error())
		       return ""
		    }
		    GasReq.Header.Add("content-type",  "application/json")
		    GasReq.Header.Add("cache-control", "no-cache")

            
            // Body
		    resGas, err_gasprice := http.DefaultClient.Do(GasReq)
		    if err_gasprice!=nil{
		       log.Println("Error gas price body return : ", err_gasprice.Error())
		       return ""
		    }
		    defer resGas.Body.Close()
            
            // Read body
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


// ***************************************************************
// Account Create
// Cmd(`{"id":3, "jsonrpc":"2.0", "method": "personal_newAccount", "params":["Gerda1000"]}`)
// ***************************************************************
func Cmd(Ip, Port ,p string) string {
    url        := "http://" + Ip + ":" + Port + "/"
    payload    :=  strings.NewReader(p)
    req, errd  :=  http.NewRequest("POST", url, payload)

    if errd!=nil{
       log.Println("Error client", errd.Error())
       return "Error client connect"
    }

    req.Header.Add("content-type",  "application/json")
    req.Header.Add("cache-control", "no-cache")

    res, err := http.DefaultClient.Do(req)
    if err!=nil{
        log.Println("Error client", err.Error())
        return "Error client connect"
    }
    defer res.Body.Close()

    body, errb := ioutil.ReadAll(res.Body)
    
    if errb!=nil{
       log.Println("Error read body", errb.Error())
       return "Error read body"
    }

    // res.Body.Close()     
    return string(body)
}


// **********************************************************
// Открытие новых счетов
// New Account
// Node 01
// 0xecb8ed418415706994614c8271ae6876391a0fd0  
// 0x22a813f889171add159fb93e545c65c109c3928c
// /ac/new/
// **********************************************************
func NewAccount_wrk(w http.ResponseWriter, req *http.Request) {
    
    FgRed   := color.New(color.Bold, color.FgRed).PrintlnFunc()
    FgGreen := color.New(color.Bold, color.FgGreen).PrintlnFunc()
    FgGreen("Open new accounts")


    //log.Println("Open new accounts \n")

    account      := ""
    accountcount := 10

    o := `{"id":4,  "jsonrpc":"2.0", "method": "personal_newAccount", "params":["`+GlobalPass+`"]}`
     // for i := 0; i < 10; i++ {
     //     r =Cmd(Node00, Portnode, `{"id":1, "jsonrpc":"2.0", "method": "personal_newAccount", "params":["Gerda1000"]}`)
     //     log.Println(Node00, r)
	
     // }
     
     /*

     log.Println("Node 02 :", Node02, "\n")
     r=Cmd(Node02, Portnode, `{"id":2, "jsonrpc":"2.0", "method": "personal_newAccount", "params":["Gerda1000"]}`)
     log.Println(Node02, r)

     
     log.Println("Node 03 :", Node03, "\n")
     r=Cmd(Node03, Portnode, `{"id":3, "jsonrpc":"2.0", "method": "personal_newAccount", "params":["Gerda1000"]}`)
     log.Println(Node03, r)
    */

     //  0x97d794d474e15715f607400a5bd325447977d0e1 - yes money


    // Node 04
    for i := 0; i < accountcount; i++ {
         fmt.Println("Node 04 :", Node04, "\n")
         account = Cmd(Node04, Portnode, o)
         

         if !strings.Contains(account,"Error"){
             D_accounts_insert(account, "Node 04","0.000000000000000000","Added")
             FgGreen(account)
         }else{
             FgRed("Error create account ") 
             D_log("Error create account", "account", "","")
         }    

         fmt.Println("...........................................\n")
    }
         

    // Node 05 
    for i := 0; i < accountcount; i++ {
         // 0x1b4540f309ca52b317dbc963c226ffa469c29280 - yes money
         fmt.Println("Node 05 :", Node05, "\n")
         account = Cmd(Node05, Portnode, o)
         //FgRed("New Account :", account) 


         if !strings.Contains(account,"Error"){
             D_accounts_insert(account, "Node 05","0.000000000000000000","Added")
             FgGreen(account)
         }else{
             FgRed("Error create account ") 
             D_log("Счет не открыт ", "account", "","")
         }    
         fmt.Println("...........................................\n")
    }     



    // Miner02 
    for i := 0; i < accountcount; i++ {
         fmt.Println("Miner02 :", Miner02, "\n")
         account = Cmd(Node05, Portnode, o)

         if !strings.Contains(account,"Error"){
             D_accounts_insert(account, "Node 05","0.000000000000000000","Added")
             FgGreen(account)
         }else{
             FgRed("Error create account ") 
             D_log("Счет не открыт ", "account", "","")
         }    
         fmt.Println("...........................................\n")
    }     
     

    // Miner03
    for i := 0; i < accountcount; i++ {
         fmt.Println("Miner03 :", Miner03, "\n")
         account = Cmd(Node05, Portnode, o)

         if !strings.Contains(account,"Error"){
             D_accounts_insert(account, "Miner 03","0.000000000000000000","Added")
             FgGreen(account)
             D_log("Счет открыт "+account, "account", "","")
         }else{
             FgRed("Error create account ", "account") 
             D_log("Счет не открыт ", "account", "","")
         }    
         fmt.Println("...........................................\n")
    }     


    // Miner04
    for i := 0; i < accountcount; i++ {
         fmt.Println("Miner04 :", Miner04, "\n")
         account = Cmd(Node04, Portnode, o)

         if !strings.Contains(account,"Error"){
             D_accounts_insert(account, "Node 04","0.000000000000000000","Added")
             FgGreen(account)
             D_log("Счет открыт "+account, "account", "","")
         }else{
             FgRed("Error create account ") 

         }    
         fmt.Println("...........................................\n")
    }     


    fmt.Println("Finished open accounts")
    D_log("Finished open accounts", "account", "","")

    // O4 NODE
    // 0x9108f1bd24823337b460603c93b171e61f66a18f
    // 0x6f3676ae51579e256e890a654ac56f3b3cdbcf06
    // 0xafb69147285e7a1e1e8ece7b5879b3f0021214a4

    // 05 Node
    // 0x11ac8b9dc5f7843e1b1f91d74611ab169448e4c4
    // 0x1e44ca252f0d6567c5db3e2404c60ba83a616d4b  
    // 0xf99297663768a0b0d183c657e06dd4060da8a217
}


// ********************************************************
// Control txtpool
// ********************************************************
func TxtPool(w http.ResponseWriter, req *http.Request){
	     r := Cmd(Node00, Portnode, `{"id":1, "jsonrpc":"2.0", "method": "txpool_content"}`)
         log.Println(Node00, r)
}


// ********************************************************
// 
// ********************************************************
func ViewT3es(w http.ResponseWriter, req *http.Request) {
 // Set color
    Whites  := color.New(color.FgWhite).PrintlnFunc()
    FgMag   := color.New(color.FgHiYellow)
    Cyan    := color.New(color.Bold, color.FgCyan)
    FgRed   := color.New(color.Bold, color.FgRed).PrintlnFunc()
    FgGreen := color.New(color.Bold, color.FgGreen).PrintlnFunc()

    FgGreen("Proxy server")
    Whites("Version: 1.01 (Testing)")


    FgMag.Print("Active node : ")
    FgRed("Host")
    Cyan.Printf("Listen proxy port %s \n", "Samp")

}


// **************************************************************
// Тестирование работы с базой
// **************************************************************
func Db_start(w http.ResponseWriter, req *http.Request) {
    // user := map[string]interface{}{"Name": "GGG", "Age": 18, "Birthday": time.Now()}
    // user:=User{Name: "Прохоров", Age: 118, Birthday: time.Now()}
    // D_insert(user)
    //D_scan()   
    tm:=time.Now().Format("03-02-206 15:04:05")
    D_scan_accounts(tm)
    D_accounts_insert("0x11ac8b9dc5f7843e1b1f91d74611ab169448e4c4-000", "014","0.000000000000","Добавлен" )

}

// *****************************************************
// Added New Accounts
// *****************************************************
func D_accounts_insert(Account, Node, Sum, Note  string){
     var accounts account

     db, err := gorm.Open("sqlite3", "db/tst.db")
     if err!=nil{
        fmt.Println("Erroro connect to db", err.Error())
     }
     defer db.Close()
    
     accounts.Title      = Account
     accounts.Node       = Node
     accounts.Createat   = time.Now().Format("02-01-2006 14:05")
     accounts.Sum        = Sum
     accounts.Note       = Note
     accounts.Descript   = "New Accounts"
     db.Create(&accounts)
     fmt.Println("Account ", Account, " added to database." )
}

// *****************************************************
// Loop Accounts test
// *****************************************************
func D_scan_accounts_test(w http.ResponseWriter, req *http.Request){
	 s:=time.Now()

	for i := 0; i < 1000; i++ {

	     s1:=time.Now()	 
	     fs:=string(i)
	     D_scan_accounts(fs)
		 f1:=time.Now()   
	     l1:=f1.Sub(s1)
	     fmt.Println(i," record ...... ",l1 )

	}


f:=time.Now()
l:=f.Sub(s)
fmt.Println("All ",l )
	
}


// *****************************************************
// Accounts - scan all 
// *****************************************************
func D_scan_accounts(Tm string){

    var acc []account
    db, err := gorm.Open("sqlite3", "db/tst.db")
    if err!=nil{
       fmt.Println("Error connect to db", err.Error())
    }
    defer db.Close()
    db.Table("accounts").Find(&acc)
   
   for _,ac:=range(acc){
        // fmt.Println(ac.Title, ac.Node)
        D_log(ac.Title, ac.Node, "Added new accounts",Tm)
        
    } 
}


// *****************************************************
// Scan Nodes
// *****************************************************
func D_scan_node(){
    db, err := gorm.Open("sqlite3", "db/tst.db")
    
    if err!=nil {
       fmt.Println("Error connect to db", err.Error())
    }
    defer db.Close()

    var nd []node
    db.Table("nodes").Find(&nd)
   
   for _,i:=range(nd){
        fmt.Println(i.Title, i.Ip)
        
    } 
}

// *****************************************************
// 
// *****************************************************
func D_insert(Dat User){
     db, err := gorm.Open("sqlite3", "db/tst.db")
     if err!=nil{
        fmt.Println("Erroro connect to db", err.Error())
     }
     defer db.Close()

     user := Dat
     db.Create(&user)
     fmt.Println("Ok")
     
     // db.Save(&user) // will set `UpdatedAt` to current time
     // db.Model(&user).Update("name", "Прохоров") // wil
}

// **************************************************
// Adding new users
// **************************************************
func DbC(){
      db, err := gorm.Open("sqlite3", "db/tst.db")
     if err!=nil{
        fmt.Println("Erroro connect to db", err.Error())
     }
     defer db.Close()

     user := User{Name: "Jinzhu", Age: 18}
     // user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}
     // db.NewRecord(user) // => returns `true` as primary key is blank
     db.Create(&user)
     //db.NewRecord(user) // => re
     fmt.Println("Ok")
}


type Account struct{
     Title string 
     Note  string

}




// **************************************************
// Adding new users
// **************************************************
func DB_Add_Account(){
      db, err := gorm.Open("sqlite3", "db/tst.db")
     if err!=nil{
        fmt.Println("Erroro connect to db", err.Error())
     }
     defer db.Close()

     accounts := Account{Title: "Jinzhu", Note: "Node"}
     // user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}
     // db.NewRecord(user) // => returns `true` as primary key is blank
     db.Create(&accounts)
     //db.NewRecord(user) // => re
     fmt.Println("Ok")
}

// *****************************************************
// 
// *****************************************************
func D_log(Acc, Node, Note, Typ string ){
     var Log logs
     db, err := gorm.Open("sqlite3", "db/tst.db")
     if err!=nil{
        fmt.Println("Erroro connect to db", err.Error())
     }
     defer db.Close()

     Log.Note = Note 
     Log.Type = Typ 
     Log.Acc  = Acc
     Log.Node = Node
     db.Create(&Log)
     // fmt.Println("Ok")
     // db.Save(&user) // will set `UpdatedAt` to current time
     // db.Model(&user).Update("name", "Прохоров") // wil
}


// *********************************************************************
// Admin panel
// *********************************************************************
func Api_admin(w http.ResponseWriter, req *http.Request) {
     SP(w,"node.html")
}

// *********************************************************************
// Open new account
// /ac/one/
// *********************************************************************
func NewAccount_one (w http.ResponseWriter, req *http.Request) {
  SP(w,"account.html")
}

// *****************************************************
// Load test panel
// *****************************************************
func Testpanel(w http.ResponseWriter, req *http.Request){
     SP(w,"testpanel.html")
}

// *****************************************************
// Start page
// *****************************************************
func SP(w http.ResponseWriter, pagename string){
    fp        := path.Join("tmp", pagename)
    tmpl, err := template.ParseFiles(fp, "tmp/main.html")
    Err(err, "Error template execute.")
    errf      := tmpl.Execute(w, nil)
    Err(errf, "Error template execute.")
}

// *****************************************************
// Show all transaction
// *****************************************************
func Block_info(w http.ResponseWriter, req *http.Request){
    // GetInfoblock("0x5b96f5a92c5074fd81cbbfbc20fff49f7530a12dd23f541c117a25b78604b914")
	// GetInfoblock("3183203")
	
	// t1:=IsTransInBlock("0x5b96f5a92c5074fd81cbbfbc20fff49f7530a12dd23f541c117a25b78604b914","4cc3a42c3")
	// t2:=IsTransInBlock("0x5b96f5a92c5074fd81cbbfbc20fff49f7530a12dd23f541c117a25b78604b914","5c3c174f95")
	// t3:=IsTransInBlock("0x5b96f5a92c5074fd81cbbfbc20fff49f7530a12dd23f541c117a25b78604b914","24b1d29c53d99b")
	// t4:=IsTransInBlock("0x5b96f5a92c5074fd81cbbfbc20fff49f7530a12dd23f541c117a25b78604b914","14b1d29c53d99b")
    // fmt.Println(t1,t2,t3, t4)

    // BlockInformationTransactions("0x5b96f5a92c5074fd81cbbfbc20fff49f7530a12dd23f541c117a25b78604b914")
    //BlockInformationTransactions("0x8c686d67dd7a185fb08a5a3662384997ec9499f539fa91efeabe8628453040e9")
    
    // Information oabout last5 block
    GetLastBlock(req)

    // Information about transaction in blocck by hash 
    // BlockInformationTransactions("0x65e6afd061262d90cb89d5e31eaf1a20b7a06bc3cfde0dae663d6ba4fec9a250")

}


// ***************************************
// Есть ли транзакция в блоке
// ***************************************
func IsTransInBlock(Hash, Tr string) bool {
     block := Cmd(Nodelocal, "8555", `{"id":3, "jsonrpc":"2.0", "method": "eth_getBlockByHash", "params":["`+Hash+`",true]}`)
     if strings.Contains(block,Tr ){
       return true
     }
     return false
}


// *************************************************************************
// Информация по блоку
// прокрутка всех транзакций в блоке
// BlockInformationTransactions("0x5b96f5a92c5074fd81cbbfbc20fff49f7530a12dd23f541c117a25b78604b914")
// *************************************************************************
func BlockInformationTransactions(Hash string){

    fmt.Println("Hash tx", Hash) 
    fmt.Println("*****************************************") 

    ff:=GetInfoblocktrasactions(Hash)
      // fmt.Println(ff)

      for _, rb:=range(ff){
      	  txs:=rb.(map[string]interface{})
      	  th:=txs["hash"].(string)
      	  ts:=txs["nonce"].(string)
      	  fr:=txs["from"].(string)
      	  to:=txs["from"].(string)
      	  bh:=txs["blockHash"].(string)
      	  bn:=txs["blockNumber"].(string)
      	  id:=txs["transactionIndex"].(string)
      	  
      	  

          fmt.Println("Transaction ", th,ts,fr,to,bh,"\n") 

          fmt.Println("Transaction ", id,bn) 
          

      }	
}




// https://github.com/miguelmota/go-coinmarketcap/tree/master/v2

// doReq HTTP client
func doReq(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}

	return body, nil
}

// makeReq HTTP request helper
func makeReq(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := doReq(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}
