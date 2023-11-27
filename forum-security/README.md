# Forum-Security

### Introduction
For this project you must take into account the security of your forum.

* You should implement a Hypertext Transfer Protocol Secure (HTTPS) protocol :

        Encrypted connection : for this you will have to generate an SSL certificate, you can think of this like a identity card for your website. You can create your certificates or use "Certificate Authorities"(CA's)

        
* The implementation of Rate Limiting must be present on this project

* You should encrypt at least the clients passwords. As a Bonus you can also encrypt the database, for this you will have to create a password for your database.

Sessions and cookies were implemented in the previous project but not under-pressure (tested in an attack environment). So this time you must take this into account.

* Clients session cookies should be unique. For instance, the session state is stored on the server and the session should present an unique identifier. This way the client has no direct access to it. Therefore, there is no way for attackers to read or tamper with session state.

### Audit

https://github.com/01-edu/public/blob/master/subjects/forum/security/audit.md

We created our own SSL certificate and most browsers during first load of this forum will issue a warning about unofficial certificate. Please, give permission and dont worry.

### Running the program

#### Terminal
Make sure you have all the necessary third-party packages installed.

You can run the code with:
```go
go run .
```
You can use the extension in Visual Studio Code - SQLite:
Open database.sql and point to "SELECT * FROM users;" right-click and select Run Query 

Or You can use the extension in Visual Studio Code - SQLite Viewer: 
Go to database.db click on Open Anyway and a window will pop up where you need to select SQLite Viewer

#### Docker
You can build the docker file with 
```
bash dockerrun.sh
```
or chmod +x dockerrun.sh and ./dockerrun.sh

after you are done you can remove the files created by docker with 
```
bash dockerclean.sh
```
or chmod +x dockerclean.sh and ./dockerclean.sh


By **[Iryna Velychko](https://01.kood.tech/git/ivelychk), [Nikol Cherneha](https://01.kood.tech/git/ncherneh), [Svitlana Bondar](https://01.kood.tech/git/sbondar), [Valeriia Nahynaliuk](https://01.kood.tech/git/vnahynal), [Anastasiia Andriievska](https://01.kood.tech/git/aandriie)**.
