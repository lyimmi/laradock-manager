package vuex

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
)

//Store is vuex store's content
type Store struct {
	Containers struct {
		Favorites           []string `json:"favoritContainers"`
		AvailableContainers []struct {
			Name     string `json:"name"`
			State    string `json:"state"`
			Favorite bool   `json:"favorite"`
		} `json:"availableContainers"`
	} `json:"Containers"`
	Settings struct {
		LaradockPath    string `json:"laradockPath"`
		TerminalPath    string `json:"terminalPath"`
		ContainerPrefix string `json:"containerPrefix"`
		DarkTheme       bool   `json:"darkTheme"`
	} `json:"Settings"`
	Status struct {
		Env           string `json:"env"`
		Docker        string `json:"docker"`
		DockerCompose string `json:"dockerCompose"`
	} `json:"Status"`
}

//State struct
type State struct {
	Store     Store
	StorePath string
}

//NewVuexState Create a new VuexState struct
func NewVuexState(storePath string) *State {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return &State{StorePath: path.Join(usr.HomeDir, storePath+"-vuex.json")}
}

//Write Writes vuex store data to file
func (t *State) Write(data string) {
	f, err := os.Create(t.StorePath)
	if err != nil {
		log.Fatal(err)
	}
	byteData := []byte(data)
	f.Write(byteData)
	if err != nil {
		log.Fatal(err)
	}
	if f.Close() != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(byteData, &t.Store)
	if err != nil {
		log.Fatal(err)
	}
}

//Read Reads vuex store data from file
func (t *State) Read() string {
	_, err := os.Stat(t.StorePath)
	if os.IsNotExist(err) {
		f, err := os.Create(t.StorePath)
		f.Write([]byte("{}"))
		if err != nil {
			log.Fatal(err)
		}
		if f.Close() != nil {
			log.Fatal(err)
		}
	}

	data, err := ioutil.ReadFile(t.StorePath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &t.Store)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

//ClearSettings delete vuex setting file
func (t *State) ClearSettings() bool {
	err := os.Remove(t.StorePath)
	if err != nil {
		log.Fatal(err)
	}
	return true
}
