package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"

	// Fonctions GitHub
	"github.com/lxn/walk"
	"github.com/sqweek/dialog"

	//Fonctions MongoDB
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Type correspondant à une liste d'assoiciations
type Query struct {
	Associations []AssociationProfilFonctionsEtSousFonctions `xml:"AssociationsProfilsFonctionsEtSousFonctions>AssociationProfilFonctionsEtSousFonctions"`
}

// Type correspondant à une association contenant un profil , des fonctions et des sous-fonctions
type AssociationProfilFonctionsEtSousFonctions struct {
	XMLName       xml.Name  `xml:"AssociationProfilFonctionsEtSousFonctions"`
	Profil        Element   `xml:"Profil"`
	Fonctions     []Element `xml:"Fonctions>Fonction" json:"fonctions"`
	SousFonctions []Element `xml:"SousFonctions>SousFonction" json:"sousfonctions"`
}

// Type correspondant à un element ( Profil , Fonctions , Sous - Fonctions)
type Element struct {
	Code string `xml:"Code,attr" json:"code"`
	Nom  string `xml:"Nom,attr" json:"nom"`
}

// Permet de lire l fichier XML inséré
func readXml() (Query, error) {
	var q Query

	filePath, err := dialog.File().Title("Select XML File").Load()
	if err != nil {
		fmt.Println("Error:", err)
		return q, err
	}

	xmlFile, err := os.Open(filePath)
	//Si une errreur a été envoyé
	if err != nil {
		fmt.Println("Error opening file:", err)
		return q, err
	}
	defer xmlFile.Close()

	b, err := io.ReadAll(xmlFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return q, err
	}

	err = xml.Unmarshal(b, &q)
	if err != nil {
		fmt.Println("Error unmarshalling XML:", err)
		return q, err
	}

	return q, err
}

// Insertion du fichier XML dans la base de données
func insertXml(q Query) {
	//Ajout des informations issues du fichier
	initClient()
	for _, association := range q.Associations {
		user := bson.D{{Key: "nom", Value: association.Profil.Nom}, {Key: "code", Value: association.Profil.Code}}

		var p Element

		//On verifie si , dans la base de données , le profil y est inscrit
		err := usersCollection.FindOne(context.TODO(), user).Decode(&p)
		//Si une errreur est détecté
		if err != nil {
			//Si cette erreur concerne le fait que la fonction dit que lke document rechercher n'existe pas
			if err == mongo.ErrNoDocuments {
				//Alors on réalisera une insertion
				usersCollection.InsertOne(context.TODO(), user)

				fmt.Printf("Ajout du profil %s réussi !!\n", association.Profil.Nom)
			}
		}

		//Verification de la présence de fonctions
		if len(association.Fonctions) != 0 {

			//Boucle d'ajout des fonctions dans la base de données
			for _, fonction := range association.Fonctions {

				user := bson.D{{Key: "nom", Value: fonction.Nom}, {Key: "code", Value: fonction.Code}}

				var fl Element
				// Supression des doublons
				err := funcsCollection.FindOne(context.TODO(), user).Decode(&fl)
				if err != nil {
					// Aucun doublon trouvé
					if err == mongo.ErrNoDocuments {
						funcsCollection.InsertOne(context.TODO(), user)
						fmt.Printf("Ajout de la fonction %s réussi !!\n", fonction.Nom)
					}
				}

			}
		}

		//Verification de la présence de sous-fonction
		if len(association.SousFonctions) != 0 {

			for _, fonction := range association.SousFonctions {
				user := bson.D{{Key: "nom", Value: fonction.Nom}, {Key: "code", Value: fonction.Code}}

				var fl Element
				// Supression des doublons
				err := funcsCollection.FindOne(context.TODO(), user).Decode(&fl)
				if err != nil {
					// Aucun doublon trouvé
					if err == mongo.ErrNoDocuments {
						funcsCollection.InsertOne(context.TODO(), user)
						fmt.Printf("Ajout de la fonction %s réussi !!\n", fonction.Nom)
					}
				}

			}

		}

		// Création de l'association
		//Suppression des doublons
		var a1 Association
		var a2 Association

		// Verifaction de la précense d'un document de la collection association avec le meme nom de profil
		as := bson.D{{Key: "nomprofil", Value: association.Profil.Nom}, {Key: "codeprofil", Value: association.Profil.Code}}
		err = AssoCollection.FindOne(context.TODO(), as).Decode(&a1)
		// verification des valeur du document de la Base de données avec le documents actuelle
		asso := bson.D{{Key: "nomprofil", Value: association.Profil.Nom}, {Key: "codeprofil", Value: association.Profil.Code}, {Key: "fonctions", Value: association.Fonctions}, {Key: "sousfonctions", Value: association.SousFonctions}}

		// Si le docuement trouvé a le meme nom que le document actuelle
		if err == nil {
			// Verifiée que les deux document ont des données identique
			err = AssoCollection.FindOne(context.TODO(), asso).Decode(&a2)
			if err != nil && err == mongo.ErrNoDocuments {
				//Mis à jour du document dans la base de données
				UpdateDoc(a1, association)
			}
			//Sinon
		} else {
			//Si le document n'est pas trouvée
			if err == mongo.ErrNoDocuments {
				//Créationn d'un nouveau document
				AssoCollection.InsertOne(context.TODO(), asso)
			}
		}

	}
}

