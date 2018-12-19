package main

import (    
    "fmt"
    "time"
    "log"
    "strings"
    "encoding/json"
    "net/http"
    "io/ioutil"

)

// ****************************************************************
//  Blocks info 
//  Pending 
//  earliest
// ****************************************************************
func GetInfoblocktrasactions(Hash string) []interface{} {
    block := Cmd(Nodelocal, "8555", `{"id":3, "jsonrpc":"2.0", "method": "eth_getBlockByHash", "params":["`+Hash+`",true]}`)
    var Md Mst

    errj:=json.Unmarshal([]byte(block), &Md)
    if errj!=nil{
       log.Println("Error ummarshaling.", errj.Error())
    }

    Fs  := Md["result"].(map[string]interface{})
    Fbg := Fs["transactions"].([]interface{})
    return Fbg
}   

// ****************************************************************
//  Blocks info 
//  Pending 
//  earliest
// ****************************************************************
func GetInfoblock(Hash string){
    // P:=[]string{"latest", "earliest", "pending"}

    block := Cmd(Nodelocal, "8555", `{"id":3, "jsonrpc":"2.0", "method": "eth_getBlockByHash", "params":["`+Hash+`",true]}`)
    // block := Cmd(Nodelocal, "8555", `{"id":3, "jsonrpc":"2.0", "method": "eth_blockNumber", "params":["` +P[0] + `"]}`)


    // Поиск по всему значению второй по скорости (75 мс)
    // -------------------------------------------
    // fmt.Println(block)
    // return
    s2:=time.Now()
     
    if strings.Contains(block,"0xd55900885e66fd9326c324bc35ed1c808b7e8362e74d22e2ac24b1d29c53d99b" ){
       fmt.Println("Yesss")
    }
    f2:=time.Now()
    r2:=f2.Sub(s2)
    fmt.Println("-----------------", r2)


    // Поиск по всему части значения самый Быстрый способ  по скорости (10 мс)
    // -------------------------------------------
    s3:=time.Now()
        
    if strings.Contains(block,"c53d99b" ){
       fmt.Println("Yesss")
    }

    f3:=time.Now()
    r3:=f3.Sub(s3)
    fmt.Println("-----------------", r3)
    s1:=time.Now()

    // Прокрутка самая долгая = 756 мс
    // -------------------------------------------
    var Md Mst
    // var Pnt []Point

    errj:=json.Unmarshal([]byte(block), &Md)
    if errj!=nil{
       log.Println("Error ummarshaling.", errj.Error())
    }



    Fs:=Md["result"].(map[string]interface{})//.transactions //"].([]interface{})[0]
    // Fg:=Fs["transactions"].([]interface{})[0].(map[string]interface{})["hash"]

    Fbg:=Fs["transactions"].([]interface{})


    for _, rb:=range(Fbg){

     	  fnd:=rb.(map[string]interface{})["hash"].(string)

         fmt.Println("-----------------------\n",fnd)
         if fnd=="0xd55900885e66fd9326c324bc35ed1c808b7e8362e74d22e2ac24b1d29c53d99b" {
            fmt.Println("Yess transaction")
            break
         }
     }


     f1:=time.Now()
     r1:=f1.Sub(s1)
     fmt.Println(r1)
	


   //   bl=eth.getBlock("pending",true);
   //    bl=eth.getBlock("latest");
   //    // console.log("NEW BLOCK:",bl.hash);
   //    console.log("hash BLOCK:",bl.number );

   //    // 3528489
   //    inf=web3.eth.getBlock(bl.number);
   //    console.log("NEW BLOCK.......:", inf.hash);	
   //    console.log("NEW BLOCK.......:", inf.transactions);	

   // var number = web3.eth.getBlockTransactionCount(inf.hash);
   //    console.log("Count transaction in Blocks : ",number);	
}	










func GetLastBlock(req *http.Request)  {
    // p       := `{"id":12, "jsonrpc":"2.0", "method": "eth_getBlock", "params":["latest"]}` 
    p       := `{"id":122, "jsonrpc":"2.0", "method": "eth_blockNumber", "params":[]}` 

    
    bl:=DoHttp(req,Nodelocal, "8555", p )
    bb:=strings.Split(bl,"0x")
    bs:=strings.Replace(bb[1],`"}`,``,1)

    fmt.Println(bs)
    
}



func DoHttp(req *http.Request,Ip,Port, Com string) string {

     
    
     url     := "http://" + Ip + ":" + Port + "/"
     payload :=  strings.NewReader(Com)
     rq, _:=  http.NewRequest("POST", url, payload)

     rq.Header.Add("content-type",  "application/json")
     rq.Header.Add("cache-control", "no-cache")

     res, err := http.DefaultClient.Do(rq)
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



