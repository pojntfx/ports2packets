package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
	"sync"
)

func main() {
	// Parse flags
	outpath := flag.String("outpath", "portstopackets.csv", "File to write the ports to")

	flag.Parse()

	// TODO: Use https://www.iana.org/assignments/service-names-port-numbers/service-names-port-numbers.xhtml as CSV for the base
	outFileContent := ""
	listenHost := "localhost"

	listenStartPort := 0
	listenEndPort := 1023

	// Listen on UDP ports and get the packets sent by nmap
	go func() {
		var wg sync.WaitGroup
		for port := listenStartPort; port <= listenEndPort; port++ {
			wg.Add(1)

			go func(wg *sync.WaitGroup, innerPort int) {
				pc, err := net.ListenPacket("udp", fmt.Sprintf("%v:%v", listenHost, innerPort))
				if err != nil {
					log.Fatal(err)
				}
				defer pc.Close()

				buf := make([]byte, 1024)
				n, _, err := pc.ReadFrom(buf)
				if err != nil {
					log.Fatal(err)
				}

				packet := base64.StdEncoding.EncodeToString(buf[:n])

				if packet != "" {
					log.Println(innerPort)

					outFileContent += fmt.Sprintf("%v,%v\n", innerPort, packet)
				}

				wg.Done()

				return
			}(&wg, port)
		}

		wg.Wait()
	}()

	// Run nmap on known UDP ports
	cmd := exec.Command("nmap", "-sU", listenHost)

	output, err := cmd.Output()
	if err != nil {
		log.Fatal(output, err)
	}

	// Write to CSV file
	ioutil.WriteFile(*outpath, []byte(outFileContent), 0644)
}
