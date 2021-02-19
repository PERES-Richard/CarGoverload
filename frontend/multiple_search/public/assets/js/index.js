const availableNodes = [
    'Avignon-solid',
    'Avignon-liquid',
    'Marseille',
    'Nice',
    'Paris'
];

const carTypes = [
    'Liquid',
    'Solid'
];

let loadingBigContainer = null;

let mainContainer = null;
let formSubmitButton = null;

let formSelectedDepartureNode = availableNodes[0];
let formSelectedDepartureDate = null;
let formSelectedSupplier = null;

(function() {
    mainContainer = document.getElementById('main-container');
    formSubmitButton = document.getElementById('middle-form-submit');
    loadingBigContainer = document.getElementById('loading-big-container');
    initAvailableNodes();
    initDateSelect();
    initSupplierSelect();
    initClickOnPlus();
    document.getElementById('middle-form').addEventListener('submit', handleForm);
})();

function handleForm(e) {
    e.preventDefault();
    const departureNode = formSelectedDepartureNode;
    const departureDate = formSelectedDepartureDate;
    const supplier = formSelectedSupplier;
    const searches = [];

    const wishes = mainContainer.querySelectorAll('.new-wish-line');
    for (let i = 0; i < wishes.length; i++) {
        const wish = wishes[i];
        let arrivalNode = wish.querySelector('#node-departure-select').value;
        let carType = wish.querySelector('#car-type-select').value;
        let wagonsNumber = wish.querySelector('#number-cars-input').value;
        searches.push({
            departureNode: departureNode,
            carType: carType,
            numberOfCars: wagonsNumber,
            arrivalNode: arrivalNode,
            dateDeparture: departureDate
        })
    }

    launchSearch(searches);
}

function checkIfFormValid() {
    if (formSelectedDepartureNode === null || formSelectedDepartureDate === null || formSelectedSupplier === null ||
        formSelectedSupplier.length === 0 || mainContainer.childNodes.length === 0) {
        formSubmitButton.disabled = true;
        return;
    }
    const wishes = mainContainer.querySelectorAll('.new-wish-line');
    for (let i = 0; i < wishes.length; i++) {
        const wish = wishes[i];
        let arrivalNode = wish.querySelector('#node-departure-select').value;
        let carType = wish.querySelector('#car-type-select').value;
        let wagonsNumber = wish.querySelector('#number-cars-input').value;
        if (arrivalNode === null || carType === null || wagonsNumber == null || arrivalNode === formSelectedDepartureNode || wagonsNumber.length === 0 || parseInt(wagonsNumber) === 0) {
            formSubmitButton.disabled = true;
            return;
        }
    }
    formSubmitButton.disabled = false;
}

function initAvailableNodes() {
    const nodeDepartureInput = document.getElementById('node-departure-select');
    availableNodes.forEach((node) => {
       const option = document.createElement('option');
       option.value = node;
       option.appendChild(document.createTextNode(node));
       nodeDepartureInput.appendChild(option);
    });
    nodeDepartureInput.addEventListener('change', function(e) {
        formSelectedDepartureNode = nodeDepartureInput.value;
        checkIfFormValid();
        console.log('Selected departure node : ' + formSelectedDepartureNode);
    });
}

function initDateSelect() {
    document.getElementById('date-departure').addEventListener('input', function(e){
        formSelectedDepartureDate = (new Date(e.target.value)).toISOString();
        checkIfFormValid();
        console.log('Selected departure date : ' + formSelectedDepartureDate);
    });
}

function initSupplierSelect() {
    document.getElementById('supplier').addEventListener('input', function(e){
        formSelectedSupplier = e.target.value;
        checkIfFormValid();
        console.log('Selected supplier : ' + formSelectedSupplier);
    });
}

function initClickOnPlus() {
    document.getElementById('plus').addEventListener('click', initNewWishLine);
}

function initNewWishLine() {
    const container = document.createElement('div');
    container.classList.add('new-wish-line');

    container.appendChild(createSelectArrivalNode());
    container.appendChild(createSelectCarType());
    container.appendChild(createNumberCarsInput());

    let button = document.createElement('div');
    button.classList.add('no-button-form');
    button.appendChild(document.createTextNode('Supprimer'));
    button.addEventListener('click', function() {
       container.remove();
        checkIfFormValid();
    });
    container.appendChild(button);

    mainContainer.appendChild(container);
    checkIfFormValid();
}

