package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

var livez bool = false
var readyz bool = false
var livenessCount int = 0
var readinessCount int = 0

func livenessDelay(delay int, fileName string) {
	log.Print("Starting livez delay...")
	time.Sleep(time.Duration(delay) * time.Second)
	livez = true
	writeFile(fileName)
	log.Print("Completed livez delay")
}

func readinessDelay(delay int, fileName string) {
	log.Print("Starting readyz delay...")
	time.Sleep(time.Duration(delay) * time.Second)
	readyz = true
	writeFile(fileName)
	log.Print("Completed readyz delay")
}

func removeFile(fileName string) {
	log.Print("Removing file: " + fileName)
	if _, err := os.Stat(fileName); err == nil {
		e := os.Remove(fileName)
		if e != nil {
			 log.Fatal(e)
		}
	}
}

func writeFile(fileName string) {
	log.Print("Writing file: " + fileName)
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		file, err := os.Create(fileName)
		if err != nil {
				log.Fatal(err)
		}
		defer file.Close()
	}
}

func main() {
	log.Print("Starting the server...")

	port := os.Getenv("PORT")
	startupFile := os.Getenv("STARTUP_FILE")
	livenessFile := os.Getenv("LIVENESS_FILE")
	readinessFile := os.Getenv("READINESS_FILE")
	listenDelay := os.Getenv("LISTEN_DELAY_SECONDS")
	livezDelay := os.Getenv("LIVENESS_DELAY_SECONDS")
	readyzDelay := os.Getenv("READINESS_DELAY_SECONDS")
	responseDelay := os.Getenv("RESPONSE_DELAY_MILLISECONDS")
	livenessSuccessMax := os.Getenv("LIVENESS_SUCCESS_MAX")
	readinessSuccessMax := os.Getenv("READINESS_SUCCESS_MAX")

	if port == "" {
		port = "8000"
		log.Print("Using default port 8000")
	} else {
		log.Print("Using port " + port)
	}

	if startupFile == "" {
		startupFile = "/tmp/startup"
		log.Print("Using default startupFile /tmp/startup")
	} else {
		log.Print("Using startupFile " + startupFile)
	}

	if livenessFile == "" {
		livenessFile = "/tmp/liveness"
		log.Print("Using default livenessFile /tmp/liveness")
	} else {
		log.Print("Using livenessFile " + livenessFile)
	}

	if readinessFile == "" {
		readinessFile = "/tmp/readiness"
		log.Print("Using default readinessFile /tmp/readiness")
	} else {
		log.Print("Using readinessFile " + readinessFile)
	}

	removeFile(startupFile)
	removeFile(livenessFile)
	removeFile(readinessFile)

	if listenDelay == "" {
		listenDelay = "10"
		log.Print("Using listen delay default of 10s")
	} else {
		log.Print("Using listen delay " + listenDelay + "s")
	}
	listenDelaySeconds, err := strconv.Atoi(listenDelay)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to convert listenDelay"))
	}

	if livezDelay == "" {
		livezDelay = "2"
		log.Print("Using live delay default of 2s")
	} else {
		log.Print("Using live delay " + livezDelay + "s")
	}
	livezDelaySeconds, err := strconv.Atoi(livezDelay)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to convert livezDelay"))
	}

	if readyzDelay == "" {
		readyzDelay = "10"
		log.Print("Using readiness delay default of 10s")
	} else {
		log.Print("Using readiness delay " + readyzDelay + "s")
	}
	readyzDelaySeconds, err := strconv.Atoi(readyzDelay)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to convert readyzDelay"))
	}

	if responseDelay == "" {
		responseDelay = "0"
		log.Print("Using response delay default of 0")
	} else {
		log.Print("Using response delay " + responseDelay + "ms")
	}
	responseDelayMilliSeconds, err := strconv.Atoi(responseDelay)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to convert responseDelay"))
	}

	if livenessSuccessMax == "" {
		livenessSuccessMax = "0"
		log.Print("Using liveness success max default of 0")
	} else {
		log.Print("Using liveness success max " + livenessSuccessMax)
	}
	livenessCountMax, err := strconv.Atoi(livenessSuccessMax)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to convert livenessSuccessMax"))
	}

	if readinessSuccessMax == "" {
		readinessSuccessMax = "0"
		log.Print("Using readiness success max default of 0")
	} else {
		log.Print("Using readiness success max " + readinessSuccessMax)
	}
	readinessCountMax, err := strconv.Atoi(readinessSuccessMax)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to convert readinessSuccessMax"))
	}

	if listenDelaySeconds > 0 {
		log.Print("Starting listen delay...")
		time.Sleep(time.Duration(listenDelaySeconds) * time.Second)
		log.Print("Completed listen delay")
	} else {
		log.Print("No listen delay")
	}
	writeFile(startupFile)

	go livenessDelay(livezDelaySeconds, livenessFile)
	go readinessDelay(readyzDelaySeconds, readinessFile)

	http.HandleFunc("/home", func(w http.ResponseWriter, _ *http.Request) {
		if responseDelayMilliSeconds != 0 {
			time.Sleep(time.Duration(responseDelayMilliSeconds) * time.Millisecond)
		}
		if readyz {
			log.Print("/home request when ready")
			fmt.Fprint(w, "/home request processed")
		} else {
			log.Print("/home request when not ready")
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	},
	)

	http.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		if responseDelayMilliSeconds != 0 {
			time.Sleep(time.Duration(responseDelayMilliSeconds) * time.Millisecond)
		}
		if readyz {
			if readinessCountMax != 0 {
				readinessCount++
				if readinessCount > readinessCountMax {
					log.Print("/readyz request after readiness success count exceeded " + strconv.Itoa(readinessCount) + "/" + strconv.Itoa(readinessCountMax))
					w.WriteHeader(http.StatusServiceUnavailable)
				} else {
					log.Print("/readyz request when ready " + strconv.Itoa(readinessCount) + "/" + strconv.Itoa(readinessCountMax))
					fmt.Fprint(w, "/readyz request processed")
				}
			} else {
				log.Print("/readyz request when ready")
				fmt.Fprint(w, "/readyz request processed")
			}
		} else {
			log.Print("/readyz request when not ready")
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	},
	)

	http.HandleFunc("/livez", func(w http.ResponseWriter, _ *http.Request) {
		if responseDelayMilliSeconds != 0 {
			time.Sleep(time.Duration(responseDelayMilliSeconds) * time.Millisecond)
		}
		if livez {
			if livenessCountMax != 0 {
				livenessCount++
				if livenessCount > livenessCountMax {
					log.Print("/livez request after liveness success count exceeded " + strconv.Itoa(livenessCount) + "/" + strconv.Itoa(livenessCountMax))
					w.WriteHeader(http.StatusServiceUnavailable)
				} else {
					log.Print("/livez request when ready " + strconv.Itoa(livenessCount) + "/" + strconv.Itoa(livenessCountMax))
					fmt.Fprint(w, "/livez request processed")
				}
			} else {
				log.Print("/livez request when ready")
				fmt.Fprint(w, "/livez request processed")
			}
		} else {
			log.Print("/livez request when not ready")
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	},
	)

	http.HandleFunc("/crash", func(w http.ResponseWriter, _ *http.Request) {
		if responseDelayMilliSeconds != 0 {
			time.Sleep(time.Duration(responseDelayMilliSeconds) * time.Millisecond)
		}
		if readyz {
			log.Print("/crash request when ready")
			fmt.Fprint(w, "/crash request processed")
		} else {
			log.Print("/crash request when not ready")
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		log.Fatal("/crash endpoint received a request, crashing...")
	},
	)

	log.Print("The service is listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
