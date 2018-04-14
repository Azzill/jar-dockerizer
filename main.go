package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"os"
	"path/filepath"
)

var jdkVersion int
var directory string
var jar string
var output string
var port []string
var expose bool
func init(){
	var ports string
	flag.IntVar(&jdkVersion,"jdk",8,"-jdk [jdk version default by 8]")
	flag.StringVar(&directory,"d","","-d [jars' directory]")
	flag.StringVar(&jar,"jar","","-jar [jar file]")
	flag.StringVar(&output,"o",".","-o [tar output directory] default current directory")
	flag.StringVar(&ports,"p","","expose ports eg: -p 1234,1235,1236")
	flag.Parse()
	port = strings.Split(ports,",")
	expose = len(ports) != 0
	if (len(jar) == 0 && len(directory) == 0) || (len(jar) > 0 && len(directory) > 0) {
		panic("you have to specific at least and no more than one jar files with -d or -jar")
	}
}
func main() {
	shell  := ""
	compress := compressUtil{}

	if len(jar) > 0 {
		fileInfo,_ := os.Stat(jar)
		outputPath := strings.TrimRight(output ,"/") + "/" + fileInfo.Name()
		if compress.compressFiles([]string{jar,makeDockerfile(fileInfo.Name())},outputPath){
			fmt.Println("successfully packed "  + jar)
			shell += "cat " + fileInfo.Name()  + ".tar | docker import - " + strings.ToLower(fileInfo.Name()) + "\n"
		}else {
			fmt.Println("failed to pack " + jar)
		}
	} else{ //jars in directory
	filepath.Walk(directory,func(path string, info os.FileInfo, err error) error{
		if strings.Index(path,".jar") == (len(path) - 4){
			fmt.Println("packing " + path )
			outputPath := strings.TrimRight(output ,"/") + "/" + info.Name()
			if compress.compressFiles([]string{path,makeDockerfile(info.Name())},outputPath){
				fmt.Println("successfully packed " + path )
				shell += "cat " + info.Name()  + ".tar | docker import - " + strings.ToLower(info.Name()) + "\n"
			}else {
				fmt.Println("failed to pack "  + info.Name())
			}
		}
		return nil
	})
	}
	fmt.Println("All docker images has been successfully packed!")
	fmt.Println("generating shell script # import.sh")
	err := ioutil.WriteFile(strings.TrimRight(output ,"/") + "/import.sh",[]byte(shell),0666)
	if err == nil{
		fmt.Println("successfully generated shell script!")
	}else {
		fmt.Println("failed to generate shell script.")
	}
	instruction :=
		"Before uploading these tarballs to portainer or using import.sh you need to follow the instructions below\n" +
		"# docker login\n" +
		"# docker pull store/oracle/serverjre:" + string(jdkVersion); fmt.Println (instruction)
	os.Remove(output + "/Dockerfile")
	return
}

func makeDockerfile(file string)string{
	outputPath := strings.TrimRight(output ,"/") + "/Dockerfile"
	content := fmt.Sprintf("FROM store/oracle/serverjre:%d\n" +
								"WORKDIR /app\n" +
								"COPY . /app\n" +
								"ENTRYPOINT [\"java\",\"-jar\",\"%s\"]\n",jdkVersion,file)
	if expose{
	for i := range port{
		content += "EXPOSE " + port[i] + "\n"
	}
	}
	err := ioutil.WriteFile(output + "/Dockerfile",[]byte(content),0666)

	if err != nil{
		panic("failed to write Dockerfile at: " + output + "/Dockerfile")
	}
	return outputPath
}

