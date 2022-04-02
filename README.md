# user-registration
## used mysql database as conatiner, backend code in golang and used go-cache for caching

## Steps to build application
```bash
# commands shown below run at the root of the project directory
1. To build golang code binary and mysql custom image
  make all
```

## Steps to run mysql container
```bash
# commands shown below run at the root of the project directory
 mkdir -p /var/tmp/mysql_data
 cd ./db
 docker-compose up
```

## Steps to run user-registration binary
```bash
# commands shown below run at the root of the project directory

./bin/userRegistration --> http server
    Usage of ./bin/userRegistration:
      -db string
          mysql database connection in <username>:<pwd>@tcp(<host>:<port>)/<db> [compulsory option]
      -port string
          HTTP server port,default port 8080 (default "8080")
   
   example:
     to get db mysql database connection
     ip = docker inspect db_mysql_db_1 | grep IPAddres 

     SELECT host, user FROM mysql.user;
     CREATE USER 'pavan'@'<ip/host>' IDENTIFIED BY 'pavan';
     GRANT ALL ON registration.* TO 'pavan'@'<ip/host>';
  
      ./bin/userRegistration -db "kumar:kumar@tcp(127.0.0.1:3307)/registration" -port "60010"
  
```
## http endpoints supported in user-registration binary
```bash
1.    SignUp endpoint -
            curl --location --request POST 'http://127.0.0.1:60010/api/signup' \
            --header 'Content-Type: application/json' \
            --data-raw '{
                "user": {
                    "firstname": "pavan",
                    "lastname": "kumar",
                    "email": "pavank830@gmail.com",
                    "password": "kumar"
                }
            }'
  
2.    login endpoint -
            curl --location --request POST 'http://127.0.0.1:60010/api/login' \
            --header 'Content-Type: application/json' \
            --data-raw '{
                "email": "pavank830@gmail.com",
                "password": "kumar"
            }'
  
3.    get user profile data endpoint -
            curl --location --request GET 'http://127.0.0.1:60010/api/profile' \
            --header 'Authorization: Bearer <JWT_Token>'
            
4.    logout endpoint -
            curl --location --request POST 'http://127.0.0.1:60010/api/logout' \
            --header 'Authorization: Bearer <JWT_Token>'            
```



## Brief on the entire application
 Two tables are created in db one is 'user' table which stores user basic details and other is 'blacklist' table which stores the blacklisted jwt tokens.
 On SignUp or login, jwt token is returned with an expiry of 24 hrs.
 To get user profile data or to logout, jwt token must be sent as Bearer token.
 On logout, the jwt token is blacklisted [to avoid from re-using the same token after logout].
