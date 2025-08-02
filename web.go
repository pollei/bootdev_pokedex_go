package main

import (
	"io"
	//"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pollei/bootdev_pokedex_go/internal/pokecache"
	//"errors"
)

/*
type pokeLocArea struct {
	name       string
	id         int
	game_index int
} */

// have to have capitilzed to export so unmarshal can reflect on it
type PokeNamedAPIResource struct {
	Name string
	Url  string
}
type PokeNamedAPIResourceList struct {
	Count    int
	Previous *string
	Next     *string
	Results  []PokeNamedAPIResource
	//results []map[string]string
}
type pokeNamedAPIResourceListResults struct {
	baseUrl    string
	linkedList []PokeNamedAPIResourceList
	currIndx   int
}

type PokeEncounter struct {
	Pokemon PokeNamedAPIResource
}
type PokeLocationArea struct {
	Id                 int
	Pokemon_encounters []PokeEncounter
}
type PokemonStat struct {
	Base_stat int
	Effort    int
	Stat      PokeNamedAPIResource
}
type PokemonType struct {
	Slot int
	Type PokeNamedAPIResource
}
type Pokemon struct {
	Name            string
	Id              int
	Base_experience int
	Height          int
	Weight          int
	Stats           []PokemonStat
	Types           []PokemonType
}

var webGLOBS = struct {
	localAreas     PokeNamedAPIResourceList
	localAreasList pokeNamedAPIResourceListResults
	cache          pokecache.Cache
}{}

func getPokeBytes(url string) ([]byte, error) {
	var retBytes []byte
	cacheByte, ok := webGLOBS.cache.Get(url)
	if ok {
		//fmt.Printf("cache hit: %s\n", url)
		return cacheByte, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("bad get")
		return retBytes, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("bad do")
		return retBytes, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if nil != err {
		fmt.Printf("bad readall")
		return retBytes, err
	}
	webGLOBS.cache.Add(url, data)
	// fmt.Println(string(data))
	return data, nil
}

func getPokeResourceList(url string) (PokeNamedAPIResourceList, error) {
	ret := PokeNamedAPIResourceList{}
	locAreaReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("bad get")
		return ret, err
	}
	locAreaReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(locAreaReq)
	if err != nil {
		fmt.Printf("bad do")
		return ret, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if nil != err {
		fmt.Printf("bad readall")
		return ret, err
	}
	//fmt.Println(string(data))
	if err := json.Unmarshal(data, &ret); err != nil {
		fmt.Printf("bad unmarshal")
		return ret, err
	}

	return ret, nil
	//res.Body.Seek(0, io.SeekStart)
	/* decoder := json.NewDecoder(res.Body)

	err = decoder.Decode(&ret)
	if err != nil {
		fmt.Printf("bad decode %v", err)
		return ret, err
	}
	return ret, nil */
}

func getExploreResult(name string) (PokeLocationArea, error) {
	var ret PokeLocationArea
	data, err := getPokeBytes("https://pokeapi.co/api/v2/location-area/" + name)
	if err != nil {
		return ret, err
	}
	if err := json.Unmarshal(data, &ret); err != nil {
		fmt.Printf("bad unmarshal")
		return ret, err
	}
	return ret, nil

}

func getPokemonResult(name string) (Pokemon, error) {
	var ret Pokemon
	data, err := getPokeBytes("https://pokeapi.co/api/v2/pokemon/" + name)
	if err != nil {
		return ret, err
	}
	if err := json.Unmarshal(data, &ret); err != nil {
		fmt.Printf("bad unmarshal")
		return ret, err
	}
	return ret, nil

}

func getNamedResourceResult(rsrcList *pokeNamedAPIResourceListResults, indx int) (PokeNamedAPIResourceList, error) {
	if indx < len(rsrcList.linkedList) {
		return rsrcList.linkedList[indx], nil
	}
	var nmdLst PokeNamedAPIResourceList
	//fmt.Println("getNamedResourceResult before zero indx")
	if indx == 0 {
		nmdLst, err := getPokeResourceList(rsrcList.baseUrl)
		if err != nil {
			return nmdLst, err
		}
		pageCnt := (nmdLst.Count / 20) + 1
		rsrcList.linkedList = make([]PokeNamedAPIResourceList, 0, pageCnt)
		//fmt.Println("getNamedResourceResult after ll make")
		rsrcList.linkedList = append(rsrcList.linkedList, nmdLst)
		//rsrcList.linkedList[0] = nmdLst
		//fmt.Printf(" nmd lst at indx0(%d) %s\n", len(nmdLst.Results), nmdLst.String())
		return nmdLst, nil
	}
	fmt.Println("getNamedResourceResult after zero indx")
	if indx < 0 {
		return nmdLst, fmt.Errorf("negative index not allowed")
	}
	for i := len(rsrcList.linkedList); i <= indx; i++ {
		if rsrcList.linkedList[i-1].Next == nil {
			return nmdLst, fmt.Errorf("can not reach index %d", indx)
		}
		nmdLst, err := getPokeResourceList(*rsrcList.linkedList[i-1].Next)
		if err != nil {
			return nmdLst, err
		}
		rsrcList.linkedList = append(rsrcList.linkedList, nmdLst)
		//rsrcList.linkedList[i] = nmdLst
	}
	fmt.Printf(" nmd lst %s\n", nmdLst.String())
	return nmdLst, nil

}

func (rl PokeNamedAPIResourceList) String() string {
	var sb strings.Builder
	for _, rslt := range rl.Results {
		sb.WriteString(rslt.Name + "\n")
	}
	return sb.String()
}
