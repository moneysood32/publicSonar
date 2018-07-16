package counter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ItemInfo contains the input data fields
type ItemInfo struct {
	Tenant string `json:"tenant"`
	ItemID string `json:"itemID"`
}

// ItemInfoSet is a temporary in-memory store for the Items
var ItemInfoSet = make(map[string][]ItemInfo)
var storageSize = 0
var maxMemoryStorageSize = 3

// ParseAndSave parses the request body and if the data isn't malformed or corrupt, saves it.
func ParseAndSave(data []byte, w http.ResponseWriter, port string) (*ItemInfo, error) {
	var newItem ItemInfo
	if err := json.Unmarshal(data, &newItem); err != nil {
		panic(err)
	}

	items := ItemInfoSet[newItem.Tenant]
	if newItem.FoundIn(items, "../../data/") {
		fmt.Println("Item already exists")
		fmt.Fprintln(w, "Item already exists")
		return nil, nil
	}
	ItemInfoSet[newItem.Tenant] = append(ItemInfoSet[newItem.Tenant], newItem)
	fmt.Println(ItemInfoSet, len(ItemInfoSet))
	storageSize++
	if storageSize >= maxMemoryStorageSize {
		if err := WriteSetToDB(ItemInfoSet, port); err != nil {
			fmt.Println(err)
		}
		storageSize = 0
		ItemInfoSet = make(map[string][]ItemInfo)
	}

	return &newItem, nil
}

// String implements stringer interface over ItemInfo,
// required to print ItemInfo in a user friendly way to the output stream.
func (i ItemInfo) String() string {
	return "Tenant : " + i.Tenant + " , " + "Item : " + i.ItemID
}

// WriteSetToDB writes input map to memory
func WriteSetToDB(userDataset map[string][]ItemInfo, port string) error {
	path := "../../data/" + strings.Join(strings.Split(port, ":"), "")
	if err := os.MkdirAll(path, 0700); err != nil {
		fmt.Println(err)
	}
	fileName := strings.Join(strings.Split(port, ":"), "") + ".txt"
	filePath := path + "/" + fileName
	var fptr *os.File
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fptr, err = os.Create(filePath)
		if err != nil {
			fmt.Println(err)
		}
		_, err = fptr.WriteString("Tenant , ItemID\r\n")
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fptr, err = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}
	}
	if err := fileWriter(fptr, userDataset); err != nil {
		fmt.Println(err)
	}

	fptr.Close()

	return nil
}

func (i ItemInfo) FoundIn(items []ItemInfo, path string) bool {
	foundInItems := i.searchInItems(items)
	foundInPath := i.searchInDB(path)
	return foundInItems || foundInPath
}

func (i ItemInfo) searchInItems(items []ItemInfo) bool {
	for _, item := range items {
		if i.Tenant == item.Tenant && i.ItemID == item.ItemID {
			return true
		}
	}
	return false
}

func (i ItemInfo) searchInDB(path string) bool {
	filePaths := allFilesInDB(path)
	for _, filePath := range filePaths {
		if i.searchInFile(filePath) {
			return true
		}
	}
	return false
}

func (i *ItemInfo) searchInFile(filePath string) bool {
	rawFileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
	}
	fileData := string(rawFileData)
	if strings.Contains(fileData, i.Tenant+" , "+i.ItemID) {
		return true
	}
	return false
}

func allFilesInDB(path string) []string {
	allFiles := make([]string, 0, 10)
	dirs, err := ioutil.ReadDir("../../data")
	if err != nil {
		fmt.Println(err)
	}
	for _, dir := range dirs {
		files, err := filepath.Glob("../../data/" + dir.Name() + "/*")
		if err != nil {
			log.Fatal(err)
		}
		allFiles = append(allFiles, files...)
	}
	return allFiles
}

func fileWriter(fptr *os.File, userDataset map[string][]ItemInfo) error {
	for user, items := range userDataset {
		for _, item := range items {
			fptr.WriteString(user + " , " + item.ItemID + "\r\n")
		}
	}
	return nil
}

func GetItemsFromAllDBs(path string, tenant string) []ItemInfo {
	itemsFromAllDBs := make([]ItemInfo, 0)
	filePaths := allFilesInDB(path)

	for _, filePath := range filePaths {
		itemsFromAllDBs = append(itemsFromAllDBs, getItemsFromDB(filePath, tenant)...)
	}
	return itemsFromAllDBs
}

func getItemsFromDB(filePath string, tenant string) []ItemInfo {
	userItemsInDB := make([]ItemInfo, 0)
	rawFileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
	}
	fileData := string(rawFileData)
	allUserItems := strings.Split(fileData, "\r\n")
	for _, allUserItem := range allUserItems {
		if strings.Contains(allUserItem, tenant) {
			rawItem := strings.Split(allUserItem, " , ")
			item := ItemInfo{
				Tenant: rawItem[0],
				ItemID: rawItem[1],
			}
			userItemsInDB = append(userItemsInDB, item)
		}
	}
	return userItemsInDB
}
