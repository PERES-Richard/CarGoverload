class CarType{
    constructor(id, name) {
        this.id = id;
        this.name = name;
    }
}

class Car{
    constructor(databaseEntry) {
        this.id = parseInt(databaseEntry.id);
        this.carType = new CarType(databaseEntry.carType.id, databaseEntry.carType.name);
    }
}

class Node{
    constructor(id, name, types) {
        this.id = id;
        this.name = name;
        this.carTypes = types;
    }
    hasCarType(id){
        for(let i = 0; i < this.carTypes.length; i++){
            if(this.carTypes[i] === id) {
                return true;
            }
        }
        return false;
    }
}

class Offer{
    constructor(databaseEntry) {
        this.id = databaseEntry.id;
        this.price = databaseEntry.price;
        this.date = new Date(databaseEntry.beginBookedDate);
        this.departure = new Node(databaseEntry.departureNode.id, databaseEntry.departureNode.name, []);
        this.arrival = new Node(databaseEntry.arrivalNode.id, databaseEntry.arrivalNode.name, []);
        this.car = new Car(databaseEntry.car);
        this.duration = databaseEntry.duration;
        this.departureTime = this.date.getHours() + ':' + this.date.getMinutes();
        this.dateArrival = new Date(databaseEntry.endingBookedDate);
        this.arrivalTime = this.dateArrival.getHours() + ':' + this.dateArrival.getMinutes();
        console.log(this)
    }
}

listNodes = []

let carTypeIdSelected = -1;
let nodeArrivalIdSelected = -1;
let nodeDepartureIdSelected = -1;
let supplierInput = null;
let carTypeIdSelect = null;
let nodeArrivalSelect = null;
let nodeDepartureSelect = null;
let dateTimeDeparture = null;
let loadingBigContainer = null;
let mainContainer = null;

function addNodesToSelect(select){
    select.innerHTML = "";
    if(select === nodeArrivalSelect){
        addOptionToSelect(select, "Aucun noeud d'arrivée sélectionné", 0, false, null)
    }else{
        addOptionToSelect(select, "Aucun noeud de départ sélectionné", 0, false, null)
    }
    select.disabled = false;

    for (let i = 0; i < listNodes.length; i++){
        const node = listNodes[i];
        if(node.hasCarType(carTypeIdSelected)){
            if(select === nodeDepartureSelect){
                addOptionToSelect(select, node.name, node.id, false, nodeDepartureIdSelected);
            }

            if(select === nodeArrivalSelect) {
                addOptionToSelect(select, node.name, node.id, false, nodeArrivalIdSelected);
            }
        }
    }
}

function addOptionToSelect(select, name, value, disabled, selectedValue){
    let option = document.createElement("option");
    option.value = value;
    option.disabled = disabled
    if (value === selectedValue)
        option.selected = true
    option.appendChild(document.createTextNode(name))
    select.appendChild(option);
}

function loadNodes(){
    let loadIdMember = new XMLHttpRequest();
    loadIdMember.open('GET', 'http://localhost/car-location/findAllNodes', true);
    loadIdMember.addEventListener('readystatechange', function() {
        if(this.readyState === 4 && this.status === 200) {
            let response = JSON.parse(this.responseText);
            response.forEach(function(node){
                listNodes.push(new Node(node.id, node.name, node.types));
                addOptionToSelect(nodeDepartureSelect, node.name, node.id, false, null);
                addOptionToSelect(nodeArrivalSelect, node.name, node.id, false, null);
            });
        }
    });
    loadIdMember.send(null);
}

function loadCarTypes(){
    let loadIdMember = new XMLHttpRequest();
    loadIdMember.open('GET', 'http://localhost/car-location/findAllCarTypes', true);
    loadIdMember.addEventListener('readystatechange', function() {
        if(this.readyState === 4 && this.status === 200) {
            let response = JSON.parse(this.responseText);
            response.forEach(function(carType){
                let carTypeObject = new CarType(carType.id, carType.name)
                addOptionToSelect(carTypeIdSelect, carTypeObject.name, carTypeObject.id, false, null);
            });
        }
    });
    loadIdMember.send(null);
}

(function(){
    loadingBigContainer = document.getElementById("loading-big-container");
    nodeDepartureSelect = document.getElementById("node-departure-select");
    nodeArrivalSelect = document.getElementById("node-arrival-select");
    carTypeIdSelect = document.getElementById("car-type-select");
    supplierInput = document.getElementById("supplier");
    mainContainer = document.getElementById("main-container");

    let buttonSubmit = document.getElementById("middle-form-submit");

    carTypeIdSelect.addEventListener("change", function(e){ // when selecting an other value
        carTypeIdSelected = parseInt(e.target.value);
        nodeArrivalIdSelected = -1
        nodeDepartureIdSelected = -1;
        buttonSubmit.disabled = false;
        addNodesToSelect(nodeDepartureSelect)
        addNodesToSelect(nodeArrivalSelect)
    });
    nodeDepartureSelect.addEventListener("change", function(e){
        nodeDepartureIdSelected = parseInt(e.target.value);
        addNodesToSelect(nodeArrivalSelect);
    });
    nodeArrivalSelect.addEventListener("change", function(e){
        nodeArrivalIdSelected = parseInt(e.target.value);
        addNodesToSelect(nodeDepartureSelect);
    });
    document.getElementById("date-departure").addEventListener("input", function(e){
        dateTimeDeparture = e.target.value;
    });
    document.getElementById("middle-form").addEventListener("submit", handleFormSubmit);
    loadCarTypes();
    loadNodes();
})()

