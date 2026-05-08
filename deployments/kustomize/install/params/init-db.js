const mongoHost = process.env.PROJECT_API_MONGODB_HOST
const mongoPort = process.env.PROJECT_API_MONGODB_PORT
const mongoUser = process.env.PROJECT_API_MONGODB_USERNAME
const mongoPassword = process.env.PROJECT_API_MONGODB_PASSWORD
const database = process.env.PROJECT_API_MONGODB_DATABASE
const collection = process.env.PROJECT_API_MONGODB_COLLECTION

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
    collections = dbInstance.getCollectionNames()
    if (collections.includes(collection)) {
      print(`Collection '${collection}' already exists in database '${database}'`)
        process.exit(0);
    }
}

// initialize
// create database and collection
const db = connection.getDB(database)
db.createCollection(collection)

// create indexes
db[collection].createIndex({ "id": 1 })

// insert sample data
let result = db[collection].insertMany([
    {
        "id": "bobulova",
        "name": "Ambulancia všeobecného lekárstva Dr. Bobuľová",
        "roomNumber": "142 - 1.posch",
        "patients": [
            {
                "id": "pat-001",
                "name": "Jožko Púčik",
                "rodneCislo": "950215/1234",
                "dateOfBirth": "1995-02-15",
                "gender": "male",
                "address": "Hlavná 12, 811 01 Bratislava",
                "phone": "+421 900 123 456",
                "email": "jozko.pucik@example.sk",
                "insurance": "VšZP",
                "bloodType": "A+",
                "allergies": "",
                "chronicConditions": "",
                "medications": "",
                "prescriptions": [
                    {
                        "id": "rx-001",
                        "medicineName": "Ibalgin",
                        "strength": "400 mg",
                        "form": "tbl.",
                        "dosage": "1-0-1-0",
                        "instructions": "Užívať po jedle",
                        "quantity": "30 tbl.",
                        "atcCode": "M01AE01",
                        "prescribedDate": "2026-04-15",
                        "validUntil": "2026-10-15",
                        "prescribedBy": "MUDr. Anna Nováková",
                        "status": "active",
                        "repeatMonths": 1,
                        "coverage": "full"
                    }
                ]
            },
            {
                "id": "pat-002",
                "name": "Bc. August Cézar",
                "rodneCislo": "780923/5678",
                "dateOfBirth": "1978-09-23",
                "gender": "male",
                "address": "Mlynská 8, 821 09 Bratislava",
                "phone": "+421 901 234 567",
                "email": "august.cezar@example.sk",
                "insurance": "Dôvera",
                "bloodType": "0-",
                "allergies": "Penicilín",
                "chronicConditions": "Hypertenzia",
                "medications": "Prestarium 5 mg",
                "prescriptions": [
                    {
                        "id": "rx-002",
                        "medicineName": "Prestarium",
                        "strength": "5 mg",
                        "form": "tbl.",
                        "dosage": "1-0-0-0",
                        "instructions": "Užívať ráno",
                        "quantity": "30 tbl.",
                        "atcCode": "C09AA04",
                        "prescribedDate": "2026-01-10",
                        "validUntil": "2027-01-10",
                        "prescribedBy": "MUDr. Peter Horváth",
                        "status": "active",
                        "repeatMonths": 12,
                        "coverage": "full"
                    }
                ]
            },
            {
                "id": "pat-003",
                "name": "Ing. Ferdinand Trety",
                "rodneCislo": "650304/9012",
                "dateOfBirth": "1965-03-04",
                "gender": "male",
                "insurance": "Union",
                "bloodType": "AB+",
                "allergies": "",
                "chronicConditions": "",
                "medications": "",
                "prescriptions": []
            }
        ]
    }
]);

if (result.writeError) {
    console.error(result)
    print(`Error when writing the data: ${result.errmsg}`)
}

// exit with success
process.exit(0);