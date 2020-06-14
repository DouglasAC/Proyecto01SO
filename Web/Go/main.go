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
	"os/exec"
)


type dataRam struct{
	Total  		int
	Consumida 	int
}

type PageVariables struct {
	Date         string
	Time         string
}

type dataCpu struct
{
	Total 	float64
}

func main(){
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/ram", obtenerRam)
	http.HandleFunc("/cpu", obtenerCpu)
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

func obtenerCpu(w http.ResponseWriter, r *http.Request) {
	//Obtener lista del porcentaje de uso del los cpus
	ps, err := exec.Command("ps", "-eo","%cpu").Output()    
    if(err != nil){
        return
    }else{ 
		//Obtener un listado 
        listaProcesos := strings.Split(string(ps),"\n") 
        var totalUso float64 = 0.0
        for i, proceso := range listaProcesos{
            if(i != 0 ){
                cpuPorCiento    := strings.Replace(proceso," ","",-1)
                cpuPC, err      := strconv.ParseFloat(cpuPorCiento, 64)
                if err == nil {
                    totalUso += cpuPC
                }
            }
        }
        jsonResponse2, errorJsonCpu  := json.Marshal(dataCpu{totalUso})
        if errorJsonCpu != nil {
            http.Error(w, errorJsonCpu.Error(), http.StatusInternalServerError)
            fmt.Fprintf(w, "jfdksljf")
            return
        }
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.WriteHeader(http.StatusOK)
        w.Write(jsonResponse2)
    }
}