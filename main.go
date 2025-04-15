package main

import (
	"fmt"
	"time"

	"gochopit/apicalls"
	"gochopit/logger"
	"gochopit/services"
	"gochopit/utils"
)

func main() {
	// Initialisation du logger
	logger.InitLogger()

	parsedUrl, err := services.ParsedUrl(utils.Link)
	if err != nil {
		fmt.Println("Erreur : ", err)
		return
	}

	CatalogApi := services.SetCatalogApi(parsedUrl)

	cookies := apicalls.GetCookie(CatalogApi)
	fcookies := apicalls.FormatedAuthCookie(cookies)
	//fmt.Printf("Cookies : %s; %s \n", fcookies.Accesstoken, fcookies.Refreshtoken)
	//fmt.Println(cookies)
	//items, err := apicalls.FetchCatalogItems(CatalogApi, fcookies)
	//if err != nil {
	//	fmt.Print(err)
	//}
	//fmt.Println(items)
	//fmt.Println(len(items))

	go services.FetchCatalogAtInterval(fcookies, parsedUrl, 5*time.Second)
	services.Fetchclean()
}
