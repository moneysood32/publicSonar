package requestProcessor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"publicSonarAssignment/src/counter"
	"strings"
)

func HandleGetRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		arr := strings.Split(r.URL.Path, "/")
		tenant := arr[2]
		itemsFromDBs := counter.GetItemsFromAllDBs("../../data/", tenant)
		tenantItems := append(counter.ItemInfoSet[tenant], itemsFromDBs...)
		itemsJSON, err := json.Marshal(tenantItems)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprintln(w, "count : "+string(len(tenantItems)+48)+" Data : "+string(itemsJSON))
	default:
		fmt.Fprintf(w, "invalid URL for GET request, try http://localhost:3001/items/{tenant}/count")
	}
}

func HandlePostRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		ItemInfo, err := counter.ParseAndSave(body, w, r.Host)
		if err != nil {
			panic(err)
		}
		fmt.Fprintln(w, ItemInfo)
	default:
		fmt.Fprintf(w, "invalid URL for POST request, try http://localhost:3001/items/")
	}
}
