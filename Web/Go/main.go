package main

import (
	"fmt"
    "net/http"
    "io/ioutil"
    "strings"
    "strconv"
	"encoding/json"
	"time"
	"html/template"
)


type dataRam struct{
	Total  		int
	Consumida 	int
}

type PageVariables struct {
	Date         string
	Time         string
}

func main(){
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/ram", obtenerRam)
	fmt.Println("puerto 3000")
	http.ListenAndServe(":3000", nil)
}

func HomePage(w http.ResponseWriter, r *http.Request){

    now := time.Now() // find the time right now
    HomePageVars := PageVariables{ //store the date and time in a struct
      Date: now.Format("02-01-2006"),
      Time: now.Format("15:04:05"),
    }

    t, err := template.ParseFiles("index.html") //parse the html file homepage.html
    if err != nil { // if there is an error
  	  fmt.Print("template parsing error: ", err) // log it
  	}
    err = t.Execute(w, HomePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
    if err != nil { // if there is an error
  	  fmt.Print("template executing error: ", err) //log it
  	}
}

func obtenerRam(w http.ResponseWriter, r *http.Request) {    
	//Leer el archivo meminfo
	archivoMeminfo, err := ioutil.ReadFile("/proc/meminfo") 
    if err != nil {
        return
	} 
	//Pasar a string    
	cadenaMeminfo := string(archivoMeminfo)
	//Arreglo de la info de meminfo
    listaInfo := strings.Split(string(cadenaMeminfo),"\n")
	//Memoria total en la popsicion 0 apartir del caracter 10 al 24
	memoriaTotal := strings.Replace((listaInfo[0])[10:24]," ","",-1)
	//Memoria habilitada en la posicion 2 apartir del caracter 14 al 24
	memoriaHabilitada := strings.Replace((listaInfo[2])[14:24]," ","",-1)
	//Pasar a entero la memiria 
    memoriaTotalKb, err1:= strconv.Atoi(memoriaTotal)
    memoriaHabilitadaKb, err2 := strconv.Atoi(memoriaHabilitada) 
    if err1 == nil && err2 == nil{
      	memoriaTotalMb := memoriaTotalKb / 1024
        memoriaHabilitadaMb := memoriaHabilitadaKb / 1024   

		estrucuraRam := dataRam { Total: memoriaTotalMb, Consumida: memoriaHabilitadaMb }
		//Pasar a jason encode
        jsonResponse, errorjson   := json.Marshal(estrucuraRam)
        if errorjson != nil {
            http.Error(w, errorjson.Error(), http.StatusInternalServerError)
            return
		}
		//Responder 
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.WriteHeader(http.StatusOK)
        w.Write(jsonResponse)
    }else{
       return
    }    
}