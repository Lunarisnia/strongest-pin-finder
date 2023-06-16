package passwordentry

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	ClientID string
	Password string
}

type ClientMonitor struct {
	client map[string]Client
	mu     sync.RWMutex
}

type Client struct {
	counter   int
	onPenalty bool
	password  string
}

var clientMonitor ClientMonitor = ClientMonitor{
	client: map[string]Client{},
	mu:     sync.RWMutex{},
}

func addClientTryCounter(clientID string) (success bool) {
	clientMonitor.mu.Lock()
	defer clientMonitor.mu.Unlock()

	if clientMonitor.client[clientID].counter < 3 {
		clientMonitor.client[clientID] = Client{
			counter:   clientMonitor.client[clientID].counter + 1,
			onPenalty: clientMonitor.client[clientID].onPenalty,
		}
		return true
	} else {
		penalizeClient(clientID)
		return false
	}
}

func penalizeClient(clientID string) {
	if !clientMonitor.client[clientID].onPenalty {
		clientMonitor.client[clientID] = Client{
			counter:   clientMonitor.client[clientID].counter,
			onPenalty: true,
		}
		go func() {
			time.Sleep(5 * time.Second)
			clientMonitor.mu.Lock()
			defer clientMonitor.mu.Unlock()
			clientMonitor.client[clientID] = Client{
				counter:   0,
				onPenalty: false,
			}
		}()
	}
}

// TODO: Make the cracker
func matchPassword(password string, clientID string) bool {
	return clientMonitor.client[clientID].password == password
}

func ProcessPasswordEntry(c *gin.Context) {
	var body LoginRequest
	c.BindJSON(&body)
	fmt.Println(body)
	clientID := body.ClientID
	success := addClientTryCounter(clientID)
	if success {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("%v tried successfully (%v)", clientID, clientMonitor.client[clientID].counter),
			"success": true,
			"guess":   matchPassword(body.Password, clientID),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("%v is on penalty", clientID),
			"success": false,
			"guess":   false,
		})
	}
}