function createSelectArrivalNode() {
    let subContainer = document.createElement('div');
    subContainer.classList.add('input-container');
    let label = document.createElement('label');
    label.htmlFor = 'node-departure-select';
    label.classList.add('color-grey');
    label.appendChild(document.createTextNode('Noeud d\'arrivée'));
    subContainer.appendChild(label);
    let div = document.createElement('div');
    div.classList.add('no-margin');
    let select = document.createElement('select');
    select.classList.add('background-main-color');
    select.classList.add('color-grey');
    select.title = 'Noeud d\'arrivée';
    select.name = 'node-departure-select';
    select.id = 'node-departure-select';
    select.addEventListener('change', checkIfFormValid);
    availableNodes.forEach(node => {
        let option = document.createElement('option');
        option.value = node;
        option.appendChild(document.createTextNode(node));
        select.appendChild(option);
    });
    div.appendChild(select);
    subContainer.appendChild(div);
    return subContainer;
}

function createSelectCarType() {
    let subContainer = document.createElement('div');
    subContainer.classList.add('input-container');
    let label = document.createElement('label');
    label.htmlFor = 'car-type-select';
    label.classList.add('color-grey');
    label.appendChild(document.createTextNode('Type de wagon'));
    subContainer.appendChild(label);
    let div = document.createElement('div');
    div.classList.add('no-margin');
    let select = document.createElement('select');
    select.classList.add('background-main-color');
    select.classList.add('color-grey');
    select.title = 'Type de wagon';
    select.name = 'car-type-select';
    select.id = 'car-type-select';
    select.addEventListener('change', checkIfFormValid);
    carTypes.forEach(node => {
        let option = document.createElement('option');
        option.value = node;
        option.appendChild(document.createTextNode(node));
        select.appendChild(option);
    });
    div.appendChild(select);
    subContainer.appendChild(div);
    return subContainer;
}

function createNumberCarsInput() {
    let subContainer = document.createElement('div');
    subContainer.classList.add('input-container');
    let label = document.createElement('label');
    label.htmlFor = 'number-cars-input';
    label.classList.add('color-grey');
    label.appendChild(document.createTextNode('Nombre de wagons désirés'));
    subContainer.appendChild(label);
    let div = document.createElement('div');
    div.classList.add('no-margin');
    let input = document.createElement('input');
    input.id = 'number-cars-input';
    input.type = 'number';
    input.min = '1';
    input.addEventListener('input', checkIfFormValid);
    div.appendChild(input);
    subContainer.appendChild(div);
    return subContainer;
}

function addLoader(){
    let loadingContainer = document.getElementById('loading-container');
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
    loadingBigContainer.innerHTML = '';
}

function displayEmptyText(){
    let container = document.createElement('div');
    container.classList.add('text-empty');
    container.appendChild(document.createTextNode('Aucun wagon de libre avec vos critères'));
    mainContainer.appendChild(container);
}

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
//
// class Offer{
//     constructor(databaseEntry) {
//         this.id = databaseEntry.id;
//         this.price = databaseEntry.price;
//         this.date = new Date(databaseEntry.beginBookedDate);
//         this.departure = new Node(databaseEntry.departureNode.id, databaseEntry.departureNode.name, []);
//         this.arrival = new Node(databaseEntry.arrivalNode.id, databaseEntry.arrivalNode.name, []);
//         this.car = new Car(databaseEntry.car);
//         this.duration = databaseEntry.duration;
//         this.departureTime = this.timeWithZeros(this.date.getHours()) + ':' + this.timeWithZeros(this.date.getMinutes());
//         this.dateArrival = new Date(databaseEntry.endingBookedDate);
//         this.arrivalTime = this.timeWithZeros(this.dateArrival.getHours()) + ':' + this.timeWithZeros(this.dateArrival.getMinutes());
//     }
//
//     timeWithZeros(time)
//     {
//         return (time < 10 ? '0' : '') + time;
//     }
//
// }


