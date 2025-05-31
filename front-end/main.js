const description = `
     <div>
      <h2>Cos'è LudoParks?</h2>
      <p>
        ParkSense è un sistema distribuito per il monitoraggio dei parchi comunali.
      </p>
      <p>
        Attraverso dispositivi edge installati nei parchi, la piattaforma raccoglie dati ambientali  
      </p>
      <p>
        in tempo reale — come temperatura, umidità, luminosità e qualità dell’aria — e li rende
      </p>
      <p>
        disponibili tramite una dashboard semplice e interattiva.
      </p>
      <ul>
        <li><p>Visualizza lo stato dei parchi in tempo reale</p></li>
        <li><p>Rileva anomalie ambientali in modo tempestivo</p></li>
        <li><p>Ottimizza la gestione e manutenzione</p></li>
        <li><p>Offre dati trasparenti ai cittadini</p></li>
      </ul>
     </div>
  `;

const measures = [
  {label: "Temperature", icon: "./assets/icons/temperature_icon.png"},
  {label: "Humidity", icon: "./assets/icons/humidity_icon.png"},
  {label: "Brightness", icon: "./assets/icons/brightness_icon.png"},
  {label: "Air Quality", icon: "./assets/icons/air_quality_icon.png"}
];

//var sensorData = null;

rendering();
document.getElementById("sidebar").classList.toggle("open");
setInterval(
   function() {
    getParksData();
   }, 60000);

function rendering() {
  getParksData();
  rendHome();
}

function getParksData() {
  fetch("http://localhost:8080/api/hello")
    .then(response => {
      if (!response.ok) {
        // Controlla se la risposta HTTP è andata a buon fine (status 200-299)
        throw new Error(`Errore HTTP! Status: ${response.status}`);
      }
      return response.json()})
    .then(data => {
      console.log(data)
      rendSidebar(data.parks)
    }).catch(err => {
      console.log(err.text);
    });
}

function rendSidebar(parks) {
  //console.log(parks)
  const park_list = document.getElementById("park-list");
  park_list.innerHTML = "";
  const home_item = document.createElement("li");
  home_item.textContent = "Home";
  home_item.onclick = () => {
    rendHome();
  };
  park_list.appendChild(home_item);

  parks.forEach(park => {
    //console.log(park);
    const item = document.createElement("li");
    item.textContent = park.name;
    item.onclick = () => {
      rendParkPhoto(park.name)
      rendTitle(park.name)
      renderComponentsForPark(park);
      renderOldDataForPark(park.olddata);
      renderTimestamp(park.timestamp)
      
      document.getElementById("description").innerHTML = "";
      //console.log("park.name")
    };
    park_list.appendChild(item);
  });

  
}

function rendHome() {
  //immagine sfondo
  document.getElementById("park-photo").src = "./assets/logo/logo1.png";
  document.getElementById("park-photo").alt = "LudoParks logo";
  //titolo
  document.getElementById("park-name").textContent = "";
  //lista delle misure
  document.getElementById("component-list").innerHTML = "";
  //vecchie misure meteo
  document.getElementById("old-data-list").innerHTML = "";
  //descrizione
  document.getElementById("description").innerHTML = description;
  //timestamp
  document.getElementById("last-update").innerHTML = "";
}

function rendParkPhoto(name) {
  lowname = name.replace(/ /g, "_").toLowerCase() + ".jpg";
  const parkPhoto = document.getElementById("park-photo");
  parkPhoto.src = "./assets/photos/" + lowname;
  parkPhoto.alt = "foto del parco " + name;
}

function rendTitle(name) {
  document.getElementById("park-name").textContent = `${name}`;
}

function renderComponentsForPark(park) {
  const component_list = document.getElementById("component-list");
  component_list.innerHTML = "";
  
  // componiamo gli elementi per mostrare le misure
  //Temperatura
  const divT = document.createElement("div");
  divT.className = "card";
  divT.innerHTML = `<div><img src="${measures[0].icon}"></div><div><strong>${measures[0].label}</strong></div><p></p><div>${getCelsius(park.temperature)}</div>`;
  component_list.appendChild(divT)

  //Umidità
  const divH = document.createElement("div");
  divH.className = "card";
  divH.innerHTML = `<div><img src="${measures[1].icon}"></div><div><strong>${measures[1].label}</strong></div><p></p><div>${getPercentage(park.humidity)}</div>`;
  component_list.appendChild(divH)

  //Luminosità
  const divB = document.createElement("div");
  divB.className = "card";
  divB.innerHTML = `<div><img src="${measures[2].icon}"></div><div><strong>${measures[2].label}</strong></div><p></p><div>${getLux(park.brightness)}</div>`;
  component_list.appendChild(divB)

  //Umidità
  const divAQ = document.createElement("div");
  divAQ.className = "card";
  divAQ.innerHTML = `<div><img src="${measures[3].icon}"></div><div><strong>${measures[3].label}</strong></div><p></p><div>${getPMI(park.airquality)}</div>`;
  component_list.appendChild(divAQ)
}

function getCelsius(temp) {
  return temp + " °C"
}

function getPercentage(hum) {
  return hum + "%"
}

function getLux(bright) {
  return bright + " lux"
}

function getPMI(aq) {
  return aq + " pmi"
}

function renderOldDataForPark(els) {
  const old_list = document.getElementById("old-data-list");
  old_list.innerHTML = "";

  els.forEach(el => {
    const div = document.createElement("div");
    div.className = "card";
    div.innerHTML = `<div><strong>${getDate(el.date)}</strong></div><div><img src="${getIcon(el.icon)}"></div><div><div style="color:blue">${el.min}</div><div style="color:red">${el.max}</div></div>`;
    old_list.appendChild(div);
  })
}

function getDate(date) {
  return date[8]+date[9]+"/"+date[5]+date[6];
}

function getIcon(icon) {
  if (icon == 1) {
    return "./assets/icons/sunny_icon.png";
  } else if (icon == 2) {
    return "./assets/icons/cloudy_icon.png";
  } else if (icon == 3) {
    return "./assets/icons/rainy_icon.png";
  } else {
    return "./assets/icons/snowy_icon.png";
  }
}

function renderTimestamp(timestamp) {
  document.getElementById("last-update").innerHTML = "last update: " + timestamp;
}
