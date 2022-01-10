package engine

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
	lovelove "hanafuda.moe/lovelove/proto"
)

type loveLoveRpcServer struct {
	lovelove.UnimplementedLoveLoveServer
	testGames map[string]testGame
}

func NewLoveLoveRpcServer() *loveLoveRpcServer {
	return &loveLoveRpcServer{
		UnimplementedLoveLoveServer: lovelove.UnimplementedLoveLoveServer{},
		testGames:                   make(map[string]testGame),
	}
}

type testGameCard struct {
	Variation int
	Month     string
	Hana      string
}

type testGame struct {
	RedHand   []testGameCard `yaml:"redHand"`
	WhiteHand []testGameCard `yaml:"whiteHand"`
	Deck      []testGameCard
	Table     []testGameCard
}

func (server *loveLoveRpcServer) LoadTestGames() {
	data, err := ioutil.ReadFile("test_games.yml")
	if err != nil {
		log.Println("Failed to find test games")
		return
	}

	yaml.Unmarshal(data, &server.testGames)
}