function launchSearch(searches){
    addLoader();
    let searchRequest = new XMLHttpRequest();
    searchRequest.open('POST', 'http://localhost/booking-process/offers', true);
    searchRequest.addEventListener('readystatechange', function() {
        if(this.readyState === 4 && this.status === 200) {
            console.log("ICI")
        }
    });
    searchRequest.setRequestHeader('Content-Type', 'application/json');
    searchRequest.send(JSON.stringify(searches));
}

// function displayOffer(offer){
//     // let container = document.createElement('div');
//     // container.classList.add('offer');
//     // container.title = 'Réserver';
//     // let information = document.createElement('div');
//     // information.classList.add('offer-information');
//     // information.classList.add('offer-container');
//     // container.classList.add('car-id-' + offer.car.id);
//     //
//     // let nodes = document.createElement('div');
//     // nodes.classList.add('offer-nodes-container');
//     //
//     // let node = document.createElement('div');
//     // node.classList.add('offer-node');
//     //
//     // let time = document.createElement('time');
//     // time.classList.add('offer-time');
//     // time.appendChild(document.createTextNode(offer.departureTime + ' - '));
//     // node.appendChild(time);
//     // node.appendChild(document.createTextNode(offer.departure.name));
//     // nodes.appendChild(node);
//     //
//     // node = document.createElement('div');
//     // node.classList.add('offer-node');
//     // time = document.createElement('time');
//     // time.classList.add('offer-time');
//     // time.appendChild(document.createTextNode(offer.arrivalTime + ' - '));
//     // node.appendChild(time);
//     // node.appendChild(document.createTextNode(offer.arrival.name));
//     // nodes.appendChild(node);
//     //
//     // let dateContainer = document.createElement('div');
//     // dateContainer.classList.add('offer-date');
//     // dateContainer.appendChild(document.createTextNode('Date de départ : ' + offer.date.toDateString()));
//     // nodes.appendChild(dateContainer);
//     //
//     // information.appendChild(nodes);
//     //
//     // let carContainer = document.createElement('div');
//     // carContainer.classList.add('offer-car');
//     //
//     // let carIcon = document.createElement('i');
//     // carIcon.classList.add('offer-car-icon');
//     // carIcon.classList.add('fas');
//     // carIcon.classList.add('fa-train');
//     // carContainer.appendChild(carIcon);
//     //
//     // let carInfoContainer = document.createElement('div');
//     // carInfoContainer.classList.add('offer-car-information')
//     //
//     // let carText = document.createElement('div');
//     // carText.appendChild(document.createTextNode('Wagon n°' + offer.car.id))
//     // carInfoContainer.appendChild(carText);
//     //
//     // let carType = document.createElement('div');
//     // carType.appendChild(document.createTextNode(offer.car.carType.name))
//     // carInfoContainer.appendChild(carType);
//     //
//     // carContainer.appendChild(carInfoContainer)
//     // information.appendChild(carContainer)
//     //
//     // container.appendChild(information);
//     // let price = document.createElement('div');
//     // price.classList.add('offer-price');
//     // price.classList.add('offer-container');
//     // price.appendChild(document.createTextNode(offer.price + '€'));
//     // container.appendChild(price);
//     //
//     // container.addEventListener('click', function(){
//     //     const res = confirm('Procéder à la réservation et au paiement ?')
//     //     if(res){
//     //         book(offer, container);
//     //     }
//     // });
//
//     // mainContainer.appendChild(container);
// }
//

//
// function book(offer, view){
//     // TODO book
//     // let bookOffer = new XMLHttpRequest();
//     // bookOffer.open('POST', 'http://localhost/booking-process/offers/payment', true);
//     // bookOffer.addEventListener('readystatechange', function() {
//     //     if(this.readyState === 4 && this.status === 200) {
//     //         launchSearch();
//     //     }
//     // });
//     // bookOffer.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
//     // bookOffer.send(JSON.stringify(
//     //     {
//     //         offerId: offer.id,
//     //         supplier: supplierInput.value
//     //     }
//     // ));
// }
