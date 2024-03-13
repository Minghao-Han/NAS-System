package main

import (
	"github.com/gin-gonic/gin"
	"nas/project/src/Utils"
	"nas/project/src/router"
	"sync"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	regularRequestRouter := router.RegularRequestRouter()
	csPorts := Utils.DefaultConfigReader().Get("FSP:csPorts").([]interface{})
	dsPorts := Utils.DefaultConfigReader().Get("FSP:dsPorts").([]interface{})
	wg := sync.WaitGroup{}
	wg.Add(len(csPorts) + len(dsPorts) + 1)
	for _, csPort := range csPorts {
		csPort := csPort
		csRouter := router.GetControlStreamRouter()
		go func() {
			router.RunOnConfig(csPort.(int), csRouter)
			wg.Done()
		}()
	}
	for _, dsPort := range dsPorts {
		dsPort := dsPort
		dsRouter := router.GetDataStreamRouter()
		go func() {
			router.RunOnConfig(dsPort.(int), dsRouter)
			wg.Done()
		}()
	}
	go func() {
		router.RunTLSOnConfig(regularRequestRouter)
		wg.Done()
	}()
	wg.Wait()
}
