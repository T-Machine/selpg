package main

import (
	"fmt"
    "io"
    "bufio"
    //"flag"
    "os"
	"os/exec"
	flag "github.com/spf13/pflag"
)

type selpgArgs struct {
	startPage	int
	endPage    int
	length    int
	pageType   bool
	destination  string
	filename string
}

func main(){
	sa := new(selpgArgs)
	getArgs(sa)
	checkArgs(sa)
	processArgs(sa)
}

func getArgs(sa *selpgArgs){
	flag.IntVarP(&(sa.startPage), "start", "s", -1, "the start page")
	flag.IntVarP(&(sa.endPage), "end", "e", -1, "the end page")
	flag.IntVarP(&(sa.length), "length", "l", 72, "the length of a page")
	flag.BoolVarP(&(sa.pageType), "type",  "f", false, "change page at \\f")
	flag.StringVarP(&(sa.destination), "destination", "d", "", "the destination")
	flag.Parse()
	if flag.NArg() > 0 {
		sa.filename = flag.Arg(0)
	} else {
		sa.filename = ""
	}
}

func checkArgs(sa *selpgArgs){
	if sa.startPage == -1 || sa.endPage == -1 {
		fmt.Fprintf(os.Stderr, "Please input start page and end page")
		os.Exit(0)
	}
	if sa.startPage < 1 {
		fmt.Fprintf(os.Stderr, "Invalid start page")
		os.Exit(0)
	}
	if sa.startPage > sa.endPage || sa.endPage < 1 {
		fmt.Fprintf(os.Stderr, "Invalid end page")
		os.Exit(0)
	}
}

func processArgs(sa *selpgArgs){
	//fin := os.Stdin 	//stand input
	fout := os.Stdout	//stand output
	current_page := 1
	current_line := 0
	var reader *bufio.Reader
	var inPipe io.WriteCloser
	var err error

	if sa.filename != "" {
		fin, err := os.Open(sa.filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Open file failed")
			os.Exit(0)
		}
		reader = bufio.NewReader(fin)
		defer fin.Close()
	} else {
		reader = bufio.NewReader(os.Stdin)
	}
	// Input from file

	if sa.destination != "" {
		cmd := exec.Command("lp", "-d" + sa.destination)
		inPipe, err = cmd.StdinPipe()	//get the pipe
		if err != nil || inPipe == nil {
			fmt.Fprintf(os.Stderr, "Failed to open pipe to %s\n", sa.destination)
			os.Exit(0)
		}
		cmd.Stdout = fout
		//cmd.Start()
		err := cmd.Run()
		if err != nil {
			os.Stderr.Write([]byte("no printer\n"))
		}
		defer inPipe.Close()
	}
	// -d	print to destination

	
	if sa.pageType {
		for {
			one_page, err := reader.ReadString('\f')

			if err != nil && err != io.EOF {
				fmt.Fprintf(os.Stderr, "Failed to read a page\n")
				os.Exit(0)
			}

			if   current_page >= sa.startPage && current_page <= sa.endPage{
				if sa.destination != "" {
					fmt.Fprintf(inPipe, "%s", one_page)
				} else {
					fmt.Fprintf(fout, "%s", one_page)
				}
			}
			//output to stdout or pipe
			if err == io.EOF || current_page > sa.endPage {
				break
			}
			current_page ++
		}
	} else {
		for {
			one_line, err := reader.ReadString('\n')

			if err != nil && err != io.EOF {
				fmt.Fprintf(os.Stderr, "Failed to read a line\n")
				os.Exit(0)
			}

			if   current_page >= sa.startPage && current_page <= sa.endPage{
				if sa.destination != "" {
					fmt.Fprintf(inPipe, "%s", one_line)
				} else {
					fmt.Fprintf(fout, "%s", one_line)
				}
			}
			//output to stdout or pipe
			current_line ++
			if current_line >= sa.length {
				current_line = 0
				current_page ++
			}
			if err == io.EOF || current_page > sa.endPage {
				break
			}
		}
	}// -l
	if current_page < sa.startPage {
		fmt.Fprintf(os.Stderr, "The start page is greater than the total pages\n")
	} else if current_page < sa.endPage {
		fmt.Fprintf(os.Stderr, "The end page is greater than the total pages\n")
	}

}