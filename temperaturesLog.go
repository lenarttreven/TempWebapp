package main

import (
        "log"
        "os"
)


func writeLog(status string, description string){
         //create file with desired read/write permissions
        var name string
        switch status{
        case "OK":
                name = "status.log"
        case "Error":
                name = "error.log"
        }
        file, err := os.OpenFile("../../log/Page/" + name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
        if err!= nil{
                        log.Fatal(err)
                }
        //defer to close when you're done with it
        defer file.Close()
         //set output of logs to file
        
        log.SetOutput(file)
        log.Println(description)
}
