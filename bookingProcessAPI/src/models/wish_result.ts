export class WishResult {
    wishId: string;
    offerPossibilities: OfferPossibility[]
}

export class OfferPossibility {
    searchId: string;
    offers: Offer[]
    TotalPrice: number;
}

export class Offer {
    bookDate: string;
    arrivalNode: NodePoint;
    departureNode: NodePoint;
    car: Car;
    distance: number;
    price: number;
}

export class NodePoint {
    name: string;
    latitude: number;
    longitude: number;
}

export class Car {
    id: number;
    carType: CarType;
}

export class CarType {
    name: string;
}