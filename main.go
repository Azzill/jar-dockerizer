package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"os"
	"path/filepath"
)

var jdkVersion string
var directory string
var jar string
var output string
var port []string
var expose bool
const (
	header = 1
	body = 2
	footer = 3
)
func init(){
	var ports string
	flag.StringVar(&jdkVersion,"jdk","8","-jdk [jdk version default by 8]")
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
	writeShellScript(&shell,"",header)
	compress := compressUtil{}
	outputPath := strings.TrimRight(output ,"/") +  "/"
	if len(jar) > 0 {
		fileInfo,_ := os.Stat(jar)
		if compress.compressFiles([]string{jar,makeDockerfile(fileInfo.Name())},outputPath + fileInfo.Name() + ".tar"){
			fmt.Println("successfully packed "  + jar)
			writeShellScript(&shell,fileInfo.Name(),body)
		}else {
			fmt.Println("failed to pack " + jar)
		}
	} else{
	filepath.Walk(directory,func(path string, info os.FileInfo, err error) error{
		if strings.Index(path,".jar") == (len(path) - 4){
			fmt.Println("packing " + path )
			if compress.compressFiles([]string{path,makeDockerfile(info.Name())},outputPath + info.Name() + ".tar"){
				fmt.Println("successfully packed " + path )
				writeShellScript(&shell,info.Name(),body)
			}else {
				fmt.Println("failed to pack "  + info.Name())
			}
		}
		return nil
	})
	}
	fmt.Println("All docker images has been successfully packed!")
	fmt.Println("generating shell script # build.sh")
	writeShellScript(&shell,"",footer)
	err := ioutil.WriteFile(strings.TrimRight(output ,"/") + "/build.sh",[]byte(shell),0666)
	if err == nil{
		fmt.Println("successfully generated shell script!")
	}else {
		fmt.Println("failed to generate shell script.")
	}
	os.Remove(output + "/Dockerfile")
	return
}

func makeDockerfile(file string)string{
	outputPath := strings.TrimRight(output ,"/") + "/Dockerfile"
	content := fmt.Sprintf("FROM openjdk:%s\n" +
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

func writeShellScript(content *string,tar string,pos int){
	switch pos {
	case header:
		*content += "docker pull openjdk:" + jdkVersion + "\nmkdir -p dockerizer_tmp\n"
	case body:
		*content += fmt.Sprintf("tar xvf %s.tar -C dockerizer_tmp\ndocker build -t %s dockerizer_tmp\nrm dockerizer_tmp/%s\nrm dockerizer_tmp/Dockerfile\n",tar,strings.ToLower(tar),tar)
	case footer:
		*content += "rm -r dockerizer_tmp"
	}
}

