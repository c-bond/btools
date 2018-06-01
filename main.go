package main

//grpc "btools/grpc"

const on bool = true

func main() {

	if on {
		// {
		// 	cmd := exec.Command("wine", os.Getenv("HOME")+"/.wine/drive_c/Program Files (x86)/Bentley/ProjectWise/bin/helloworld.exe")
		// 	err := cmd.Run()
		// 	log.Println(err)
		// }
		//mgr.StartManager()
		StartServer()
	}
}