function handleFormSubmit(e){
    removeLoader()
    e.preventDefault();
    if (carTypeIdSelected === -1){
        alert("Vous devez choisir un type de wagon");
        return;
    }
    if (nodeDepartureIdSelected === -1){
        alert("Vous devez choisir un noeud de départ");
        return;
    }
    if (nodeArrivalIdSelected === -1){
        alert("Vous devez choisir un noeud d'arrivé");
        return;
    }
    if (nodeArrivalIdSelected === nodeDepartureIdSelected){
        alert("Vous devez choisir des nodes de départ et d'arrivée différents !");
        return;
    }
    if (dateTimeDeparture === null){
        alert("Vous devez choisir une date et une heure de départ");
        return;
    }
    if(supplierInput.value === null || supplierInput.value.length === 0){
        alert("Le nom du fournisseur doit être renseigné");
        return;
    }
    launchSearch();
}

function addLoader(){
    let loadingContainer = document.getElementById("loading-container");
    if(loadingContainer !== null)
        return;
    loadingContainer = document.createElement('div');
    loadingContainer.id = 'loading-container';
    const loader = document.createElement('div');
    loader.classList.add('loader');
    loadingContainer.appendChild(loader);
    loadingBigContainer.appendChild(loadingContainer);
}

function removeLoader(){
    loadingBigContainer.innerHTML = "";
}

function launchSearch(){
    mainContainer.innerHTML = "";
    addLoader();
    let fetchOffers = new XMLHttpRequest();
    fetchOffers.open('GET', 'http://localhost/booking-process/offers?' +
        'supplier=' + supplierInput.value +
        '&carTypeId=' + carTypeIdSelected +
        '&arrivalNodeId=' + nodeArrivalIdSelected +
        '&departureNodeId=' + nodeDepartureIdSelected +
        '&dateTimeDeparture=' + (new Date(dateTimeDeparture)).toISOString(), true);
    fetchOffers.addEventListener('readystatechange', function() {
        if(this.readyState === 4 && this.status === 201) {
            removeLoader();
            let response = JSON.parse(this.responseText);
            removeLoader();
            response.forEach(function(offer){
                displayOffer(new Offer(offer))
            });
        }
    });
    fetchOffers.send(null);
}

function displayOffer(offer){
    let container = document.createElement("div");
    container.classList.add("offer");
    container.title = "Réserver";
    let information = document.createElement("div");
    information.classList.add("offer-information");
    information.classList.add("offer-container");

    let nodes = document.createElement("div");
    nodes.classList.add("offer-nodes-container");

    let node = document.createElement("div");
    node.classList.add("offer-node");

    let time = document.createElement("time");
    time.classList.add("offer-time");
    time.appendChild(document.createTextNode(offer.departureTime + ' - '));
    node.appendChild(time);
    node.appendChild(document.createTextNode(offer.departure.name));
    nodes.appendChild(node);

    node = document.createElement("div");
    node.classList.add("offer-node");
    time = document.createElement("time");
    time.classList.add("offer-time");
    time.appendChild(document.createTextNode(offer.arrivalTime + ' - '));
    node.appendChild(time);
    node.appendChild(document.createTextNode(offer.arrival.name));
    nodes.appendChild(node);

    let dateContainer = document.createElement("div");
    dateContainer.classList.add("offer-date");
    dateContainer.appendChild(document.createTextNode(offer.date.toDateString()));
    nodes.appendChild(dateContainer);

    information.appendChild(nodes);

    let carContainer = document.createElement("div");
    carContainer.classList.add("offer-car");

    let carIcon = document.createElement("i");
    carIcon.classList.add("offer-car-icon");
    carIcon.classList.add("fas");
    carIcon.classList.add("fa-train");
    carContainer.appendChild(carIcon);

    let carInfoContainer = document.createElement("div");
    carInfoContainer.classList.add("offer-car-information")

    let carText = document.createElement("div");
    carText.appendChild(document.createTextNode("Wagon n°" + offer.car.id))
    carInfoContainer.appendChild(carText);

    let carType = document.createElement("div");
    carType.appendChild(document.createTextNode(offer.car.carType.name))
    carInfoContainer.appendChild(carType);

    carContainer.appendChild(carInfoContainer)
    information.appendChild(carContainer)

    container.appendChild(information);
    let price = document.createElement("div");
    price.classList.add("offer-price");
    price.classList.add("offer-container");
    price.appendChild(document.createTextNode(offer.price + "€"));
    container.appendChild(price);

    container.addEventListener("click", function(){
        const res = confirm("Procéder à la réservation et au paiement ?")
        if(res){
            book(offer, container);
        }else{

        }
    });

    mainContainer.appendChild(container);
}

function book(offer, view){
    let bookOffer = new XMLHttpRequest();
    bookOffer.open('POST', 'http://localhost/booking-process/offers/payment', true);
    bookOffer.addEventListener('readystatechange', function() {
        if(this.readyState === 4 && this.status === 200) {
            removeLoader();
            let response = JSON.parse(this.responseText);
            console.log(response);
            console.log(view)
            view.remove();
        }
    });
    bookOffer.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
    bookOffer.send(JSON.stringify(
        {
            offerId: offer.id,
            supplier: supplierInput.value
        }
    ));
}
