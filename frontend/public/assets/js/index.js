class CarType{
    constructor(databaseEntry) {
        this.id = parseInt(databaseEntry.id);
        this.name = databaseEntry.name;
    }
}

class Node{
    constructor(databaseEntry) {
        this.id = parseInt(databaseEntry.id);
        this.name = databaseEntry.name;
        let carTypes = []
        databaseEntry.availableCarTypes.forEach(function(carType){
            carTypes.push(new CarType(carType));
        })
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
    addLoader();
    let loadIdMember = new XMLHttpRequest();
    loadIdMember.open('POST', 'http://localhost/booking-process/offers', true);
    loadIdMember.addEventListener('readystatechange', function() {
        if(this.readyState === 4 && this.status === 200) {
            removeLoader();
            let response = JSON.parse(this.responseText);
            console.log(response)
        }
    });
    loadIdMember.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
    loadIdMember.send('departureNodeId=' + nodeDepartureIdSelected +
        "&arrivalNodeId=" + nodeArrivalIdSelected +
        "&dateTimeDeparture=" + dateTimeDeparture +
        "&supplier=" + supplierInput.value +
        "&carTypeId=" + carTypeIdSelected);
}
