const mongoHost = process.env.NAPMAP_API_MONGODB_HOST
const mongoPort = process.env.NAPMAP_API_MONGODB_PORT

const mongoUser = process.env.NAPMAP_API_MONGODB_USERNAME
const mongoPassword = process.env.NAPMAP_API_MONGODB_PASSWORD

const database = process.env.NAPMAP_API_MONGODB_DATABASE
const collection = process.env.NAPMAP_API_MONGODB_COLLECTION

const retrySeconds = parseInt(process.env.RETRY_CONNECTION_SECONDS || "5") || 5;

// try to connect to mongoDB until it is not available
let connection;
while(true) {
    try {
        connection = Mongo(`mongodb://${mongoUser}:${mongoPassword}@${mongoHost}:${mongoPort}`);
        break;
    } catch (exception) {
        print(`Cannot connect to mongoDB: ${exception}`);
        print(`Will retry after ${retrySeconds} seconds`)
        sleep(retrySeconds * 1000);
    }
}

// if database and collection exists, exit with success - already initialized
const databases = connection.getDBNames()
if (databases.includes(database)) {
    const dbInstance = connection.getDB(database)
    const collections = dbInstance.getCollectionNames()
    if (collections.includes(collection)) {
        print(`Collection '${collection}' already exists in database '${database}'`)
        process.exit(0);
    }
}

// initialize
const db = connection.getDB(database)
db.createCollection(collection)

// create indexes
db[collection].createIndex({ "id": 1 })
db[collection].createIndex({ "city": 1 })
db[collection].createIndex({ "status": 1 })

// insert sample data
let result = db[collection].insertMany([
    {
        "id": "st-001",
        "name": "NAPMap Bratislava Nivy",
        "stationType": "CHARGING",
        "fuels": ["ELECTRIC"],
        "operatorName": "NAPMap Energy s.r.o.",
        "address": "Mlynské nivy 16",
        "city": "Bratislava",
        "country": "SK",
        "lat": 48.1486,
        "lng": 17.1077,
        "openingHours": "24/7",
        "maxPowerKw": 150,
        "connectors": ["CCS2", "Type2", "CHAdeMO"],
        "services": ["WC", "Parking"],
        "status": "ACTIVE"
    },
    {
        "id": "st-002",
        "name": "NAPMap Petržalka",
        "stationType": "CHARGING",
        "fuels": ["ELECTRIC"],
        "operatorName": "NAPMap Energy s.r.o.",
        "address": "Einsteinova 25",
        "city": "Bratislava",
        "country": "SK",
        "lat": 48.1322,
        "lng": 17.1067,
        "openingHours": "06:00-22:00",
        "maxPowerKw": 50,
        "connectors": ["Type2"],
        "services": ["Parking"],
        "status": "ACTIVE"
    },
    {
        "id": "st-003",
        "name": "NAPMap Trnava Centrum",
        "stationType": "CHARGING",
        "fuels": ["ELECTRIC"],
        "operatorName": "NAPMap Energy s.r.o.",
        "address": "Hlavná 5",
        "city": "Trnava",
        "country": "SK",
        "lat": 48.3774,
        "lng": 17.5885,
        "openingHours": "24/7",
        "maxPowerKw": 100,
        "connectors": ["CCS2", "Type2"],
        "services": ["WC", "Food"],
        "status": "ACTIVE"
    },
    {
        "id": "st-004",
        "name": "NAPMap Nitra Zobor",
        "stationType": "CHARGING",
        "fuels": ["ELECTRIC"],
        "operatorName": "NAPMap Energy s.r.o.",
        "address": "Štefánikova 8",
        "city": "Nitra",
        "country": "SK",
        "lat": 48.3069,
        "lng": 18.0869,
        "openingHours": "06:00-20:00",
        "maxPowerKw": 75,
        "connectors": ["CCS2", "Type2"],
        "services": ["Parking"],
        "status": "ACTIVE"
    },
    {
        "id": "st-005",
        "name": "NAPMap Žilina Station",
        "stationType": "CHARGING",
        "fuels": ["ELECTRIC"],
        "operatorName": "NAPMap Energy s.r.o.",
        "address": "Vysokoškolákov 52",
        "city": "Žilina",
        "country": "SK",
        "lat": 49.2194,
        "lng": 18.7408,
        "openingHours": "24/7",
        "maxPowerKw": 150,
        "connectors": ["CCS2", "Type2", "CHAdeMO"],
        "services": ["WC", "Food", "Parking"],
        "status": "ACTIVE"
    },
    {
        "id": "st-006",
        "name": "NAPMap Banská Bystrica",
        "stationType": "CHARGING",
        "fuels": ["ELECTRIC"],
        "operatorName": "NAPMap Energy s.r.o.",
        "address": "Námestie SNP 12",
        "city": "Banská Bystrica",
        "country": "SK",
        "lat": 48.7395,
        "lng": 19.1533,
        "openingHours": "24/7",
        "maxPowerKw": 100,
        "connectors": ["CCS2", "Type2"],
        "services": ["WC"],
        "status": "ACTIVE"
    },
    {
        "id": "st-007",
        "name": "NAPMap Prešov Hub",
        "stationType": "CHARGING",
        "fuels": ["ELECTRIC"],
        "operatorName": "NAPMap Energy s.r.o.",
        "address": "Masarykova 20",
        "city": "Prešov",
        "country": "SK",
        "lat": 49.0018,
        "lng": 21.2395,
        "openingHours": "06:00-22:00",
        "maxPowerKw": 75,
        "connectors": ["CCS2", "Type2"],
        "services": ["Parking"],
        "status": "ACTIVE"
    },
    {
        "id": "st-008",
        "name": "NAPMap Košice Terminal",
        "stationType": "CHARGING",
        "fuels": ["ELECTRIC"],
        "operatorName": "NAPMap Energy s.r.o.",
        "address": "Staničné námestie 1",
        "city": "Košice",
        "country": "SK",
        "lat": 48.7164,
        "lng": 21.2611,
        "openingHours": "24/7",
        "maxPowerKw": 150,
        "connectors": ["CCS2", "Type2", "CHAdeMO"],
        "services": ["WC", "Food", "Parking"],
        "status": "ACTIVE"
    },
    {
        "id": "st-009",
        "name": "NAPMap Senec H2 Point",
        "stationType": "REFUELING",
        "fuels": ["HYDROGEN"],
        "operatorName": "H2 Slovakia a.s.",
        "address": "Diaľničná cesta 12",
        "city": "Senec",
        "country": "SK",
        "lat": 48.2194,
        "lng": 17.3997,
        "openingHours": "06:00-22:00",
        "maxPowerKw": null,
        "connectors": [],
        "services": ["WC", "Food", "Parking"],
        "status": "ACTIVE"
    },
    {
        "id": "st-010",
        "name": "NAPMap Zvolen CNG",
        "stationType": "REFUELING",
        "fuels": ["CNG"],
        "operatorName": "CNG Slovensko s.r.o.",
        "address": "Bystrický rad 2",
        "city": "Zvolen",
        "country": "SK",
        "lat": 48.5747,
        "lng": 19.1369,
        "openingHours": "24/7",
        "maxPowerKw": null,
        "connectors": [],
        "services": ["WC", "Parking"],
        "status": "ACTIVE"
    }
]);

if (result.writeError) {
    console.error(result)
    print(`Error when writing the data: ${result.errmsg}`)
}

// exit with success
process.exit(0);
