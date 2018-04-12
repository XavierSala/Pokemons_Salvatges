package main

import (
	"fmt"
	"sort"

	db "./pokemon/database"
	l "./pokemon/lluita"
	m "./pokemon/models"
)

func main() {
	var connexio db.BaseDeDades

	_, err := connexio.Connecta()
	checkErr(err)

	// Tria els 2 equips de 2 jugadors
	pokemons, err := connexio.TriaPokemons(4, 200)
	checkErr(err)

	tmp := make([][]m.Pokemon, 2)
	tmp[0] = pokemons[:2]
	tmp[1] = pokemons[2:]

	combat := l.Lluita{Equips: tmp}

	fmt.Println("-----------------")
	fmt.Println("Lluita entre:")
	fmt.Print(combat.String())
	fmt.Println("-----------------")

	combat.Combat(connexio)

	fmt.Println("-----------------")
	fmt.Println("Resultat:")
	fmt.Println(combat.String())
	fmt.Println("-----------------")

	connexio.Desconnecta()

}

func ordenaPerPes(llista []m.Pokemon) []m.Pokemon {
	sort.Slice(llista,
		func(i, j int) bool { return llista[i].Pes < llista[j].Pes })
	return llista
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
