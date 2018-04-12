package lluita

import (
	"errors"
	"fmt"
	"log"

	b "../database"

	m "../models"
)

// Lluita entre dos equips de Pokemon
type Lluita struct {
	Equips [][]m.Pokemon
}

// String serveix per imprimir els resultats
func (l *Lluita) String() string {
	resultat := ""
	for i, equip := range l.Equips {
		resultat += fmt.Sprintf("Equip %d:\n", i)
		for _, p := range equip {
			resultat += fmt.Sprintf("  ->  %s (%d)\n", p.Nom, p.Atac)
		}
	}
	return resultat
}

func (l *Lluita) purga() {

	for i, t := range l.Equips {
		tmp := make([]m.Pokemon, 0)
		for _, p := range t {
			if p.Vida > 0 {
				tmp = append(tmp, p)
			} else {
				log.Printf("ELIMINAT %s (equip %d)", p.Nom, i)
			}
		}
		l.Equips[i] = tmp
	}

}

// TriaQuiComenca determina quin és l'equip que ataca
func (l *Lluita) triaQuiComenca() int {
	if l.Equips[0][0].Pes > l.Equips[1][0].Pes {
		return 1
	}
	return 0
}

// TriaJugador tria un dels jugadors de l'equip
func (l *Lluita) triaCombatent(equip []m.Pokemon) (int, error) {
	for i, pokemon := range equip {
		if pokemon.Vida >= 0 {
			return i, nil
		}
	}
	return -1, errors.New("No queden Pokemons")
}

// Combat Pokemon
func (l *Lluita) Combat(db b.BaseDeDades) (bool, error) {

	qui := l.triaQuiComenca()

	for len(l.Equips[0]) > 0 && len(l.Equips[1]) > 0 {

		combatents := make([]int, 0)

		for i := range l.Equips {
			quinPokemon, err := l.triaCombatent(l.Equips[i])
			if err == nil {
				combatents = append(combatents, quinPokemon)
			}
		}

		if len(combatents) == 2 {
			qui = l.batalla(qui, combatents, db)
			l.purga()
		}
	}
	return true, nil
}

// Batalla entre els dos Pokemon que ataquen
func (l *Lluita) batalla(comenca int, combatents []int, bd b.BaseDeDades) int {

	var ForcaAtac [2]int32
	//  Localitza com canvien els atacs
	ForcaAtac[0] = int32(float32(l.Equips[0][combatents[0]].Atac) * bd.ModificaAtac(l.Equips[0][combatents[0]], l.Equips[1][combatents[1]]))
	ForcaAtac[1] = int32(float32(l.Equips[1][combatents[1]].Atac) * bd.ModificaAtac(l.Equips[1][combatents[1]], l.Equips[0][combatents[0]]))

	// Eliminar els Pokemons si no poden guanyar perquè no fan res a l'altre els elimino
	if ForcaAtac[0] == 0 && ForcaAtac[1] == 0 {
		log.Println("... Els combatents es moren d'avorriment perquè no es fan res")
		l.Equips[0][combatents[0]].Vida = 0
		l.Equips[1][combatents[1]].Vida = 0
		l.purga()
	} else {

		for l.Equips[0][combatents[0]].Vida > 0 && l.Equips[1][combatents[1]].Vida > 0 {
			atacat := (comenca + 1) % 2

			ataca := combatents[comenca]
			defensa := combatents[atacat]

			// l.Equips[atacat][ataca].Vida -= l.Equips[comenca][defensa].Atac
			l.Equips[atacat][ataca].Vida -= ForcaAtac[comenca]

			log.Printf("%s(%d) ataca %d a %s(%d)",
				l.Equips[comenca][ataca].Nom,
				l.Equips[comenca][ataca].Vida,
				ForcaAtac[comenca],
				l.Equips[atacat][defensa].Nom,
				l.Equips[atacat][defensa].Vida)

			comenca = (comenca + 1) % 2
		}
	}
	return comenca
}