// Supresion d'une fonction dans une association
func remove(s []Element, i int) []Element {
	if i < 0 || i >= len(s) {
		// index out of range, return the original slice
		return s
	}
	s[i] = s[len(s)-1]
	// Ensure that the slice is not empty before truncating it
	if len(s) > 0 {
		return s[:len(s)-1]
	}
	return s
}

// Permet de suprimmer les fonctions en double ( WIP)
func deleteduplicate(q Query) Query {
	for _, association := range q.Associations {

		if len(association.Fonctions) != 0 {
			// Remove duplicates from Fonctions
			for j, fonction := range association.Fonctions {
				for i, search := range association.Fonctions {
					if fonction.Nom == search.Nom && i != j {
						association.Fonctions = remove(association.Fonctions, i)
					}
				}
			}
		}

		if len(association.SousFonctions) != 0 {
			// Remove duplicates from SousFonctions
			for j, sousFonction := range association.SousFonctions {
				for i, search := range association.SousFonctions {
					if sousFonction.Nom == search.Nom && i != j {
						association.SousFonctions = remove(association.SousFonctions, i)
					}
				}
			}
		}
	}
	return q
}

// Procedures de mis à jour des données de l'associations
func UpdateDoc(a1 Association, a2 AssociationProfilFonctionsEtSousFonctions) {

	// Create Base.xml
	if err := writeXMLToFile("Base.xml", AssociationConvertJsontoXML(a1)); err != nil {
		fmt.Println("Error creating/updating Base.xml: ", err)
		return
	}

	// Create New.xml
	if err := writeXMLToFile("New.xml", a2); err != nil {
		fmt.Println("Error creating/updating New.xml: ", err)
		return
	}

	// Prompt user for confirmation
	msgBox := walk.MsgBox(nil, "Update Confirmation", "Des modifications ont été réalisées sur l'association du profil suivant : \n"+a2.Profil.Nom+".\nVoulez-vous mettre à jour les données existantes avec les nouvelles ?", walk.MsgBoxYesNo|walk.MsgBoxIconQuestion)
	if msgBox == walk.DlgCmdYes {
		updateDatabase(a2)
	}
}

// Génère un fichier XML
func writeXMLToFile(filename string, data interface{}) error {
	xmlFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer xmlFile.Close()

	xmlFile.WriteString(xml.Header)
	encoder := xml.NewEncoder(xmlFile)
	encoder.Indent("", "\t")
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

// Mis à jour de la base de données
func updateDatabase(a2 AssociationProfilFonctionsEtSousFonctions) {
	filter := bson.D{{Key: "nomprofil", Value: a2.Profil.Nom}}
	replacement := bson.D{
		{Key: "nomprofil", Value: a2.Profil.Nom},
		{Key: "codeprofil", Value: a2.Profil.Code},
		{Key: "fonctions", Value: a2.Fonctions},
		{Key: "sousfonctions", Value: a2.SousFonctions}, // Use underscore instead of space in field name
	}

	if _, err := AssoCollection.ReplaceOne(context.TODO(), filter, replacement); err != nil {
		log.Fatal(err)
	}
}

func AssociationConvertJsontoXML(a1 Association) AssociationProfilFonctionsEtSousFonctions {
	c1 := AssociationProfilFonctionsEtSousFonctions{
		Profil: Element{
			Code: a1.CodeProfil,
			Nom:  a1.NomProfil,
		},
		Fonctions:     make([]Element, len(a1.Fonctions)),
		SousFonctions: make([]Element, len(a1.SousFonctions)),
	}

	for i, fonction := range a1.Fonctions {
		c1.Fonctions[i].Nom = fonction.Nom
		c1.Fonctions[i].Code = fonction.Code
	}

	for i, fonction := range a1.SousFonctions {
		c1.SousFonctions[i].Nom = fonction.Nom
		c1.SousFonctions[i].Code = fonction.Code
	}

	return c1
}
