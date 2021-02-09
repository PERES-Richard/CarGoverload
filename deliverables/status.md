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



# Week 06: YELLOW

## Ce que l'on a fait

- Intégration des services entre eux via Kafka (début du flow vers la fin)
- Création de l'avant dernier service
- Uniformisation des entitées qui transitent
- Ajout de la logique de "range" de recherche

## Ce que l'on a prévu

- Ajouter le tout dernier service au projet
- Continuer de verifier / intégrer les services entre eux (du début vers la fin)

(*Plus tard*)
- Finir la logique des services restants (BookingProcess API, SearchingAggregator, ItineraryDispatcher et OrderCreator)
- Changer le front end ou à défaut le remplacer par des appels HTTP

## Les blocages et risques

- Integration des services restants via Kafka plus compliqué que prévu
- Logique des services restants trop complexe pour être implémenté en 1 ou 2 semaines
- Comme discuté, refactor pas assez incrémentale et donc pas d'avancement démontrable effortless
