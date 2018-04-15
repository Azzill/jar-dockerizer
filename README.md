# jar-dockerizer
  ## a dockerize tool for java8
## Installation
      # go.exe build -i -o dockerize main.go compress.go
## Usage 
  #### -d string 
      -d [jars' directory]
  #### -jar string
      -jar [jar file]
  #### -jdk int
      -jdk [jdk version default by 8] (only support 8)
  #### -o string
      -o [tar output directory] default current directory (default ".")
  #### -p string
      expose ports eg: -p 1234,1235,1236
## Example
      # dockerize -jar myjar.jar
      # dockerize -jar myjar.jar -p 80,8080
      # dockerize -jar myjar.jar -o dir
      # dockerize -d /var/path -o /var/path
## What will be built in the output directory
#### the tarballs compressed with Dockerfile and jar
#### the "build.sh" for docker to import these tarballs
## Download base image in Docker
#### Download jre8 base image from https://store.docker.com/images/oracle-serverjre-8
#### You might need to sign up an account
## Build images
      # ./build.sh



