package main

import (
	"log"

	// Fonctions GitHub
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type UpdatePage struct {
	*walk.MainWindow
	// Champ du nom
	NameEdit *walk.TextLabel
	// Champ du code
	CodeEdit *walk.TextLabel
	//  Liste des fonctions
	FunctionList *walk.ListBox
	// Liste des fonctions utilisée par le profil conserné
	FunctionListEditBox *walk.ListBox
	// Liste des Sous - fonctions utilisée par le profil conserné
	SubFunctionListEditBox *walk.ListBox
	// Modèle pour la liste des fonctions
	FunctionListModel *Model
	// Modèle pour la liste des fonctions du profil
	FunctionListEditModel *Model
	// Modèle pour la liste des sous-fonctions du profil
	SubFunctionListEditModel *Model
}

// Initialisationn de la page de mise à jour
func UpadtePageLaunch(p1 string) {
	profil, err := getProfilByName(p1)
	if err != nil {

	} else {
		up := &UpdatePage{SubFunctionListEditModel: NewSubFonctionUpdateModel(profil.Nom), FunctionListEditModel: NewFonctionUpdateModel(profil.Nom), FunctionListModel: NewFonctionModel()}
		if _, err := (MainWindow{
			AssignTo: &up.MainWindow,
			Title:    "Modification Profil : " + p1,
			Size:     Size{Width: 720, Height: 520},
			Layout:   VBox{MarginsZero: true},
			Children: []Widget{
				Label{
					Text:          "Mise à jour du profil " + p1,
					TextAlignment: AlignCenter,
				},
				HSplitter{
					Children: []Widget{

						VSplitter{
							Children: []Widget{
								TextLabel{
									AssignTo: &up.NameEdit,
									Text:     "Nom : " + profil.Nom,
								},
								TextLabel{
									AssignTo: &up.CodeEdit,
									Text:     "Code : " + profil.Code,
								},
								Label{
									Text: "Fonctions ",
								},
								ListBox{
									AssignTo: &up.FunctionListEditBox,
									Model:    up.FunctionListEditModel,
									//OnCurrentIndexChanged: mw.lb_CurrentIndexChanged,
									//OnItemActivated:       mw.lb_ItemActivated,
								},
								Label{
									Text: "Sous - Fonctions",
								},
								ListBox{
									AssignTo: &up.SubFunctionListEditBox,
									Model:    up.SubFunctionListEditModel,
									//OnCurrentIndexChanged: mw.lbfunc_CurrentIndexChanged,
									//						OnItemActivated:       mw.lb_ItemActivated,
								},
							},
						},
						VSplitter{

							Children: []Widget{
								Label{
									Text: "Liste des fonctions",
								},
								ListBox{
									AssignTo: &up.FunctionList,
									Model:    up.FunctionListModel,
									//OnCurrentIndexChanged: mw.lbfunc_CurrentIndexChanged,
									//OnItemActivated: up.fucncAddActivated,
								},
								PushButton{
									Text: "+",
								},
							},
						},
					},
				},
				PushButton{
					Text: "Enregistrer",
				},
			},
		}.Run()); err != nil {
			log.Fatal(err)
		}
	}

	//up.NameEdit.SetText(p1)
}

// Ajoute une fonction dans la liste des fonctions du profil
func (up *UpdatePage) AddFunctionInList(f1 string) {
	fonction, err := getFonctionByName(f1)
	if err != nil {

	} else {

		name := fonction.Nom
		if fonction.Code != "" {
			name += " / " + fonction.Code
		}
		up.FunctionListEditModel.Items[len(up.FunctionListEditModel.Items)] = Item{name: name, value: fonction.Nom}
	}
}

// Ajout de la fonction dans la liste des sous-fonction du profil
func (up *UpdatePage) AddSubFunctionInList(f1 string) *Model {
	function, err := getFonctionByName(f1)
	if err != nil {
		return up.SubFunctionListEditModel
	} else {

		items := make([]Item, len(up.SubFunctionListEditModel.Items))
		for i, fonction := range up.SubFunctionListEditModel.Items {

			items[i] = Item{name: fonction.name, value: fonction.value}
		}

		name := function.Nom
		if function.Code != "" {
			name += " / " + function.Code
		}
		items[len(up.SubFunctionListEditModel.Items)] = Item{name: name, value: function.Nom}
		return &Model{Items: items}
	}
}

/*
func (up *UpdatePage) fucncAddActivated() {
	value := up.FunctionListModel.Items[up.FunctionList.CurrentIndex()].value

	//walk.MsgBox(up, "Value", value, walk.MsgBoxIconInformation)
	//UpadtePageLaunch(value)

	up.AddFunctionInList(value)
}
*/
