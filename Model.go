package main

import (
	"log"

	"github.com/lxn/walk"
)

// Type représentant un élément d'une liste avec sa valeur et son nom
type Item struct {
	name  string
	value string
}

// Correspond à un model pour une liste
type Model struct {
	walk.ListModelBase
	Items []Item
}

// Création du modele de la liste des profils
func NewProfilModel() *Model {
	initClient()
	profiles := getAllprofils()
	m := &Model{Items: make([]Item, len(profiles))}

	for i, profile := range profiles {
		if profile.Code == "" {
			m.Items[i] = Item{name: profile.Nom, value: profile.Nom}
		} else {
			m.Items[i] = Item{name: profile.Nom + " / " + profile.Code, value: profile.Nom}
		}

	}

	return m
}

// Création du modele de la liste des fonctions
func NewFonctionModel() *Model {
	profiles := getAllFonctions()
	m := &Model{Items: make([]Item, len(profiles))}

	for i, profile := range profiles {
		if profile.Code == "" {
			m.Items[i] = Item{name: profile.Nom, value: profile.Nom}
		} else {
			m.Items[i] = Item{name: profile.Nom + " / " + profile.Code, value: profile.Nom}
		}

	}

	return m
}

// Création d'un model en fonction du nom du profil pour y lister les fonctions dont il dispose
func NewFonctionUpdateModel(p1 string) *Model {
	itemprofil, err := getProfilByName(p1)
	if err != nil {
		log.Println("Error retrieving profile:", err)
		return &Model{Items: nil} // Return an empty model on error
	}

	Asso, err := getFunctionsByProfil(itemprofil)
	if err != nil {
		log.Println("Error retrieving functions for profile:", err)
		return &Model{Items: nil} // Return an empty model on error
	}

	// Populate the model with FonctionItems
	items := make([]Item, len(Asso.Fonctions))
	for i, fonction := range Asso.Fonctions {
		name := fonction.Nom
		if fonction.Code != "" {
			name += " / " + fonction.Code
		}
		items[i] = Item{name: name, value: fonction.Nom}
	}

	return &Model{Items: items}
}

// Création d'un model en fonction du nom du profil pour y lister les sous - fonctions dont il dispose
func NewSubFonctionUpdateModel(p1 string) *Model {
	itemprofil, err := getProfilByName(p1)
	if err != nil {
		log.Println("Error retrieving profile:", err)
		return &Model{Items: nil} // Return an empty model on error
	}

	Asso, err := getFunctionsByProfil(itemprofil)
	if err != nil {
		log.Println("Error retrieving functions for profile:", err)
		return &Model{Items: nil} // Return an empty model on error
	}

	// Populate the model with FonctionItems
	items := make([]Item, len(Asso.SousFonctions))
	for i, fonction := range Asso.SousFonctions {
		name := fonction.Nom
		if fonction.Code != "" {
			name += " / " + fonction.Code
		}
		items[i] = Item{name: name, value: fonction.Nom}
	}

	return &Model{Items: items}
}

// Affichage du nombre d'éléments d'un modele
func (m *Model) ItemCount() int {
	return len(m.Items)
}

// Récupération de la valeur d'élément du modele
func (m *Model) Value(index int) interface{} {
	return m.Items[index].name
}
