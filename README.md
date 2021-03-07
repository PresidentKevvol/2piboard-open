# 2PiBoard by The Hornet's Nest

2PiBoard is a homebrew messaging board imitating the styles of classical message boards originating in 1990s-2000s such as 4chan. It is made with a Golang net/http backend and made because Kevvol is bored.


The code is open sourced here to let anyone interested in making a 'modern' message board tinker with it.
You can also [see the entire webapp in action here](https://2piboard.hornetsnestguild.com/).


## If you want to set up your own 2PiBoard

Requirements:
- go 1.15+
- packages needed listed in go.mod (no special compiling accomodation required)
- PostgreSQL database

.env variables:
|entry			|purpose	|
|-----------------------|---------------|
|database_username	|Username to an user in the PostgreSQL database server|
|database_password	|Password for the database account|
|database_url		|URL to the database server|
|database_name		|Name of the database within the PostgreSQL server|
|host_name		|The host name of the web server (what domain is this site hosted on), can include port to be like `www.example.com:8880`|
|ssl_cert		|the SSL certificate file for using HTTPS, optional|
|ssl_key		|the SSL certificate key for using HTTPS, optional|

Steps:
1. Setup the PostgreSQL server by creating a user and a database belonging to that user
2. Fill in needed parameters in .env according to the table above
3. Run the SQL code in `setups/create_tables.sql` using `\i` in the psql shell under said user in the SQL server
4. `cd` into the working directory and run `go build` to compile the application
5. Setup the compiled application as a background service and start it (varies with platform)
6. (optional) Link traffic to the server application with a Apache or Nginx reverse proxy

*(on second thought I really should just make a setup.sh for easy setup some time)*

Attribution:
If you take our code and make a customized/improved variation version you just need to attribute us (The Hornet's Nest) in the home page and footer
As usuals we are not responsible for anything users post on the board(s)

If you have any questions, ask us (Kevvol) in [The Hornet's Nest discord server](https://discord.gg/AbPuABJ)


Good luck in your ventures in creating/running a modern message board!

-Kevvol Mar. 2021