
db.createUser(
    {
       user: "admin", 
       pwd: "31@$t1c$h1ftp@zz", 
       roles:["userAdminAnyDatabase"]
    }
)

db.createUser(
    {
      user: "elasticshift",
      pwd: "3l@$t1c$h1ft",
      roles: [ "readWrite", "dbAdmin" ]
    }
 )