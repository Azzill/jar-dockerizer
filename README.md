# jar-dockerizer
## This tool generates tarballs and shell script for docker to import images easier
## Usage of dockerize: 
  ### -d string 
        -d [jars' directory]
  ### -jar string
        -jar [jar file]
  ### -jdk int
        -jdk [jdk version default by 8] (default 8)
  ### -o string
        -o [tar output directory] default current directory (default ".")
  ### -p string
        expose ports eg: -p 1234,1235,1236