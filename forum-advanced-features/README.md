# Forum advanced-features

### Introduction
This project consists in a subtask of Forum project that allows :
* possibility to edit and delete posts and comments
* get notifications when other users liked/disliked/commented your posts
* user activity page that shows which posts and comments you made, and posts/comments that you liked or disliked


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


By **[Nikol Cherneha](https://01.kood.tech/git/ncherneh), [Svitlana Bondar](https://01.kood.tech/git/sbondar), [Valeriia Nahynaliuk](https://01.kood.tech/git/vnahynal), [Anastasiia Andriievska](https://01.kood.tech/git/aandriie), [Iryna Velychko](https://01.kood.tech/git/ivelychk)**.
