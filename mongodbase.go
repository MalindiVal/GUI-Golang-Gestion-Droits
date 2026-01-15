package main

import (
	"context"
	"strings"

	//Fonctions MongoDB
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Collection de la base de données MongoDB
var (
	usersCollection *mongo.Collection
	funcsCollection *mongo.Collection
	AssoCollection  *mongo.Collection
)

// Type correspondant à un élément de la liste tirée d'une requete ver mongoDB ( Association)
type Association struct {
	// Le nom du profil de l'associaition
	NomProfil string `json:"nomprofil"`
	//Le code du profil de l'association
	CodeProfil string `json:"codeprofil"`
	// l'ensemble des fonction de l'association
	Fonctions []Element `json:"fonctions"`
	//L'ensemble des sous-fonction de l'association
	SousFonctions []Element `json:"sousfonctions"`
	AssociationProfilFonctionsEtSousFonctions
}

// Initialisation du la connexion avec la base de données
func initClient() {

	// Connection à la base de données
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	// En cas de problème
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	usersCollection = client.Database("malindi").Collection("Profils")
	funcsCollection = client.Database("malindi").Collection("Fonctions")
	AssoCollection = client.Database("malindi").Collection("Associations")
}

// fonction permmettant de récupérér l'ensemble des données de la collection des fonctions
func getAllFonctions() []Element {
	filter := bson.D{}
	var r string
	cursor, err := funcsCollection.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.TODO()) // Close the cursor when done
	var results []Element
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		res, err := bson.MarshalExtJSON(result, false, false)
		if err != nil {
			panic(err)
		}
		r += string(res) // Concatenate the JSON strings
	}
	return results
}

// Récupération de l'ensemble des profils
func getAllprofils() []Element {
	filter := bson.D{}
	var r string
	cursor, err := usersCollection.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.TODO()) // Close the cursor when done
	var results []Element
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		res, err := bson.MarshalExtJSON(result, false, false)
		if err != nil {
			panic(err)
		}
		r += string(res) // Concatenate the JSON strings
	}
	return results
}

// Récupèration d'un profil en fonction du nom
func getProfilByName(name string) (Element, error) {
	filter := bson.D{{Key: "nom", Value: name}}
	cursor := usersCollection.FindOne(context.TODO(), filter)
	var result Element
	err := cursor.Decode(&result)
	return result, err
}

// Récupération du profils en foncions de son code
func getProfilByCode(code string) (Element, error) {
	filter := bson.D{{Key: "code", Value: strings.ToUpper(code)}}
	cursor := usersCollection.FindOne(context.TODO(), filter)
	var result Element
	err := cursor.Decode(&result)
	return result, err
}

// Récuperation d'une fonction en fonction de son code
func getFonctionByCode(code string) (Element, error) {
	filter := bson.D{{Key: "code", Value: strings.ToUpper(code)}}
	cursor := funcsCollection.FindOne(context.TODO(), filter)
	var result Element
	err := cursor.Decode(&result)
	return result, err
}

// Récupération d'unre fonction avec son nom
func getFonctionByName(name string) (Element, error) {
	filter := bson.D{{Key: "nom", Value: name}}
	cursor := funcsCollection.FindOne(context.TODO(), filter)
	var result Element
	err := cursor.Decode(&result)
	return result, err
}

// récupération du profil en fonction des fonctions qu'il est autorisé à utiliser
func getProfilByFonction(name string) ([]Association, error) {
	var results []Association
	filter := bson.D{{Key: "fonctions.nom", Value: name}}
	var r string
	cursor, err := AssoCollection.Find(context.TODO(), filter)
	if err != nil {
		return results, err
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &results); err != nil {
		return results, err
	}
	for _, result := range results {
		res, err := bson.MarshalExtJSON(result, false, false)
		if err != nil {
			return results, err
		}
		r += string(res)
		//fmt.Println(result)
	}
	return results, nil
}

// Récupération des profil utilisant la sous fonction
func getProfilBySubFonction(name string) ([]Association, error) {
	var results []Association
	filter := bson.D{{Key: "sousfonctions.nom", Value: name}}
	var r string
	cursor, err := AssoCollection.Find(context.TODO(), filter)
	if err != nil {
		return results, err
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &results); err != nil {
		return results, err
	}
	for _, result := range results {
		res, err := bson.MarshalExtJSON(result, false, false)
		if err != nil {
			return results, err
		}
		r += string(res)
		//Lfmt.Println(result)
	}
	return results, nil
}

// Récuperation des fonctions utilisée par un profil choisi
func getFunctionsByProfil(p1 Element) (Association, error) {
	var result Association
	filter := bson.D{{Key: "nomprofil", Value: p1.Nom}, {Key: "codeprofil", Value: p1.Code}}
	err := AssoCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return result, err
	}
	//fmt.Println(result)
	return result, nil // Return the results if no error occurred
}

func deleteprofil(p1 Element) error {
	filter := bson.D{{Key: "nom", Value: p1.Nom}, {Key: "code", Value: p1.Code}}
	_, err := usersCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	filter = bson.D{{Key: "nomprofil", Value: p1.Nom}, {Key: "codeprofil", Value: p1.Code}}
	_, err = AssoCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}

func deleteAllasociations() error {
	filter := bson.D{}
	_, err := usersCollection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return err
	}
	_, err = AssoCollection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return err
	}
	_, err = funcsCollection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}

func deletefunction(p1 Element) error {
	filter := bson.D{{Key: "nom", Value: p1.Nom}, {Key: "code", Value: p1.Code}}
	_, err := funcsCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	filter = bson.D{{Key: "fonctions.nom", Value: p1.Nom}, {Key: "fonctions.code", Value: p1.Code}}
	cursor, err := AssoCollection.Find(context.TODO(), filter)
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var association Association
		if err := cursor.Decode(&association); err != nil {
			return err
		}

		// Remove the function from the function list of the association
		updateFilter := bson.D{{Key: "fonctions.nom", Value: p1.Nom}}
		update := bson.D{
			{Key: "$pull", Value: bson.D{
				{Key: "fonctions", Value: bson.D{
					{Key: "nom", Value: p1.Nom},
					{Key: "code", Value: p1.Code},
				}},
			}},
		}
		_, err := AssoCollection.UpdateOne(context.TODO(), updateFilter, update)
		if err != nil {
			return err
		}
	}
	if err := cursor.Err(); err != nil {
		return err
	}

	filter = bson.D{{Key: "sousfonctions.nom", Value: p1.Nom}, {Key: "sousfonctions.code", Value: p1.Code}}
	cursor, err = AssoCollection.Find(context.TODO(), filter)
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var association Association
		if err := cursor.Decode(&association); err != nil {
			return err
		}

		// Remove the function from the function list of the association
		updateFilter := bson.D{{Key: "sousfonctions.nom", Value: p1.Nom}}
		update := bson.D{
			{Key: "$pull", Value: bson.D{
				{Key: "sousfonctions", Value: bson.D{
					{Key: "nom", Value: p1.Nom},
					{Key: "code", Value: p1.Code},
				}},
			}},
		}
		_, err := AssoCollection.UpdateOne(context.TODO(), updateFilter, update)
		if err != nil {
			return err
		}
	}
	if err := cursor.Err(); err != nil {
		return err
	}

	return nil
}
