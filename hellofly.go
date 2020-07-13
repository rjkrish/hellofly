package main

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

/*Call is a call record*/
type Call struct {
	ID            int
	Name          string
	Email         string
	Phone         string
	RequestedTime time.Time
	Attempts      int8
	Completed     bool
	CompletedTime time.Time
}

/*Queue holds all the received call requests
TODO: make this a priority queue based on created time
*/
var Queue = make([]*Call, 0, 10)
var cid = 100000
var qlock = sync.Mutex{}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	//for testing
	addToQueue(&Call{Name: "Alice", Email: "", Phone: "6505061000"})

	r.GET("/", handleIndex)
	r.POST("/calls", handleCall)
	r.GET("/calls", handleQueue)
	r.POST("/calls/:cid/complete", handleCallCompletion)
	r.Run()
}

func handleIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}

func handleQueue(c *gin.Context) {
	c.HTML(http.StatusOK, "queue.tmpl", gin.H{"Q": Queue})
}

func handleCall(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	phone := c.PostForm("phone")
	log.Println("Incoming call request from", name)
	addToQueue(&Call{Name: name, Email: email, Phone: phone})
	c.HTML(http.StatusOK, "call.tmpl", gin.H{"Name": name})
}

func handleCallCompletion(c *gin.Context) {
	i, e := strconv.Atoi(c.Param("cid"))
	if e != nil {
		log.Println("Error:", e)
	} else {
		removeFromQueue(i)
	}
	handleQueue(c)
}

func addToQueue(c *Call) {
	qlock.Lock()

	c.RequestedTime = time.Now()
	c.Attempts = 0
	cid = cid + 1
	c.ID = cid
	Queue = append(Queue, c)

	qlock.Unlock()
}

func removeFromQueue(cid int) {
	qlock.Lock()
	for c := range Queue {
		if Queue[c].ID == cid {
			Queue[c].CompletedTime = time.Now()
			Queue[c].Completed = true
		}
	}

	qlock.Unlock()
}
