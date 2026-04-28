const mongoHost = process.env.NAPMAP_API_MONGODB_HOST
const mongoPort = process.env.NAPMAP_API_MONGODB_PORT

const mongoUser = process.env.NAPMAP_API_MONGODB_USERNAME
const mongoPassword = process.env.NAPMAP_API_MONGODB_PASSWORD

const database = process.env.NAPMAP_API_MONGODB_DATABASE
const collection = process.env.NAPMAP_API_MONGODB_COLLECTION

const retrySeconds = parseInt(process.env.RETRY_CONNECTION_SECONDS || "5") || 5;

// retry mongo connection
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

const db = connection.getDB(database)
if (!db.getCollectionNames().includes(collection)) {
    db.createCollection(collection)
}

db[collection].createIndex({ "id": 1 }, { unique: true })
db[collection].createIndex({ "city": 1 })
db[collection].createIndex({ "status": 1 })

const seed = [
    {
        "id": "st-001",
        "name": "ZSE Drive Eurovea",
        "stationType": "CHARGING",
        "fuels": ["ELECTRIC"],
        "operatorName": "ZSE Drive",
        "address": "Pribinova 8",
        "city": "Bratislava",
        "country": "SK",
        "lat": 48.1417,
        "lng": 17.1216,
        "openingHours": "24/7",
        "maxPowerKw": 150,
        "connectors": ["CCS2", "Type2", "CHAdeMO"],
        "services": ["WC", "Parking"],
        "status": "ACTIVE"
    },
    {
        "id": "st-002",
        "name": "Greenway Petržalka",
        "stationType": "CHARGING",
        "fuels": ["ELECTRIC"],
        "operatorName": "Greenway",
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
        "name": "ZSE Drive Trnava Hlavná",
        "stationType": "CHARGING",
        "fuels": ["ELECTRIC"],
        "operatorName": "ZSE Drive",
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
        "name": "Eon Drive Nitra",
        "stationType": "CHARGING",
        "fuels": ["ELECTRIC"],
        "operatorName": "Eon Drive Slovensko",
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
        "name": "Tesla Supercharger Žilina",
        "stationType": "CHARGING",
        "fuels": ["ELECTRIC"],
        "operatorName": "Tesla",
        "address": "Vysokoškolákov 52",
        "city": "Žilina",
        "country": "SK",
        "lat": 49.2194,
        "lng": 18.7408,
        "openingHours": "24/7",
        "maxPowerKw": 150,
        "connectors": ["CCS2", "Tesla"],
        "services": ["WC", "Food", "Parking"],
        "status": "ACTIVE"
    },
    {
        "id": "st-006",
        "name": "Greenway Banská Bystrica",
        "stationType": "CHARGING",
        "fuels": ["ELECTRIC"],
        "operatorName": "Greenway",
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
        "name": "ZSE Drive Prešov",
        "stationType": "CHARGING",
        "fuels": ["ELECTRIC"],
        "operatorName": "ZSE Drive",
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
        "name": "ZSE Drive Košice Aupark",
        "stationType": "CHARGING",
        "fuels": ["ELECTRIC"],
        "operatorName": "ZSE Drive",
        "address": "Námestie osloboditeľov 1",
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
        "name": "Slovnaft H2 Senec",
        "stationType": "REFUELING",
        "fuels": ["HYDROGEN"],
        "operatorName": "Slovnaft",
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
        "name": "SPP CNG Zvolen",
        "stationType": "REFUELING",
        "fuels": ["CNG"],
        "operatorName": "SPP CNG",
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
]

let inserted = 0
for (const station of seed) {
    const res = db[collection].updateOne(
        { id: station.id },
        { $setOnInsert: station },
        { upsert: true }
    )
    if (res.upsertedCount > 0) inserted += 1
}
print(`Seed completed: ${inserted}/${seed.length} new stations inserted`)

process.exit(0);
