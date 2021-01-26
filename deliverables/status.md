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
