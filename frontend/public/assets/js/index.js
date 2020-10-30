class CarType{
    constructor(databaseEntry) {
        this.id = parseInt(databaseEntry.id);
        this.name = databaseEntry.name;
    }
}

class Car{
    constructor(databaseEntry) {
        this.id = parseInt(databaseEntry.id);
        this.carType = new CarType(databaseEntry.carType);
    }
}

class Node{
    constructor(databaseEntry) {
        this.id = parseInt(databaseEntry.id);
        this.name = databaseEntry.name;
        let carTypes = []
        if (databaseEntry.availableCarTypes !== undefined){
            databaseEntry.availableCarTypes.forEach(function(carType){
                carTypes.push(new CarType(carType));
            });
        }
        this.carTypes = carTypes;
    }
    hasCarType(id){
        for(let i = 0; i < this.carTypes.length; i++){
            if(this.carTypes[i].id === id) {
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
        this.date = new Date(databaseEntry.date);
        this.departure = new Node(databaseEntry.departureNode);
        this.arrival = new Node(databaseEntry.arrivalNode);
        this.car = new Car(databaseEntry.car);
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

function addNodesToSelect(select, firstInit = false){
    select.innerHTML = "";
    if(select === nodeArrivalSelect){
        addOptionToSelect(select, "Aucun noeud d'arrivée sélectionné", 0, true, null)
    }else{
        addOptionToSelect(select, "Aucun noeud de départ sélectionné", 0, true, null)
    }
    select.disabled = false;

    for (let i = 0; i < listNodes.length; i++){
        const node = listNodes[i];
        if(node.hasCarType(carTypeIdSelected)){
            if(i === 0 && firstInit){
                if(select === nodeArrivalSelect){
                    nodeArrivalIdSelected = node.id;
                }else{
                    nodeDepartureIdSelected = node.id;
                }
            }
            if(select === nodeDepartureSelect && node.id !== nodeArrivalIdSelected){
                addOptionToSelect(select, node.name, node.id, false, nodeDepartureIdSelected);
            }

            if(select === nodeArrivalSelect && node.id !== nodeDepartureIdSelected) {
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
    loadIdMember.open('GET', 'http://localhost/car-booking/getAllNodes', true);
    loadIdMember.addEventListener('readystatechange', function() {
        if(this.readyState === 4 && this.status === 200) {
            let response = JSON.parse(this.responseText);
            response.forEach(function(node){
               listNodes.push(new Node(node));
                addOptionToSelect(nodeDepartureSelect, node.name, node.id, false, null);
                addOptionToSelect(nodeArrivalSelect, node.name, node.id, false, null);
            });
        }
    });
    loadIdMember.send(null);
}

function loadCarTypes(){
    let loadIdMember = new XMLHttpRequest();
    loadIdMember.open('GET', 'http://localhost/car-booking/getAllCarTypes', true);
    loadIdMember.addEventListener('readystatechange', function() {
        if(this.readyState === 4 && this.status === 200) {
            let response = JSON.parse(this.responseText);
            response.forEach(function(carType){
                let carTypeObject = new CarType(carType)
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
        buttonSubmit.disabled = false;
        addNodesToSelect(nodeDepartureSelect, true)
        addNodesToSelect(nodeArrivalSelect, true)
    });
    nodeDepartureSelect.addEventListener("change", function(e){
        nodeDepartureIdSelected = parseInt(e.target.value);
        addNodesToSelect(nodeArrivalSelect);
        addNodesToSelect(nodeDepartureSelect);
    });
    nodeArrivalSelect.addEventListener("change", function(e){
        nodeArrivalIdSelected = parseInt(e.target.value);
        addNodesToSelect(nodeDepartureSelect);
        addNodesToSelect(nodeArrivalSelect);
    });
    document.getElementById("date-departure").addEventListener("input", function(e){
        dateTimeDeparture = e.target.value;
    });
    document.getElementById("middle-form").addEventListener("submit", handleFormSubmit);
    loadNodes();
    loadCarTypes();
})()

function handleFormSubmit(e){
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
    fetchOffers.open('POST', 'http://localhost/booking-process/offers', true);
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
    fetchOffers.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
    fetchOffers.send(JSON.stringify({

    }));
    // setTimeout(function(){
    //     displayOffers([
    //         new Offer({
    //             id: "g15ad65fed5",
    //             price: "5.0",
    //             date: new Date().toISOString(),
    //             car: {
    //                 id: 1,
    //                 carType: {
    //                     name: "Liquid",
    //                     id: 2
    //                 }
    //             },
    //             departureNode: {
    //                 name: "Nice",
    //                 id: 1
    //             },
    //             arrivalNode: {
    //                 name: "Marseille",
    //                 id: 2
    //             }
    //         })
    //     ]);
    // }, 50);
}

function displayOffer(offer){
    console.log(offer)
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
    node.appendChild(document.createTextNode(offer.departure.name));
    nodes.appendChild(node);

    node = document.createElement("div");
    node.classList.add("offer-node");
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
        if(this.readyState === 4 && this.status === 201) {
            removeLoader();
            let response = JSON.parse(this.responseText);
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
