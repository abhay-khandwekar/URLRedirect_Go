Solution to the excercise as given on:

https://courses.calhoun.io/lessons/les_goph_04

Exercise details

The goal of this exercise is to create an http.Handler that will look at the path of any incoming web request and determine if it should redirect the user to a new page, much like URL shortener would.

Solution Contains:
------------------

- YAML handler
- YAML file paramater from command-line
- CSV handler
- CSV file parameter from command-line
- Database handler, using BOLT Db
- Chanined handlers (in fallback fashion) - 
  Db-Handler -> CSV-handler -> YAML-handler -> Map-handler -> Default-handler 