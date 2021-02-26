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
let buttonPay = null;
let wishesContainer = document.getElementById("offers-result-container");

let formSelectedDepartureNode = availableNodes[0];
let formSelectedDepartureDate = null;
let formSelectedSupplier = null;

let selectedWishes = {};
let wishId = null;
let searchesInWish = 0;

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
        this.price = databaseEntry.price;
        this.distance = databaseEntry.distance;
        this.departure = new Node(databaseEntry.departureNode.id, databaseEntry.departureNode.name, []);
        this.arrival = new Node(databaseEntry.arrivalNode.id, databaseEntry.arrivalNode.name, []);
        this.date = new Date(databaseEntry.dateDeparture);
        this.car = new Car(databaseEntry.car);
        this.departureTime = this.timeWithZeros(this.date.getHours()) + ':' + this.timeWithZeros(this.date.getMinutes());

        this.duration = this.distance / 200; // 200km / 200km*h = 1h
        this.arrivalTime = this.date.addHours(this.duration);
        this.arrivalTime = this.timeWithZeros(this.arrivalTime.getHours()) + ':' + this.timeWithZeros(this.arrivalTime.getMinutes());
    }

    timeWithZeros(time)
    {
        return (time < 10 ? '0' : '') + time;
    }
}

class OfferPossibility{
    constructor(databaseEntry) {
        this.searchId = databaseEntry.searchId;
        this.offers = [];
        databaseEntry.offers.forEach(offer => {
            this.offers.push(new Offer(offer));
        })
    }
}

class WishPossibilities{
    constructor(databaseEntry) {
        wishId = databaseEntry.wishId;
        this.wishId = databaseEntry.wishId;
        this.offerPossibilities = [];
        databaseEntry.offerPossibilities.forEach(offerPossibility => {
            this.offerPossibilities.push(new OfferPossibility(offerPossibility));
        })
        searchesInWish = this.offerPossibilities.length;
    }
}

(function() {
    Date.prototype.addHours = function(h) {
        this.setTime(this.getTime() + (h*60*60*1000));
        return this;
    }
    mainContainer = document.getElementById('main-container');
    formSubmitButton = document.getElementById('middle-form-submit');
    loadingBigContainer = document.getElementById('loading-big-container');
    buttonPay = document.getElementById('button-pay');
    buttonPay.addEventListener('click', function() {
        let payRequest = new XMLHttpRequest();
        payRequest.open('POST', 'http://localhost/booking-process/booking/payment', true);
        payRequest.addEventListener('readystatechange', function() {
            if(this.readyState === 4 && this.status === 201) {
                const response = this.responseText;
                if (response === "OK") {
                    alert("Merci d'avoir fait confiance à CarGoverload");
                }
            }
        });
        payRequest.setRequestHeader('Content-Type', 'application/json');
        const result = {
            wishId: wishId,
            wishes: []
        };
        for (let key in selectedWishes) {
            const selectedWish = selectedWishes[key];
            console.log(selectedWish)
            result.wishes.push({
                searchId: key,
                carId: selectedWish.car.id,
                departureNode: selectedWish.departure.name,
                arrivalNode: selectedWish.arrival.name,
                dateDeparture: selectedWish.date.toISOString(),
            });
        }
        console.log(result);
        payRequest.send(JSON.stringify(result));

    });
    initAvailableNodes();
    initDateSelect();
    initSupplierSelect();
    initClickOnPlus();
    document.getElementById('middle-form').addEventListener('submit', handleForm);
})();

