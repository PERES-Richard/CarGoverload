#Current flag : Week 44 - GREEN

# Week 40 : GREEN

## Ce que l'on a fait
Cette semaine nous avons réfléchi au scénario principal ainsi qu'aux scénarios alternatifs de notre système. Nous avons également créé un diagramme de cas d'utilisation.

## Ce que l'on a prévu
Pour la semaine prochaine nous avons prévu de mettre en place le diagramme de composant en fonction des retours sur les scénarios. 

## Les blocages et risques
Nous pensons que le diagramme de use case est plutôt faible / inutile. Il faudra peut être en refaire un en fonction des retours

# Week 41: GREEN

## Ce que l'on a fait
Cette semaine nous avons revu notre scénario principal suite aux retours. Nous l'avons détaillé et affiné pour qu'il soit plus compréhensible. Nous avons également finalisé le diagramme de composant. Nous avons enfin définis la roadmap et le planning.

## Ce que l'on a prévu
Nous prévoyons de commencer à mettre en place le code.

## Les blocages et risques
Nous avons longuement réfléchi à l'architecture du diagramme de composant, cependant, si les nouveaux retours sont négatifs concernant celui-ci ou les scénarios, nous aurons surement du mal à rattraper ce retard.

# Week 42: GREEN

## Ce que l'on a fait

Cette semaine nous avons initialisé chaque composant en golang, avec une dockerisation du tout.

## Ce que l'on a prévu

Pour la semaine prochaine, nous finirons d'intégrer la CI (Circle CI), et commencerons à créer un produit répondant à un scénario très minimal avec une route /book basique pour créer une véritable communication entre les composants.

## Les blocages et risques

Le risque principal provient de la technologie choisie qui est nouvelle pour nous (à la fois concernant Golang et CircleCI). En dehors de cela, il n'y a pas de problème particulier en vue pour l'instant.

# Week 43: GREEN

## Ce que l'on a fait

- Intégration de tous les composants (OK à 90%)
- Refactor de la structure de nos composants
- Début de la persistance de nos réservations
- Avancement de la logique métier nécessaire pour la POC (Notamment au niveau de la recherche et de la proposition d'offres disponibles à nos fournisseurs)
 
## Ce que l'on a prévu

- Finir l'intégration
- Finir la persistance (que ce soit au niveau des réservations ou bien du stockage des offres proposés aux fournisseurs)
- Finir la configuration de la CI 
- Continuer la logique métier

## Les blocages et risques

Pas de problème prévisible pour le moment


# Week 44: GREEN

## Ce que l'on a fait

- Intégration de tous les composants OK
- Persistance de nos réservations OK
- Changement de MongoDB à PostgreSQL
- Configuration de la CI OK
 
## Ce que l'on a prévu

- Persistance des offres
- Mise en place de la CLI/Front
- Amélioration de la logique de la recherche et de la disponibilité
- Verification que le scenario pour la démo est entièrement réalisable et valide, et rajout de logique métier si besoin est.

## Les blocages et risques

Léger contretemps au niveau de la CLI (qui avait été prévue initialement pour la semaine 42), mais qui ne devrait pas trop impacter le reste du développement.