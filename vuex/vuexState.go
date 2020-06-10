package vuex

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
)

//Store store's content
type Store struct {
	Settings map[string]string `json:"settings"`
}

//State struct
type State struct {
	StorePath string
}

//NewVuexState Create a new VuexState struct
func NewVuexState(storePath string) *State {
	return &State{StorePath: storePath}
}

//homeDir return user's home directory
func homeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

//Write Writes vuex store data to file
func (t *State) Write(data string) bool {
	homeDir, err := homeDir()
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(path.Join(homeDir, t.StorePath+"-vuex.json"))
	if err != nil {
		log.Fatal(err)
	}
	f.Write([]byte(data))
	if err != nil {
		log.Fatal(err)
	}
	if f.Close() != nil {
		log.Fatal(err)
	}
	return true
}

//Read Reads vuex store data from file
func (t *State) Read() string {
	storePath := t.StorePath + "-vuex.json"
	homeDir, err := homeDir()
	if err != nil {
		log.Fatal(err)
	}
	_, err = os.Stat(storePath)
	if os.IsNotExist(err) {
		f, err := os.Create(path.Join(homeDir, storePath))
		f.Write([]byte("{}"))
		if err != nil {
			log.Fatal(err)
		}
		if f.Close() != nil {
			log.Fatal(err)
		}
	}

	data, err := ioutil.ReadFile(path.Join(homeDir, storePath))
	if err != nil {
		log.Fatal(err)
	}

	return string(data)
}
