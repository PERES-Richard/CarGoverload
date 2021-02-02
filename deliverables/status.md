# Week 04: GREEN

## Ce que l'on a fait

- création des user stories
- réfactor de notre diagramme d’architecture pour prendre en compte notre nouveau scénario, passage d’une architecture service en micro-service événementiel (motivé par le fait que l’on va faire les calculs d’itinéraire en parallèle et que l’on souhaite découplé notre architecture)

## Ce que l'on a prévu

- Init de nos nouveaux services : ItineraryDispatcher, SearchingAgregator, MultiSearchingAgregator, OffersCreator, OrderValidator
- Refactor des services existants
- Intégration de tous les services au bus Kafka

## Les blocages et risques

- manquer de temps car notre changement d’architecture est important



# Week 05: GREEN

## Ce que l'on a fait

- Mise en place du bus kafka et intégration des services existants (sauf BookingProcess)

## Ce que l'on a prévu

- Refactoring côté métier de chaque service, ajout des nouveaux services présent dans l'architecture

## Les blocages et risques

- Mise en pratique de l'architecture conçue au préalable, donc on doit se préparer à des problèmes imprévus (flot de données, intégration...)