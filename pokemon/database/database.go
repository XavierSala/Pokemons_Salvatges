package database

import (
	"database/sql"
	"fmt"

	m "../models"
	// Driver MySQL
	_ "github.com/go-sql-driver/mysql"
)

const (
	host     = "127.0.0.1"
	port     = 3306
	user     = "root"
	password = "ies2010"
	dbname   = "pokemons"
)

const cercaLluitadors = `SELECT p.pokemon_id, p.nom, p.pes, pp.VALOR as atac
FROM POKEMONS p
INNER JOIN POKEMON_PODER pp ON p.pokemon_id = pp.pokemon_id
INNER JOIN PODERS po ON pp.PODER_ID = po.PODER_ID
WHERE po.NOM = 'Atac'
ORDER BY RAND() LIMIT ?`

const cercaTipus = `SELECT pt.tipus_id, t.nom
FROM POKETIPUS pt
INNER JOIN TIPUS t ON pt.TIPUS_ID = t.TIPUS_ID
WHERE pt.POKEMON_ID = ?`

const cercaModificadorAtac = `SELECT EFECTE
FROM TIPUS_ATAC
WHERE TIPUS_ATACANT_ID=? AND TIPUS_ATACAT_ID=?`

// BaseDeDades és la interfície bàsica
type BaseDeDades struct {
	db *sql.DB
}

// Connecta amb la base de dades
func (b *BaseDeDades) Connecta() (bool, error) {
	var err error

	myInfo := fmt.Sprintf("%s:%s@/%s", user, password, dbname)

	b.db, err = sql.Open("mysql", myInfo)
	if err != nil {
		return false, err
	}

	// Es fa servir Ping perquè Open no crea la connexió, només en
	// comprova els paràmetres
	err = b.db.Ping()
	if err != nil {
		panic(err)
	}

	return true, nil
}

// Desconnecta de la base de dades
func (b *BaseDeDades) Desconnecta() {
	b.db.Close()
}

// TriaPokemons escul el número de Pokemons que s'hi especifiquen
func (b *BaseDeDades) TriaPokemons(numero int, vida int32) ([]m.Pokemon, error) {
	var pokemons []m.Pokemon
	rows, err := b.db.Query(cercaLluitadors, 4)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var pokemon m.Pokemon
		err = rows.Scan(&pokemon.ID, &pokemon.Nom, &pokemon.Pes, &pokemon.Atac)
		if err != nil {
			return nil, err
		}
		pokemon.Vida = vida
		pokemon.Tipus, err = b.localitzaTipus(pokemon.ID)
		if err != nil {
			return nil, err
		}
		pokemons = append(pokemons, pokemon)
	}
	return pokemons, nil
}

func (b *BaseDeDades) localitzaTipus(id uint32) ([]m.TipusPokemon, error) {
	var tipus []m.TipusPokemon
	rows, err := b.db.Query(cercaTipus, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tip m.TipusPokemon
		err = rows.Scan(&tip.ID, &tip.Nom)
		if err != nil {
			return nil, err
		}
		tipus = append(tipus, tip)
	}
	return tipus, nil
}

// ModificaAtac busca com s'ha de modificar l'atac
func (b *BaseDeDades) ModificaAtac(ataca, atacat m.Pokemon) float32 {

	forca := float32(1.0)

	for _, poderAtac := range ataca.Tipus {
		for _, poderDefensa := range atacat.Tipus {
			row := b.db.QueryRow(cercaModificadorAtac, poderAtac.ID, poderDefensa.ID)
			var mod float32
			err := row.Scan(&mod)
			if err == nil {
				forca = forca * mod
			}
		}

	}
	return forca
}
