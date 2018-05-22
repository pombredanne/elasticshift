conn = new Mongo("ubvmmgo01:27017") // Change host and port accordingly
printjson(conn)

// conn = new Mongo() // Change host and port accordingly
esh = db.getDB('elasticshift')
//esh = db.getSiblingDB('elasticshift')
printjson(esh)

esh.createUser(
    {
       user: "admin", 
       pwd: "31@$t1c$h1ftp@zz", 
       roles:["dbAdmin"]
    }
)

esh.createUser(
    {
      user: "elasticshift",
      pwd: "3l@$t1c$h1ft",
      roles: [ "readWrite" ]
    }
)

// if database is protected with auth credentials
// db.auth(<username>, <password>)
printjson('showing collections')
cursor = esh.collection.find();
while ( cursor.hasNext() ) {
    printjson('collections...')
   printjson( cursor.next() );
}