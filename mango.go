/*@Author :Manasvini Banavara Suryanarayana
*@SJSU ID : 010102040
*CMPE 273 Assignment #2
*/
package main

import (
    "fmt"
    "./httprouter"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "net/url"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "os"
    "strconv"
)

type Coordinate struct {
   Lat float64 `json:"lat"`
   Lng float64 `json:"lng"`
}


type Details struct {
   ID int32  `json:"id"`
   Name string `json:"name"`
   Address string `json:"address"`
   City string `json:"city"`
   State string `json:"state"`
   Zip string `json:"zip"`
   Coordinate Coordinate `json:"coordinate"`
 }
  
  type UpdDetails struct {
   Address string `json:"address"`
   City string `json:"city"`
   State string `json:"state"`
   Zip string `json:"zip"`
   Coordinate Coordinate `json:"coordinate"`
 }
  

func getmethod(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
  var id string = p.ByName("id")
  id2, err9 :=strconv.Atoi(id)
  if err9 != nil {
        fmt.Println(err9)
    }
  fmt.Println(id)
    //database connection
    sess, err := mgo.Dial("mongodb://admin:admin@ds043714.mongolab.com:43714/cmpe273")
    if err != nil {
    fmt.Printf("Can't connect to mongo, go error %v\n", err)
    os.Exit(1)
    }
    defer sess.Close()

  sess.SetSafe(&mgo.Safe{})
 collection := sess.DB("cmpe273").C("table")


  var responsemesg Details

  err = collection.Find(bson.M{"id":id2}).One(&responsemesg)
  if err != nil {
    fmt.Printf("got an error finding a doc %v\n",err)
    os.Exit(1)
  }

    //converting response body struct to json format
   respjson, err5 := json.Marshal(responsemesg)
   if err5 != nil {
        fmt.Println(err5)
    }
     
    rw.Header().Set("Content-Type","application/json")
    rw.WriteHeader(200)
    //sending back response
    fmt.Fprintf(rw, "%s", respjson)
    
}

