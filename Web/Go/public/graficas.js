var grafRam = new Morris.Line({
    // ID of the element in which to draw the chart.
    element: 'graficaRam',
    // Chart data records -- each entry in this array corresponds to a point on
    // the chart.
    data: [
      
    ],
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
          console.log(dat.Total)
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