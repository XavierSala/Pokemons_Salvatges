package models

// Pokemon defineix un bitxo d'aquests
type Pokemon struct {
	ID    uint32
	Nom   string
	Pes   float32
	Atac  int32
	Vida  int32
	Tipus []TipusPokemon
}
