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

let carTypeIdSelect = null;
let nodeArrivalSelect = null;
let nodeDepartureSelect = null;

function addSelectNodes(select, carTypeId){
    select.innerHTML = "";
    if(select === nodeArrivalSelect){
        addOptionToSelect(select, "Aucun noeud d'arrivée sélectionné", 0, true)
    }else{
        addOptionToSelect(select, "Aucun noeud de départ sélectionné", 0, true)
    }
    select.disabled = false;

    listNodes.forEach(function (node){
        if(carTypeId === 0 || node.hasCarType(carTypeId)){
           addOptionToSelect(select, node.name, node.id);
        }
    });
}

function addOptionToSelect(select, name, value, disabled){
    let option = document.createElement("option");
    option.value = value;
    option.disabled = disabled
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
                addOptionToSelect(nodeDepartureSelect, node.name, node.id, false);
                addOptionToSelect(nodeArrivalSelect, node.name, node.id, false);
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
                addOptionToSelect(carTypeIdSelect, carTypeObject.name, carTypeObject.id);
            });
        }
    });
    loadIdMember.send(null);
}

(function(){
    carTypeIdSelect = document.getElementById("car-type-select");

    let buttonSubmit = document.getElementById("middle-form-submit");

    carTypeIdSelect.addEventListener("change", function(e){ // when selecting an other value
        buttonSubmit.disabled = false;
        addSelectNodes(nodeDepartureSelect, parseInt(e.target.value))
        addSelectNodes(nodeArrivalSelect, parseInt(e.target.value))
    });
    nodeDepartureSelect = document.getElementById("node-departure-select");
    nodeArrivalSelect = document.getElementById("node-arrival-select");
    loadNodes();
    loadCarTypes();
})()
