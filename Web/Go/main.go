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

type dataCpu struct{
	Total 	float64
}

type dataProcoso struct{
    Running   			int
    Interruptible    	int
    Uninterruptible     int
    Zombie    			int
    Stopped         	int
    Swapping            int
    Total          		int
    Tabla          		string
}

type killProc struct{
	Numero 		int
}

func main(){
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/ram", obtenerRam)
	http.HandleFunc("/cpu", obtenerCpu)
	http.HandleFunc("/procesos", obtenerProcesos)
	http.HandleFunc("/kill", killProceso)
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
        josnRespuestaCpu, errorJsonCpu  := json.Marshal(dataCpu{totalUso})
        if errorJsonCpu != nil {
            http.Error(w, errorJsonCpu.Error(), http.StatusInternalServerError)
            fmt.Fprintf(w, "jfdksljf")
            return
        }
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.WriteHeader(http.StatusOK)
        w.Write(josnRespuestaCpu)
    }
}

func obtenerProcesos(w http.ResponseWriter, r *http.Request) {
	//Leer archivo cpu de la parte de modulos
	archivoProc , err := ioutil.ReadFile("/proc/cpu_201503935") 
    if err != nil {
        return
	}
	//Obtner la lista
	cadena := string(archivoProc)
	listaProcesos := strings.Split(cadena,"\n") 
	tabla := "<table class=\"table table-bordered\" >\n"
	tabla += "<thead class=\"thead-dark\">\n"
    tabla += "<tr>\n"
    tabla += "<th scope=\"col\">#</th>\n"
    tabla += "<th scope=\"col\">PDI</th>\n"
    tabla += "<th scope=\"col\">Nombre</th>\n"
    tabla += "<th scope=\"col\">Usuario</th>\n"
	tabla += "<th scope=\"col\">Estado</th>\n"
	tabla += "<th scope=\"col\">Ram</th>\n"
	tabla += "<th scope=\"col\">Kill</th>\n"
    tabla += "</tr>\n"
    tabla += "</thead>\n"
	tabla += "<tbody>\n"
	x := 1
	Running := 0
    Interruptible := 0
    Uninterruptible := 0
    Zombie := 0
    Stopped := 0
    Swampping := 0
	for i, proceso := range listaProcesos{
		if(i > 2){
			fila := strings.Split(string(proceso), ",")
			
			if(strings.HasPrefix(fila[0],"Hijo")){
				//Hijo de systemd,1, PID: 362, Nombre: systemd-journal, User: 0, Estado: Interruptible
				tabla += "<tr class=\"bg-secondary text-white\">\n"
				tabla += "<th scope=\"row\">" + strconv.Itoa(x) + "</td>\n"
				pdi := strings.Split(fila[2],":")
				tabla += "<th scope=\"row\">" + pdi[1] + "</td>\n"
				nombre := strings.Split(fila[3], ":")
				tabla += "<th scope=\"row\">" + nombre[1] + "</td>\n"
				us := strings.Split(fila[4],":")
				usuario := getUsuario(string(us[1]))
				tabla += "<th scope=\"row\">" + usuario + "</td>\n"
				estado := strings.Split(fila[5], ":")
				tabla += "<th scope=\"row\">" + estado[1] + "</td>\n"
				if(estado[1] == " Running"){
					Running++
				}else if(estado[1] == " Interruptible"){
					Interruptible++
				}else if(estado[1] == " Uninterruptible"){
					Uninterruptible++
				}else if(estado[1] == " Zombie"){
					Zombie++
				} else if(estado[1] == " Stopped"){
					Stopped++
				}else if(estado[1] == " Swapping"){
					Swampping++
				}
				tabla += "<th scope=\"row\">-</td>\n";
				tabla += "<th scope=\"col\"><button type=\"button\" class=\"btn btn-danger\" onclick=\"parar('" + pdi[1] + "')\">Kill</button>\n</th>\n"

			}else{
				//PID: 1, Nombre: systemd, User: 0, Estado: Interruptible
				tabla += "<tr class=\"bg-primary text-white\">\n"
				tabla += "<th scope=\"row\">" + strconv.Itoa(x) + "</td>\n"
				pdi := strings.Split(fila[0],":")
				tabla += "<th scope=\"row\">" + pdi[1] + "</td>\n"
				nombre := strings.Split(fila[1], ":")
				tabla += "<th scope=\"row\">" + nombre[1] + "</td>\n"
				us := strings.Split(fila[2],":")
				usuario := getUsuario(string(us[1]))
				tabla += "<th scope=\"row\">" + usuario + "</td>\n"
				estado := strings.Split(fila[3], ":")
				tabla += "<th scope=\"row\">" + estado[1] + "</td>\n"
				if(estado[1] == " Running"){
					Running++
				}else if(estado[1] == " Interruptible"){
					Interruptible++
				}else if(estado[1] == " Uninterruptible"){
					Uninterruptible++
				}else if(estado[1] == " Zombie"){
					Zombie++
				} else if(estado[1] == " Stopped"){
					Stopped++
				}else if(estado[1] == " Swapping"){
					Swampping++
				}
				tabla += "<th scope=\"row\">-</td>\n";
				tabla += "<th scope=\"col\"><button type=\"button\" class=\"btn btn-danger\" onclick=\"parar('" + pdi[1] + "')\">Kill</button></th>\n"

			}
			x++
		}
	}
	josnRespuestaCpu, errorJsonCpu  := json.Marshal(dataProcoso{Running,Interruptible,Uninterruptible,Zombie,Stopped,Swampping,x-1,tabla})
	if errorJsonCpu != nil {
		http.Error(w, errorJsonCpu.Error(), http.StatusInternalServerError)
		fmt.Fprintf(w, "jfdksljf")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(josnRespuestaCpu)
	return
}

func getUsuario(user string) string{
	if( user ==" 0"){
		return "root"
	} else if( user ==" 1"){
		return "daemon"
	} else if( user ==" 2"){
		return "bin"
	} else if( user ==" 3"){
		return "sys"
	} else if( user ==" 4"){
		return "sync"
	} else if( user ==" 5"){
		return "games"
	} else if( user ==" 6"){
		return "man"
	} else if( user ==" 7"){
		return "lp"
	} else if( user ==" 8"){
		return "mail"
	} else if( user ==" 9"){
		return "news"
	} else if( user ==" 10"){
		return "uucp"
	} else if( user ==" 13"){
		return "proxy"
	} else if( user ==" 33"){
		return "www-data"
	} else if( user ==" 34"){
		return "backup"
	} else if( user ==" 38"){
		return "list"
	} else if( user ==" 39"){
		return "irc"
	} else if( user ==" 41"){
		return "gnats"
	} else if( user ==" 65534"){
		return "nobody"
	} else if( user ==" 100"){
		return "systemd-network"
	} else if( user ==" 101"){
		return "systemd-resolve"
	} else if( user ==" 102"){
		return "syslog"
	} else if( user ==" 103"){
		return "messagebus"
	} else if( user ==" 104"){
		return "_apt"
	} else if( user ==" 105"){
		return "uuidd"
	} else if( user ==" 106"){
		return "avahi-autoipd"
	} else if( user ==" 107"){
		return "usbmux"
	} else if( user ==" 108"){
		return "dnsmasq"
	} else if( user ==" 109"){
		return "rtkit"
	} else if( user ==" 110"){
		return "cups-pk-helper"
	} else if( user ==" 111"){
		return "speech-dispatcher"
	} else if( user ==" 112"){
		return "whoopsie"
	} else if( user ==" 113"){
		return "kernoops"
	} else if( user ==" 114"){
		return "saned"
	} else if( user ==" 115"){
		return "pulse"
	} else if( user ==" 116"){
		return "avahi"
	} else if( user ==" 117"){
		return "colord"
	} else if( user ==" 118"){
		return "hplip"
	} else if( user ==" 119"){
		return "geoclue"
	} else if( user ==" 120"){
		return "gnome-initial-setup"
	} else if( user ==" 121"){
		return "gdm"
	} else if( user ==" 1000"){
		return "douglas"
	} 
	return user
}

func killProceso(w http.ResponseWriter, r *http.Request) {  
	//Obtner json 
    data, errorkill := ioutil.ReadAll(r.Body)
    defer r.Body.Close()

    if errorkill != nil{
		fmt.Println("ERROR 01")
        http.Error(w,errorkill.Error(),500)
        return
	}
	//obtner data
	var procesoKill killProc
	
    errorUnma := json.Unmarshal(data, &procesoKill)

    if errorUnma != nil{
		fmt.Println("ERROR 02")
        http.Error(w,errorUnma.Error(),500)
        return
	}
	//obtener numero del proceso
    proceso := strconv.Itoa(procesoKill.Numero)
	
	//ejecutar comando
    errorComando := exec.Command("kill","-9", proceso).Run()    
    if errorComando != nil {
		fmt.Println("ERROR 03")
        http.Error(w,errorComando.Error(),500)
        return
    }

    return
}