package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	// Fonctions GitHub
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type MyMainWindow struct {
	*walk.MainWindow
	// Le titre de la première liste
	titledata *walk.Label
	// Le titre de la dexieme liste
	titledata2 *walk.Label
	// Le model pour la liste des profils
	model *Model
	// Le modele pour la liste des fonctions
	funcmodel *Model
	// La barre de recherche
	search *walk.TextEdit
	// La liste des profils
	lb *walk.ListBox
	//La liste des Fonctions
	lbfunc *walk.ListBox
	// Les resultats  des la prèmiere zone
	te *walk.TextLabel
	// Les resultats de la deuxième zone
	te2        *walk.TextLabel
	resultType string
}

// Création du menu
func initMainMenu() {
	mw := &MyMainWindow{model: NewProfilModel(), funcmodel: NewFonctionModel()}
	if _, err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "Recherche Droits",
		MinSize:  Size{Width: 720, Height: 520},
		Size:     Size{Width: 1080, Height: 800},
		Layout:   VBox{MarginsZero: true},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					VSplitter{
						Children: []Widget{
							Label{
								Text: "Recherche",
							},

							TextEdit{
								MaxSize:       Size{Height: 20},
								MinSize:       Size{Height: 20},
								AssignTo:      &mw.search,
								OnTextChanged: mw.Search,
							},
							Label{
								Text: "Listes ",
							},
							Label{
								Text:          "Profils ",
								TextAlignment: AlignCenter,
							},
							ListBox{
								AssignTo:              &mw.lb,
								Model:                 mw.model,
								OnCurrentIndexChanged: mw.lb_CurrentIndexChanged,
								OnItemActivated:       mw.lb_ProfilActivated,
							},
							Label{
								Text:          "Fonctions",
								TextAlignment: AlignCenter,
							},
							ListBox{
								AssignTo:              &mw.lbfunc,
								Model:                 mw.funcmodel,
								OnCurrentIndexChanged: mw.lbfunc_CurrentIndexChanged,
								OnItemActivated:       mw.lb_FonctionActivated,
							},
						},
					},
					VSplitter{

						Children: []Widget{
							Label{
								Text: "Resultats ",
							},
							PushButton{
								Text:      "Imprimer les resultats",
								OnClicked: mw.PrintResults,
							},
							PushButton{
								Text:      "Ajouter des Profils et Fonctions",
								OnClicked: mw.AddFunctionProfils,
							},
							PushButton{
								Text:      "Supprimer tout",
								OnClicked: mw.DeleteAll,
							},
							Label{
								AssignTo:      &mw.titledata,
								TextAlignment: AlignCenter,
							},
							ScrollView{
								Layout: VBox{MarginsZero: true},
								Children: []Widget{
									TextLabel{
										AssignTo: &mw.te,
										//ReadOnly: true,
									},
								},
							},
							Label{
								AssignTo:      &mw.titledata2,
								TextAlignment: AlignCenter,
							},
							ScrollView{
								Layout: VBox{MarginsZero: true},
								Children: []Widget{
									TextLabel{
										AssignTo: &mw.te2,
										//ReadOnly: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}
}

// fonction permettant d'enregistrer le résultats dan un fichier texte
func (mw *MyMainWindow) PrintResults() {

	dirName := "./resultats"

	// Check if the directory exists
	_, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		// Directory does not exist, create it
		if err := os.Mkdir(dirName, os.FileMode(0777)); err != nil {
			walk.MsgBox(mw, "Erreur", "Une erreur s'est produite lors de la création du répertoire", walk.MsgBoxIconInformation)
			fmt.Println(err)
		}
	} else if err != nil {
		// Error occurred while checking directory existence
		walk.MsgBox(mw, "Erreur", "Une erreur s'est produite lors de la vérification du répertoire", walk.MsgBoxIconInformation)
		fmt.Println(err)
	}

	var name string

	if mw.resultType == "Fonction" {
		if mw.lbfunc.CurrentIndex() != -1 {
			name = dirName + "/" + mw.funcmodel.Items[mw.lbfunc.CurrentIndex()].value + ".txt"
		}
	}

	if mw.resultType == "Profil" {
		if mw.lb.CurrentIndex() != -1 {
			name = dirName + "/" + mw.model.Items[mw.lb.CurrentIndex()].value + ".txt"
		}

	}
	if name != "" {
		f, err := os.Create(name)
		if err != nil {
			walk.MsgBox(mw, "Erreur", "Une erreur s'est produite lors de la création du fichier", walk.MsgBoxIconInformation)
			log.Fatal(err)
		}
		defer f.Close()

		_, err = f.WriteString(mw.titledata.Text())
		if err != nil {
			log.Fatal(err)
		}

		_, err = f.WriteString(mw.te.Text())
		if err != nil {
			log.Fatal(err)
		}

		_, err = f.WriteString(mw.titledata2.Text())
		if err != nil {
			log.Fatal(err)
		}

		_, err = f.WriteString(mw.te2.Text())
		if err != nil {
			log.Fatal(err)
		}

		walk.MsgBox(mw, "Terminé !!!", "Un fichier .txt sous le nom de "+mw.model.Items[mw.lb.CurrentIndex()].value+" a été crée avec succées dossier résultats", walk.MsgBoxIconInformation)
	} else {
		walk.MsgBox(mw, "Terminé !!!", "La sélection est vide", walk.MsgBoxIconInformation)

	}
}

// Création de la fenetre principal
func main() {
	initMainMenu()
}

// Fonction de recherhce d'un profil ou d'une fonction avec le nom/code
func (mw *MyMainWindow) Search() {
	q := mw.search.Text()
	// Search is performed only when the field is not empty
	if q == "" {
		return
	}

	// Define a map to store the indexes for faster lookup
	indexMap := make(map[string]int)
	for i, profile := range mw.model.Items {
		indexMap[profile.value] = i
	}

	// Search for profile by name
	res, err := getProfilByName(q)
	if err == nil {
		if index, ok := indexMap[res.Nom]; ok {
			mw.lb.SetCurrentIndex(index)
		}
		return
	}

	// Search for profile by code
	res, err = getProfilByCode(q)
	if err == nil {
		if index, ok := indexMap[res.Nom]; ok {
			mw.lb.SetCurrentIndex(index)
		}
		return
	}

	// Search for function by name
	res, err = getFonctionByName(q)
	if err == nil {
		for i, profile := range mw.funcmodel.Items {
			if profile.value == res.Nom {
				mw.lbfunc.SetCurrentIndex(i)
				break
			}
		}
		return
	}

	// Search for function by code
	res, err = getFonctionByCode(q)
	if err == nil {
		for i, profile := range mw.funcmodel.Items {
			if profile.value == res.Nom {
				mw.lbfunc.SetCurrentIndex(i)
				break
			}
		}
		return
	}
}

// Qunad l'élément seléctionné de la liste des profils change Mis à jour de la liste des  résultats pour y afficher les fonctions liée au profil sélectionnné
func (mw *MyMainWindow) lb_CurrentIndexChanged() {
	mw.resultType = "Profil"
	i := mw.lb.CurrentIndex()
	var s strings.Builder
	var r strings.Builder
	if i < 0 || i >= len(mw.model.Items) {
		r.WriteString("dépassement de l'index")
	} else {
		item := &mw.model.Items[i]
		itemprofil, err := getProfilByName(item.value)
		if err != nil {
			r.WriteString("Profil non trouvé")
			mw.te.SetTextColor(walk.RGB(255, 0, 0))
		} else {
			//mw.funcmodel = NewFonctionModel(item.value)
			Asso, err := getFunctionsByProfil(itemprofil)
			if err != nil {
				r.WriteString("Fonctions introuvables")
				mw.te.SetTextColor(walk.RGB(255, 0, 0))
				fmt.Println(err)
			} else {
				mw.titledata.SetText(fmt.Sprintf("--- Fonctions de %s : %d ---\r\n", Asso.NomProfil, len(Asso.Fonctions)))
				if len(Asso.Fonctions) > 0 {
					for i, fonction := range Asso.Fonctions {
						if fonction.Code == "" {
							r.WriteString(fmt.Sprintf(" %d | Nom : %s\n", i, fonction.Nom))
						} else {
							r.WriteString(fmt.Sprintf(" %d | Nom : %s / Code : %s \r\n", i, fonction.Nom, fonction.Code))
						}
					}
				} else {
					r.WriteString(" !!!!! Aucune fonction enregistrée !!!!!\r\n")
					mw.te.SetTextColor(walk.RGB(255, 0, 0))
				}

				if len(Asso.SousFonctions) > 0 {
					/*if mw.titledata2 == nil {
						mw.titledata2 = &walk.Label{}
						mw.te2 = &walk.TextEdit{}
					}*/
					mw.titledata2.SetText(fmt.Sprintf("--- Sous - Fonctions de %s : %d ---\r\n", Asso.NomProfil, len(Asso.SousFonctions)))
					for i, fonction := range Asso.SousFonctions {
						if fonction.Code == "" {
							s.WriteString(fmt.Sprintf(" %d | Nom : %s \n\n", i, fonction.Nom))
						} else {
							s.WriteString(fmt.Sprintf(" %d | Nom : %s / Code : %s \r\n", i, fonction.Nom, fonction.Code))
						}
					}

				} else {
					mw.titledata2.SetText("")
					s.WriteString("")
				}

			}

		}
	}
	mw.te2.SetText(s.String())
	mw.te.SetText(r.String())
}

// Quand l'élément sélection de la liste des fonctions change. Mis à jour de la liste des  résultats pour y afficher les profils liée au fonction sélectionnnée
func (mw *MyMainWindow) lbfunc_CurrentIndexChanged() {
	mw.resultType = "Fonction"
	i := mw.lbfunc.CurrentIndex()

	var r strings.Builder
	var s strings.Builder
	
	if i < 0 || i >= len(mw.funcmodel.Items) {
		r.WriteString("dépassement de l'index")
	} else {
		item := &mw.funcmodel.Items[i]
		//mw.funcmodel = NewFonctionModel(item.value)
		//fmt.Println(item.value)
		Asso, err := getProfilByFonction(item.value)
		if err != nil {
			r.WriteString("Profils introuvables")
		} else {
			mw.titledata.SetText(fmt.Sprintf("--- Utilisateurs principaux de la fonction %s : %d ---\r\n", item.value, len(Asso)))
			if len(Asso) > 0 {
				for i, profile := range Asso {
					if profile.CodeProfil == "" {
						r.WriteString(fmt.Sprintf(" %d | Nom : %s\r\n", i, profile.NomProfil))
					} else {
						r.WriteString(fmt.Sprintf(" %d | Nom : %s/ Code : %s \r\n", i, profile.NomProfil, profile.CodeProfil))
					}
				}
			} else {
				r.WriteString(" !!!!! Aucun profil utilise cette fonction !!!!!\r\n")
			}

			Asso, err = getProfilBySubFonction(item.value)
			if err != nil {

			} else {
				if len(Asso) > 0 {
					mw.titledata2.SetText(fmt.Sprintf("--- Utilisateurs Secondaires de la fonction %s : %d ---\r\n", item.value, len(Asso)))
					//r.WriteString()
					for i, profile := range Asso {
						if profile.CodeProfil == "" {
							s.WriteString(fmt.Sprintf(" %d | Nom : %s\r\n", i, profile.NomProfil))
						} else {
							s.WriteString(fmt.Sprintf(" %d | Nom : %s/ Code : %s \r\n", i, profile.NomProfil, profile.CodeProfil))
						}
					}
				} else {
					s.WriteString("")
					mw.titledata2.SetText("")
				}
			}

		}
	}
	mw.te.SetText(r.String())
	mw.te2.SetText(s.String())
}

func (mw *MyMainWindow) lb_ProfilActivated() {
	// Rendre inactive la fenetre principale
	mw.SetEnabled(false)

	value := mw.model.Items[mw.lb.CurrentIndex()].value
	item, err := getProfilByName(value)
	if err != nil {
		walk.MsgBox(mw, "Une erreur s'est produit", "Probleme de recherche du profil", walk.MsgBoxIconInformation)
		mw.SetEnabled(true)
		return
	}

	msgBox := walk.MsgBox(nil, " Confirmation", "Voulez - vous supprimer le profil suivant : \n"+value, walk.MsgBoxYesNo|walk.MsgBoxIconQuestion)
	if msgBox == walk.DlgCmdYes {
		err := deleteprofil(item)
		if err != nil {
			mw.SetEnabled(true)
			return
		}
		mw.Close()

		// Display a message box for successful transfer
		walk.MsgBox(mw, "Suppression terminée", "Suppression terminée", walk.MsgBoxIconInformation)

		// Initialize the main menu
		initMainMenu()
	} else {
		mw.SetEnabled(true)
	}

}

func (mw *MyMainWindow) lb_FonctionActivated() {
	mw.SetEnabled(false)

	value := mw.funcmodel.Items[mw.lbfunc.CurrentIndex()].value
	item, err := getFonctionByName(value)
	if err != nil {
		walk.MsgBox(mw, "Une erreur s'est produit", "Probleme de recherche de la fonction", walk.MsgBoxIconInformation)
		mw.SetEnabled(true)
		return
	}
	msgBox := walk.MsgBox(nil, " Confirmation", "Voulez - vous supprimer la fonction suivant : \n"+value, walk.MsgBoxYesNo|walk.MsgBoxIconQuestion)
	if msgBox == walk.DlgCmdYes {
		err := deletefunction(item)
		if err != nil {
			mw.SetEnabled(true)
			return
		}
		mw.Close()

		// Display a message box for successful transfer
		walk.MsgBox(mw, "Suppression terminée", "Suppression terminée", walk.MsgBoxIconInformation)

		// Initialize the main menu
		initMainMenu()
	} else {
		mw.SetEnabled(true)
	}
}

// Permet d'ajouter des profils et des fonctions à la liste via un fichier XML
func (mw *MyMainWindow) AddFunctionProfils() {
	mw.SetEnabled(false)
	q, err := readXml()
	if err != nil {
		walk.MsgBox(mw, "Echec du Transfer", "Echec de la conversion !!!", walk.MsgBoxIconInformation)
		mw.SetEnabled(true)
		return
	}

	if q.Associations != nil {
		insertXml(deleteduplicate(q))
		mw.Close()

		// Display a message box for successful transfer
		walk.MsgBox(mw, "Fin du Transfer", "Transfert terminée !!!!!", walk.MsgBoxIconInformation)

		// Initialize the main menu
		initMainMenu()
	} else {
		// Display a message box for invalid XML format
		walk.MsgBox(mw, "Echec du Transfer", "Format XML invalide !!!", walk.MsgBoxIconInformation)
		mw.SetEnabled(true)
	}
}

func (mw *MyMainWindow) DeleteAll() {
	msgBox := walk.MsgBox(nil, " Confirmation", "Voulez - vous tout supprimer ?", walk.MsgBoxYesNo|walk.MsgBoxIconQuestion)
	if msgBox == walk.DlgCmdYes {
		err := deleteAllasociations()
		if err != nil {
			mw.SetEnabled(true)
			return
		}
		mw.Close()

		// Display a message box for successful transfer
		walk.MsgBox(mw, "Suppression terminée", "Suppression terminée", walk.MsgBoxIconInformation)

		// Initialize the main menu
		initMainMenu()
	} else {
		mw.SetEnabled(true)
	}
}