func postmethod(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
   fmt.Println("1")
   //creating struct to read request json 
   type ReqBody struct {
        Name string `json:"name"`
        Address string `json:"address"`
        City string `json:"city"`
        State string `json:"state"`
        Zip string `json:"zip"`
    }  
    var x ReqBody

    //fetch the body from request 
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
       fmt.Println(err)
   }
    
    //converting json body to struct of type ReqBody
    err1 := json.Unmarshal(body, &x)
    if err1 != nil {
       fmt.Println(err1)
   }
  //creating google api query
   var x1 string = x.Address+","+x.City+","+x.State
    fmt.Println("Query sent to google API : " + x1)
    var input string = url.QueryEscape(x1)
    resp, err2 := http.Get("http://maps.google.com/maps/api/geocode/json?address="+input+"&sensor=false")

    if err2 != nil {
        fmt.Println(err2)
    }
    defer resp.Body.Close()

    body1, err3 := ioutil.ReadAll(resp.Body)
    if err3 != nil {
        fmt.Println(err3)
    }

    var googmap interface{}
    err4 := json.Unmarshal(body1, &googmap)
    if err4 != nil {
        fmt.Println(err4)
    }

    res1 := googmap.(map[string]interface{})
    i := res1["results"]
    m1 := i.([]interface{})
    x8 := m1[0]
    x2 := x8.(map[string]interface{})
     y1 := x2["geometry"]
     y2 := y1.(map[string]interface{})
     y3 := y2["location"]
     fmt.Println(y3)
     y4 := y3.(map[string]interface{})
     latval := y4["lat"]
     lngval := y4["lng"]
     
   
     fmt.Println("3")
    //constructing struct for sending back response body
     cord := Coordinate{
      Lat: latval.(float64),
      Lng: lngval.(float64),
     }
    fmt.Println("4")
    var idval int32 = bson.NewObjectId().Counter()
    fmt.Println(idval)
    test := Details{
        ID:  idval,
        Name: x.Name,
        Address: x.Address,
        City: x.City,
        State: x.State,
        Zip: x.Zip,
        Coordinate: cord,
      }

      //database connection
    sess, err := mgo.Dial("mongodb://admin:admin@ds043714.mongolab.com:43714/cmpe273")
    if err != nil {
    fmt.Printf("Can't connect to mongo, go error %v\n", err)
    os.Exit(1)
    }
    defer sess.Close()

  sess.SetSafe(&mgo.Safe{})
  collection := sess.DB("cmpe273").C("table")
 
  err = collection.Insert(test)
  if err != nil {
    fmt.Printf("Can't insert document: %v\n", err)
    os.Exit(1)
  }

    //converting response body struct to json format
   respjson, err5 := json.Marshal(test)
   if err5 != nil {
        fmt.Println(err5)
    }
     
    rw.Header().Set("Content-Type","application/json")
    rw.WriteHeader(201)
    //sending back response
    fmt.Fprintf(rw, "%s", respjson)
     
}
func putmethod(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    //fmt.Fprintf(rw, "Hello, %s!\n", p.ByName("id"))
  //fetching id from request
  var id string = p.ByName("id")
  id2, err9 :=strconv.ParseInt(id,10,32)

  if err9 != nil {
        fmt.Println(err9)
    }

    //fetching data from request body

    //creating struct to read request json 
   type ReqBody struct {
        Address string `json:"address"`
        City string `json:"city"`
        State string `json:"state"`
        Zip string `json:"zip"`
    }  
    var x ReqBody

    //fetch the body from request 
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
       fmt.Println(err)
   }
    
    //converting json body to struct of type ReqBody
    err1 := json.Unmarshal(body, &x)
    if err1 != nil {
       fmt.Println(err1)
   }

   //creating google api query
   var x1 string = x.Address+","+x.City+","+x.State
    fmt.Println("Query sent to google API : " + x1)
    var input string = url.QueryEscape(x1)
    resp, err2 := http.Get("http://maps.google.com/maps/api/geocode/json?address="+input+"&sensor=false")

    if err2 != nil {
        fmt.Println(err2)
    }
    defer resp.Body.Close()

    body1, err3 := ioutil.ReadAll(resp.Body)
    if err3 != nil {
        fmt.Println(err3)
    }

    var googmap interface{}
    err4 := json.Unmarshal(body1, &googmap)
    if err4 != nil {
        fmt.Println(err4)
    }

    res1 := googmap.(map[string]interface{})
    i := res1["results"]
    m1 := i.([]interface{})
    x8 := m1[0]
    x2 := x8.(map[string]interface{})
     y1 := x2["geometry"]
     y2 := y1.(map[string]interface{})
     y3 := y2["location"]
     fmt.Println(y3)
     y4 := y3.(map[string]interface{})
     latval := y4["lat"]
     lngval := y4["lng"]
     

     //constructing struct for sending back response body
     cord := Coordinate{
      Lat: latval.(float64),
      Lng: lngval.(float64),
     }
    fmt.Println("4")
    
    test := UpdDetails{
        Address: x.Address,
        City: x.City,
        State: x.State,
        Zip: x.Zip,
        Coordinate: cord,
      }

    //database connection
    sess, err := mgo.Dial("mongodb://admin:admin@ds043714.mongolab.com:43714/cmpe273")
    if err != nil {
    fmt.Printf("Can't connect to mongo, go error %v\n", err)
    os.Exit(1)
    }
    defer sess.Close()

  sess.SetSafe(&mgo.Safe{})
  collection := sess.DB("cmpe273").C("table")

  //update the records for particular ID
  err10 := collection.Update(bson.M{"id":id2},bson.M{"$set": test})
  if err10 != nil {
    fmt.Printf("Can't update document: %v\n", err10)
    os.Exit(1)
  }

  //fetch the record to get name
  var responsemesg Details

  err11 := collection.Find(bson.M{"id":id2}).One(&responsemesg)
  if err11 != nil {
    fmt.Printf("got an error finding a doc %v\n",err11)
    os.Exit(1)
  }

    //converting response body struct to json format
   respjson, err5 := json.Marshal(responsemesg)
   if err5 != nil {
        fmt.Println(err5)
    }
     
    rw.Header().Set("Content-Type","application/json")
    rw.WriteHeader(201)
    //sending back response
    fmt.Fprintf(rw, "%s", respjson)

}
func delmethod(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    

  //fetching id from request
  var id string = p.ByName("id")
  id2, err9 :=strconv.ParseInt(id,10,32)

  if err9 != nil {
        fmt.Println(err9)
    }

  //database connection
    sess, err := mgo.Dial("mongodb://admin:admin@ds043714.mongolab.com:43714/cmpe273")
    if err != nil {
    fmt.Printf("Can't connect to mongo, go error %v\n", err)
    os.Exit(1)
    }
    defer sess.Close()

  sess.SetSafe(&mgo.Safe{})
  collection := sess.DB("cmpe273").C("table")

  //delete the particular id record

  err10 := collection.Remove(bson.M{"id":id2})
  if err10 != nil {
    fmt.Printf("Can't Remove document: %v\n", err10)
    os.Exit(1)
  }

  rw.Header().Set("Content-Type","application/json")
    rw.WriteHeader(200)

}


func main() {
    mux := httprouter.New()
    mux.GET("/locations/:id", getmethod)
    mux.POST("/locations", postmethod)
    mux.PUT("/locations/:id", putmethod)
    mux.DELETE("/locations/:id", delmethod)
    server := http.Server{
            Addr:        "0.0.0.0:8080",
            Handler: mux,
    }

    server.ListenAndServe()
}