function checkIfEnableButtonPay() {
    let totalPrice = 0;
    let size = 0;
    for (let key in selectedWishes) {
        size++;
        totalPrice += selectedWishes[key].price;
    }
    buttonPay.innerText = "Payer " + totalPrice.toFixed(2) + "€";
    if (size === searchesInWish) {
        buttonPay.disabled = false;
        return;
    }
    buttonPay.disabled = true;
}

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
            numberOfCars: parseInt(wagonsNumber),
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
    nodeDepartureInput.addEventListener('change', function() {
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

function launchCheckSearchResult(wishId){
	let checkSearch = setInterval(()=>{
		let searchRequest = new XMLHttpRequest();
		searchRequest.open('GET', 'http://localhost/booking-process/offers/' + wishId, true);
		searchRequest.addEventListener('readystatechange', function() {
			if(this.readyState === 4 && this.status === 200) {
				if(this.responseText != null && this.responseText.length > 0){
                    clearInterval(checkSearch);
					const response = JSON.parse(this.responseText); //Todo display offer
                    selectedWishes = {};
					buildOffers(response)
				}
			}
		});
		searchRequest.setRequestHeader('Content-Type', 'application/json');
		searchRequest.send();
	}, 1000);
}

function launchSearch(searches){
    buttonPay.disabled = true;
    buttonPay.innerText = "Payer";
    wishesContainer.innerHTML = "";
    addLoader();
    let searchRequest = new XMLHttpRequest();
    searchRequest.open('POST', 'http://localhost/booking-process/offers', true);
    searchRequest.addEventListener('readystatechange', function() {
        if(this.readyState === 4 && this.status === 201) {
            const response = JSON.parse(this.responseText);
            const wishId = response.wishId;
            launchCheckSearchResult(wishId);
        }
    });
    searchRequest.setRequestHeader('Content-Type', 'application/json');
    searchRequest.send(JSON.stringify(searches));
}

function buildOffers(offerPossibilities) {
    const wishPossibilities = new WishPossibilities(offerPossibilities);
    if (wishPossibilities.offerPossibilities.length === 0) {
        displayEmptyText();
    } else {
        displayWish(wishPossibilities);
    }
}

function displayWish(wishPossibilities) {
    const wishId = document.createElement("div");
    wishId.classList.add("full-width");
    wishId.id = "wish-id";
    wishId.appendChild(document.createTextNode("WishId: " + wishPossibilities.wishId));
    wishesContainer.appendChild(wishId);

    const searchesContainer = document.createElement("div");
    searchesContainer.classList.add("searches-container");
    wishPossibilities.offerPossibilities.forEach(search => {
        displaySearch(searchesContainer, search);
    });
    wishesContainer.appendChild(searchesContainer);
    removeLoader();
}

function displaySearch(container, search) {
    const div = document.createElement("div");
    div.classList.add("search-container");

    const searchId = document.createElement("div");
    searchId.classList.add("full-width");
    searchId.appendChild(document.createTextNode("SearchId: " + search.searchId));
    div.appendChild(searchId);

    search.offers.forEach(offer => {
        displayOffer(div, search.searchId, offer);
    });
    container.appendChild(div);
}

function displayOffer(bigContainer, searchId, offer){
    let container = document.createElement('div');
    container.classList.add('offer');
    container.title = 'Réserver';
    let information = document.createElement('div');
    information.classList.add('offer-information');
    information.classList.add('offer-container');
    container.classList.add('car-id-' + offer.car.id);

    let nodes = document.createElement('div');
    nodes.classList.add('offer-nodes-container');

    let node = document.createElement('div');
    node.classList.add('offer-node');

    let time = document.createElement('time');
    time.classList.add('offer-time');
    time.appendChild(document.createTextNode(offer.departureTime + ' - '));
    node.appendChild(time);
    node.appendChild(document.createTextNode(offer.departure.name));
    nodes.appendChild(node);

    node = document.createElement('div');
    node.classList.add('offer-node');
    time = document.createElement('time');
    time.classList.add('offer-time');
    time.appendChild(document.createTextNode(offer.arrivalTime + ' - '));
    node.appendChild(time);
    node.appendChild(document.createTextNode(offer.arrival.name));
    nodes.appendChild(node);

    let dateContainer = document.createElement('div');
    dateContainer.classList.add('offer-date');
    dateContainer.appendChild(document.createTextNode('Date de départ : ' + offer.date.toDateString()));
    nodes.appendChild(dateContainer);

    information.appendChild(nodes);

    let carContainer = document.createElement('div');
    carContainer.classList.add('offer-car');

    let carIcon = document.createElement('i');
    carIcon.classList.add('offer-car-icon');
    carIcon.classList.add('fas');
    carIcon.classList.add('fa-train');
    carContainer.appendChild(carIcon);

    let carInfoContainer = document.createElement('div');
    carInfoContainer.classList.add('offer-car-information')

    let carText = document.createElement('div');
    carText.appendChild(document.createTextNode('Wagon n°' + offer.car.id))
    carInfoContainer.appendChild(carText);

    let carType = document.createElement('div');
    carType.appendChild(document.createTextNode(offer.car.carType.name))
    carInfoContainer.appendChild(carType);

    carContainer.appendChild(carInfoContainer)
    information.appendChild(carContainer)

    container.appendChild(information);
    let price = document.createElement('div');
    price.classList.add('offer-price');
    price.classList.add('offer-container');
    price.appendChild(document.createTextNode(offer.price.toFixed(2) + '€'));
    container.appendChild(price);

    container.addEventListener('click', function() {
        selectedWishes[searchId] = offer;
        const previous = bigContainer.querySelector(".offer-price-selected");
        if (previous != null)
            previous.classList.remove('offer-price-selected');
        price.classList.add('offer-price-selected');
        checkIfEnableButtonPay();
    });
    bigContainer.appendChild(container);
}


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
