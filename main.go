package main

import (
	"encoding/base64"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

func main() {
	// Parse flags
	inFilePath := flag.String("in", "/etc/ports2packets/service-names-port-numbers.csv", "Path to the CSV input file containing the registered services. Download from https://www.iana.org/assignments/service-names-port-numbers/service-names-port-numbers.xhtml")
	outFilePath := flag.String("out", "ports2packets.csv", "Path of the generated CSV output file")

	flag.Parse()

	// Read registered ports from CSV file
	csvFile, err := os.Open(*inFilePath)
	if err != nil {
		log.Fatal(err)
	}

	registeredPorts := make(map[int]bool)

	csvReader := csv.NewReader(csvFile)
	currentLine := 0
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		// Skip header
		if currentLine != 0 {
			// We only care about UDP ports, as we can use net.Dial to check if TCP ports are open
			if record[2] == "udp" {
				if record[1] != "" {
					rangeEnd := strings.Split(record[1], "-")

					startPort, err := strconv.Atoi(rangeEnd[0])
					if err != nil {
						log.Fatal(err)
					}

					if !registeredPorts[startPort] {
						registeredPorts[startPort] = true
					}

					if len(rangeEnd) > 1 {
						endPort, err := strconv.Atoi(rangeEnd[1])
						if err != nil {
							log.Fatal(err)
						}

						delta := endPort - startPort

						for i := 1; i <= delta; i++ {
							if !registeredPorts[startPort+i] {
								registeredPorts[startPort+i] = true
							}
						}
					}
				}
			}
		}

		currentLine++
	}

	// Listen on UDP ports and get the packets sent by nmap
	outFileContent := "port,packet\n"
	listenHost := "localhost"

	go func() {
		var wg sync.WaitGroup
		for port := range registeredPorts {
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
					outFileContent += fmt.Sprintf("%v,%v\n", innerPort, packet)
				}

				wg.Done()

				return
			}(&wg, port)
		}

		wg.Wait()
	}()

	// Run nmap on known UDP ports
	cmd := exec.Command("nmap", "–host-timeout", "2000", "-sU", listenHost)

	output, err := cmd.Output()
	if err != nil {
		log.Fatal(output, err)
	}

	// Write to CSV file
	if err := ioutil.WriteFile(*outFilePath, []byte(outFileContent), 0644); err != nil {
		log.Fatal(err)
	}
}
