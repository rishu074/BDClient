package workers

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"

	Conf "github.com/NotRoyadma/BDClient/config"
	"github.com/NotRoyadma/BDClient/logger"
	Stats "github.com/NotRoyadma/BDClient/stats"
	"github.com/gorilla/websocket"
)

func StartUploadWorker() {
	logger.WriteLog("Started Wroker workers/upload.go")
	var addr = flag.String("addr", Conf.Conf.Remote, "http service address")

	// Defer theese things to be safe
	defer logger.WriteLog("Ended Wroker workers/upload.go")
	defer Stats.SetStats(false)

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/api/upload"}
	logger.WriteLog("connecting to " + u.String())

	headers := http.Header{}
	headers.Set("token", Conf.Conf.Token)
	headers.Set("node", Conf.Conf.Node)

	ws, res, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		if err == websocket.ErrBadHandshake {
			p := make([]byte, 200)
			_, _ = res.Body.Read(p)
			logger.WriteLog("workers/upload.go 38 " + err.Error() + string(p))
		}
		logger.WriteERRLog("workers/upload.go 38 " + err.Error())
	}

	logger.WriteLog("connected to " + u.String())
	defer ws.Close()

	type FirstMessageStruct struct {
		Event string
	}

	dataTosend, _ := json.Marshal(FirstMessageStruct{
		Event: "initiate_file",
	})

	ws.WriteMessage(websocket.TextMessage, dataTosend)
	logger.WriteLog("Sent Message to Remote")

	logger.WriteLog("Waiting for Remote...")

	var backFireResponse interface{}

	for {
		ws.SetReadDeadline(time.Now().Add(45 * time.Minute))
		_, strdata, err := ws.ReadMessage()
		if err != nil {
			return
		}
		err = json.Unmarshal(strdata, &backFireResponse)
		if err != nil {
			panic(err)
		}
		break
	}

	logger.WriteLog("Got Data from remote...")
	backFireResponseJson := backFireResponse.(map[string]interface{})

	if backFireResponseJson["Event"] != "initiate_subfolders" {
		logger.WriteERRLog("Unexpected Remote data..")
		logger.WriteERRLog(backFireResponse.(string))
		return
	}

	logger.WriteLog("Preparing subfolders")

	ChunksToSend := backFireResponseJson["Chunk"]
	subfolders, _ := os.ReadDir(Conf.Conf.DataDirectory)

	defer endSharing(ws)

	for _, subfolder := range subfolders {
		if !subfolder.IsDir() {
			return
		}
		logger.WriteLog("Processing subfolder.." + subfolder.Name())

		SubFolderName := subfolder.Name()

		type SeverResponse struct {
			Event string
			Name  string
		}

		dataTosend, _ := json.Marshal(SeverResponse{
			Event: "subfolder_start",
			Name:  SubFolderName,
		})

		ws.WriteMessage(websocket.TextMessage, dataTosend)
		logger.WriteLog("Sent start sb to Remote.")

		// wait for surver
		var backFireResponse interface{}

		for {
			ws.SetReadDeadline(time.Now().Add(45 * time.Minute))
			_, strdata, err := ws.ReadMessage()
			if err != nil {
				return
			}
			err = json.Unmarshal(strdata, &backFireResponse)
			if err != nil {
				panic(err)
			}
			break
		}

		// `subfolder_data_start` event foldername `123`
		if backFireResponse.(map[string]interface{})["Event"].(string) != "subfolder_data_start" {
			logger.WriteERRLog("Unexpected Remote data..")
			logger.WriteERRLog(backFireResponse.(string))
			return
		}

		Filename := backFireResponse.(map[string]interface{})["Filename"].(string)

		ReadingFile, _ := os.Open(Conf.Conf.DataDirectory + "/" + subfolder.Name() + "/" + Filename)
		defer ReadingFile.Close()
		FileInfo, _ := ReadingFile.Stat()

		TotalSize := FileInfo.Size()
		logger.WriteLog("Initializing File transfer..." + strconv.Itoa(int(TotalSize)))
		var CurrentBytes int64 = 0

		type EndMessageStruct struct {
			Event string
		}

		type ChunkMessage struct {
			Event string
			Chunk []byte
		}

		for {
			var NextChunk []byte
			if int(TotalSize)-int(CurrentBytes) < int(ChunksToSend.(float64)) {
				NextChunk = make([]byte, int(int(TotalSize)-int(CurrentBytes)))
			} else {
				NextChunk = make([]byte, int(ChunksToSend.(float64)))
			}

			n, err := ReadingFile.Read(NextChunk)
			if err != nil || n == 0 {
				dataTosend, _ = json.Marshal(EndMessageStruct{
					Event: "end_s_chunk",
				})
				ws.WriteMessage(websocket.TextMessage, dataTosend)
				logger.WriteLog("Send end data to remote endpoint")

				//protect the memory leak
				dataTosend = nil
				NextChunk = nil
				break
			} else if err == io.EOF {
				dataTosend, _ = json.Marshal(EndMessageStruct{
					Event: "end_s_chunk",
				})
				ws.WriteMessage(websocket.TextMessage, dataTosend)
				logger.WriteLog("Send end data to remote endpoint")

				//protect the memory leak
				dataTosend = nil
				NextChunk = nil
				break
			}

			dataTosend, _ = json.Marshal(ChunkMessage{
				Event: "subfolder_chunk_data",
				Chunk: NextChunk,
			})
			logger.WriteLog("Sending data chunk" + strconv.Itoa(len(NextChunk)))
			ws.WriteMessage(websocket.TextMessage, dataTosend)
			CurrentBytes += int64(len(NextChunk))

			//protect the memory
			dataTosend = nil
			NextChunk = nil

			// get confirmance
			var backFireResponse interface{}
			for {
				ws.SetReadDeadline(time.Now().Add(45 * time.Minute))
				_, strdata, err := ws.ReadMessage()
				if err != nil {
					return
				}
				err = json.Unmarshal(strdata, &backFireResponse)
				if err != nil {
					panic(err)
				}
				break
			}

			JsonResponse := backFireResponse.((map[string]interface{}))
			if JsonResponse["Event"].(string) != "subfolder_chunk_data_ack" {
				logger.WriteERRLog("Unexpected Remote data..")
				logger.WriteERRLog(backFireResponse.(string))
				return
			}

			// the server is requesting new chunks means resend them = rerun the loop
		}

	}
}

func endSharing(c *websocket.Conn) {
	type EndMessageStructS struct {
		Event string
	}

	dataTosend, _ := json.Marshal(EndMessageStructS{
		Event: "end_sharing",
	})

	c.WriteMessage(websocket.TextMessage, dataTosend)
	logger.WriteLog("Ended sharing safely")
}
