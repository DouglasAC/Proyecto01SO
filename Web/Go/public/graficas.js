var grafRam = new Morris.Line({
    // ID of the element in which to draw the chart.
    element: 'graficaRam',
    // Chart data records -- each entry in this array corresponds to a point on
    // the chart.
    data: [],
    // The name of the data record attribute that contains x-values.
    xkey: 'num',
    // A list of names of data record attributes that contain y-values.
    ykeys: ['value'],
    // Labels for the ykeys -- will be displayed when you hover over the
    // chart.
    labels: ['Value'],
    resize: true
  });
  
  let dataRam = [];
  let contRam = 1;

  function leerRam()
  {
    $.ajax({
        async: false,
        contentType: 'application/json;  charset=utf-8',
        type: "GET",
        url: "/ram",
        data: JSON.stringify({
            
        }),
        si: false,
        success: function (data, textStatus, jqXHR) {
          let datos = JSON.stringify(data);
          let dat = JSON.parse(datos)
  
          document.getElementById("totalRam").innerHTML = "Total de memoria Ram: "+ dat.Total;
          document.getElementById("memoriaC").innerHTML = "Total de memoria consumida: "+(dat.Total-dat.Consumida);
          document.getElementById("memoriaP").innerHTML = "Porcentaje de Consumo: "+(((dat.Total-dat.Consumida)*100)/dat.Total)+"%";
          let nuevo = {num: ''+contRam++ +'', value:  (dat.Total-dat.Consumida)};
          dataRam.push(nuevo);
          grafRam.setData(dataRam);
          if(dataRam.length>15)
          {
            dataRam.shift()
          }
        },
        error: function (jqXHR, textStatus, errorThrown) {
            console.log("Error")
        }

    });
  }

  setInterval('leerRam()',1000);

  function cambiar(nombre)
  {
    if(nombre=="ram")
    {
      document.getElementById("ram").style.display="block";
      document.getElementById("cpu").style.display="none";
      document.getElementById("procesos").style.display="none";
      document.getElementById("btnCpu").className="tabslink";
      document.getElementById("btnProcesos").className="tabslink";
      document.getElementById("btnRam").className="tabslink active";
     
    }else if(nombre=="cpu")
    {
      document.getElementById("ram").style.display="none";
      document.getElementById("cpu").style.display="block";
      document.getElementById("procesos").style.display="none";
      document.getElementById("btnRam").className="tabslink";
      document.getElementById("btnProcesos").className="tabslink";
      document.getElementById("btnCpu").className="tabslink active";
     
    }else{
      document.getElementById("ram").style.display="none";
      document.getElementById("cpu").style.display="none";
      document.getElementById("procesos").style.display="block";
      document.getElementById("btnRam").className="tabslink";
      document.getElementById("btnCpu").className="tabslink";
      document.getElementById("btnProcesos").className="tabslink active";
    }
    
  }

  var grafCpu = new Morris.Line({
    // ID of the element in which to draw the chart.
    element: 'graficaCpu',
    // Chart data records -- each entry in this array corresponds to a point on
    // the chart.
    data: [],
    // The name of the data record attribute that contains x-values.
    xkey: 'num',
    // A list of names of data record attributes that contain y-values.
    ykeys: ['cpu'],
    // Labels for the ykeys -- will be displayed when you hover over the
    // chart.
    labels: ['Uso Cpu'],
    resize: true
  });

  let dataCpu = [];
  let contCpu = 1;

  function leerCpu()
  {
    $.ajax({
        async: false,
        contentType: 'application/json;  charset=utf-8',
        type: "GET",
        url: "/cpu",
        data: JSON.stringify({
            
        }),
        si: false,
        success: function (data, textStatus, jqXHR) {
          let datos = JSON.stringify(data);
          let dat = JSON.parse(datos)
          
          document.getElementById("cpuP").innerHTML = "Porcentaje de Consumo de los cpus: "+ dat.Total+ " %";
          
          let nuevo = {num: ''+contCpu++ +'', cpu: dat.Total};
          dataCpu.push(nuevo);
          grafCpu.setData(dataCpu);
          if(dataCpu.length>15)
          {
            dataCpu.shift()
          }
        },
        error: function (jqXHR, textStatus, errorThrown) {
            console.log("Error")
        }

    });
  }

  setInterval('leerCpu()',1000);

  function leerProcesos()
  {
    $.ajax({
      async: false,
      contentType: 'application/json;  charset=utf-8',
      type: "GET",
      url: "/procesos",
      data: JSON.stringify({
          
      }),
      si: false,
      success: function (data, textStatus, jqXHR) {
        let datos = JSON.stringify(data);
        let dat = JSON.parse(datos)
        
        document.getElementById("Running").innerHTML = "Running: "+ dat.Running;
        document.getElementById("Interruptible").innerHTML = "Interruptible: "+ dat.Interruptible;
        document.getElementById("Zombie").innerHTML = "Zombie: "+ dat.Zombie;
        document.getElementById("Stopped").innerHTML = "Stopped: "+ dat.Stopped;
        document.getElementById("Swapping").innerHTML = "Swapping: "+ dat.Swapping;
        document.getElementById("Total").innerHTML = "Total: "+ dat.Total;
        
        document.getElementById("tabla").innerHTML = dat.Tabla;
      },
      error: function (jqXHR, textStatus, errorThrown) {
          console.log("Error")
      }

    });
  }
  setInterval('leerProcesos()',2000);

  function parar(codigo)
  {
    console.log(parseInt(codigo))
    $.ajax({
      async: false,
      contentType: 'application/json;  charset=utf-8',
      type: "POST",
      url: "/kill",
      data: JSON.stringify({
          Numero: parseInt(codigo)
      }),
      success: function (data, textStatus, jqXHR) {
        leerProcesos()
      },
      error: function (jqXHR, textStatus, errorThrown) {
          console.log("Error")
      }

    });
  }