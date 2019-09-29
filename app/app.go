package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	// color variables (in bytecode form)
	green       = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white       = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow      = string([]byte{27, 91, 57, 48, 59, 52, 51, 109})
	red         = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue        = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta     = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan        = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset       = string([]byte{27, 91, 48, 109})

	// list of status messages (to be sent as response
	// for e.g. [{"status_code" : 0, "response_code" : 200, "message" : "Everything's alright!"},...])
	statMsgData = GetJSONArrayData("status.json")
)

// route handler function for info route
func sendInfo(c echo.Context) error {
	reqM := c.Request()
	resM := c.Response()
	return c.JSON(200, map[string]string{
		"name":        "Children_app_api",
		"developer":   "Yuil Tripathee",
		"version":     "v1.0",
		"status_code": fmt.Sprintf("%d", resM.Status),
		"time":        time.Now().Format("2006/01/01 - 15:04:05"),
		"protocol":    reqM.Proto,
		"ip":          c.RealIP(),
		"method":      reqM.Method,
		"url":         fmt.Sprintf("%s", reqM.URL),
		"bytes_out":   fmt.Sprintf("%d", resM.Size),
		"server_type": "Testing",
	})
}

// route handler to stream image
func sendImage(c echo.Context) error {
	audForID := strings.ToLower(c.Param("id"))
	// sample validation (constraint: parameter must be single character between a - z)
	if len(audForID) > 1 || audForID < "a" || audForID > "z" {
		return c.JSON(500, statMsgData[0])
	}
	audForAddr := fmt.Sprintf("res/image/%s.jpg", audForID)
	return c.File(audForAddr)
}

// route handler for sample DB data
// route handler function for database querying
func sendSampleDBData(c echo.Context) error {
	status, response := getDataDBbyIndex(sampleTable, "id", c.QueryParam("id"))
	return c.JSON(status, response)
}

// route handler function for local JSON data
func sendSampleLocalData(c echo.Context) error {
	status, response_data := 200, GetJSONObjectData("res/sample.json")
	// this system is required only for local JSON (not from database)
	// in database querying system, this is implemented earlier
	response := statMsgData[2]
	response["data"] = response_data
	return c.JSON(status, response)
}

func main() {
	// block to confirm if the runtime enviroment is for DEVELOPMENT or PRODUCTION
	var ENVConfig string
	CLIConfig := os.Args
	
	// Initialization of go-echo server
	e := echo.New()

	// debug mode (optional)
	// e.Debug = true 

	// just to hide the echo framework banner
	// e.HideBanner = true	

	// Adding trailing slash to request URI
	// e.Pre(middleware.AddTrailingSlash())

	// if args supplied (just for logger colorization)
	// also to log out state of back-end (PRODUCTION or DEVELOPMENT)
	if len(CLIConfig) > 1 {
		switch CLIConfig[1] {
		case "DEV":
			ENVConfig = "DEV"
		case "PROD":
			ENVConfig = "PROD"
		default:
			fmt.Println("Invalid arguments supplied.\nExiting the program.")
			os.Exit(0)
		}
	} else {
		// if no args supplied
		ENVConfig = "DEV"
	}
		
	
	// name definition for the runtime application (along with the runtime enviroment variant)
	name := fmt.Sprintf("R&D-%s", ENVConfig)
		
	// coloration block
	// tailored (TBC: colored) logger adapting to the different runtime environment
	switch ENVConfig {
	case "DEV":
		// Debug version of LOG
		e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
			var reqMethod string
			var resStatus int
			var statusColor, methodColor, resetColor string
			// request and response object
			req := c.Request()
			res := c.Response()
			// rendering variables for response status and request method
			resStatus = res.Status
			reqMethod = req.Method
			// for response status
			switch {
			case resStatus >= http.StatusOK && resStatus < http.StatusMultipleChoices:
				statusColor = green
			case resStatus >= http.StatusMultipleChoices && resStatus < http.StatusBadRequest:
				statusColor = white
			case resStatus >= http.StatusBadRequest && resStatus < http.StatusInternalServerError:
				statusColor = yellow
			default:
				statusColor = red
			}
			// for request method
			switch reqMethod {
			case "GET":
				methodColor = blue
			case "POST":
				methodColor = cyan
			case "PUT":
				methodColor = yellow
			case "DELETE":
				methodColor = red
			case "PATCH":
				methodColor = green
			case "HEAD":
				methodColor = magenta
			case "OPTIONS":
				methodColor = white
			default:
				methodColor = reset
			}
			// reset to return to the normal terminal color variables (kinda default)
			resetColor = reset
			// print formatting the custom logger tailored for DEVELOPMENT environment
			fmt.Printf("\n[%s] %v |%s %3d %s| %8s | %10s |%s %-7s %s %s",
				name, // name of server (APP) with the environment
				time.Now().Format("2006/01/02 - 15:04:05"), // TIMESTAMP for route access
				statusColor, resStatus, resetColor, // response status
				req.Proto,                          // protocol
				c.RealIP(),                         // client IP
				methodColor, reqMethod, resetColor, // request method
				req.URL, // request URI (path)
			)
		}))
	default:
		// Production version of LOG
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: fmt.Sprintf("\n%s | ${host} | ${time_custom} | ${status} | ${latency_human} | ${remote_ip} | ${bytes_in} bytes_in | ${bytes_out} bytes_out | ${method} | ${uri} ",
				name,
			),
			CustomTimeFormat: "2006/01/02 15:04:05", // custom readable time format
			Output:           os.Stdout,             // output method
		}))
	}

	// list of endpoint routes
	APIRoute := e.Group("/api")
	// grouping routes for version 1.0 API
	v1route := APIRoute.Group("/v1")
	v1route.GET("/", sendInfo) 	// sample route to send normal JSON from program variable (route for API info)
	v1route.GET("/img/:id", sendImage)                 // sample route to stream image (same apply for music, video and other file types)
	v1route.GET("/source/DB/", sendSampleDBData)       // sample route to demonstrate data transfer from database (MySQL here)
	v1route.GET("/source/local/", sendSampleLocalData) // sample route to demonstrate data transfer from local JSON file

	// static route for dummy landing page
	e.Static("static", "static")

	// stores routes available in the system in a JSON file
	data, err := json.MarshalIndent(e.Routes(), "", "  ")
	if err != nil {
		fmt.Println(err) // error handling in Golang works this way
	}
	ioutil.WriteFile("routes.json", data, 0644)

	// firing up the server
	e.Logger.Fatal(e.Start(":3000"))
}

// reference for generic JSON data format implementation: https://gist.github.com/aogooc/3e81c8752deba00ecab07255d94619fe
// reference for strict JSON data format implementation: https://yayprogramming.com/mysql-to-struct-in-go-language/
