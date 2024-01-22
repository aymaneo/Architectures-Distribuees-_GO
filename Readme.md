# Architectures-Distribuees_GO

## Présentation
Ce dépôt contient le code pour une API dédiée à la surveillance des conditions météorologiques dans un aéroport. L'API, écrite en Go, permet de récupérer et d'agréger des données provenant de divers capteurs simulés.

## Fonctionnalités
- **Récupération des Mesures** : Permet d'obtenir des mesures en temps réel ou historiques pour des périodes spécifiques.
- **Calcul de la Moyenne** : Fournit la valeur moyenne des mesures pour une journée donnée.
- **Accès aux Moyennes Globales** : Accès aux moyennes de toutes les mesures sur une journée spécifique.

## Dépendances
- **InfluxDB** : Base de données NoSQL pour le stockage des mesures.
- **Gorilla Mux** : Routeur HTTP pour le service web.

## Endpoints
- `GET /api/mesure/{mesures}/{iata}/{start}/{end}` : Récupère les mesures entre deux dates.
- `GET /api/mesure/{mesures}/{iata}/{start}` : Calcule la moyenne des mesures pour une date donnée.
- `GET /api/allMeans/{iata}/{start}` : Récupère les moyennes pour tous les types de mesures pour une date donnée.

## Comment démarrer
1. **Initialisation** : Clonez le dépôt et installez les dépendances.
2. **Configuration** : Configurez votre instance InfluxDB et assurez-vous qu'elle est accessible depuis l'API.
3. **Lancement** : Exécutez le serveur avec `go run main.go`.

L'API est documentée avec Swagger 
