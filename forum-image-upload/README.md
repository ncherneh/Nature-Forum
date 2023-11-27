# Forum-image-upload

### Introduction
This project consists in creating a web forum that allows :
* communication between users
* associating categories to posts
* liking and disliking posts and comments
* filtering posts
* upload image 

### Running the program

#### Terminal
Make sure you have all the necessary third-party packages installed.

You can run the code with:
```
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


By **[Nikol Cherneha](https://01.kood.tech/git/ncherneh), [Svitlana Bondar](https://01.kood.tech/git/sbondar), [Valeriia Nahynaliuk](https://01.kood.tech/git/vnahynal), [Anastasiia Andriievska](https://01.kood.tech/git/aandriie), [Iryna Velychko](https://01.kood.tech/git/ivelychk)**.
