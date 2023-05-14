package main

import (
	"fmt"
	"reflect"
)

func StructComparison(old, new Recipes) {
	if reflect.DeepEqual(old, new) {
		fmt.Println("equal")
		return
	}
	PrintChangesRecipes(old, new)
}

func PrintChangesRecipes(old, new Recipes) {
	idxcake := Checkcake(old, new)
	Checktime(old, new, idxcake)
	Checkingredient(old, new, idxcake)
	CheckChangeIngredients(old, new, idxcake)

}

func Checkcake(old, new Recipes) map[int]int {
	idxcakes := make(map[int]int)
	oldcake := make(map[string]int)
	for i, cake := range old.Cake {
		oldcake[cake.Name] = i
	}
	newcake := make(map[string]int)
	for i, cake := range new.Cake {
		newcake[cake.Name] = i
	}
	for cakename, idxnewcake := range newcake {
		idxoldcake, ok := oldcake[cakename]
		if !ok {
			fmt.Printf("ADDED cake \\\"%s\\\"\n", cakename)
		} else {
			idxcakes[idxoldcake] = idxnewcake
		}
	}
	for cakename := range oldcake {
		if _, ok := newcake[cakename]; !ok {
			fmt.Printf("REMOVED cake \\\"%s\\\"\n", cakename)
		}
	}
	return idxcakes
}

func Checktime(old, new Recipes, idx map[int]int) {
	for idxoldcake, idxnewcake := range idx {
		if old.Cake[idxoldcake].Time != new.Cake[idxnewcake].Time {
			fmt.Printf(
				"CHANGED cooking time for cake \\\"%s\\\" - \\\"%s\\\" instead of \\\"%s\\\" \n",
				old.Cake[idxoldcake].Name, new.Cake[idxnewcake].Time, old.Cake[idxoldcake].Time,
			)
		}
	}
}

func Checkingredient(old, new Recipes, idx map[int]int) {
	for idxoldcake, idxnewcake := range idx {
		oldingredient := make(map[string]int)
		for i, ingr := range old.Cake[idxoldcake].Ingredients {
			oldingredient[ingr.IngredientName] = i
		}
		newingredient := make(map[string]int)
		for i, ingr := range new.Cake[idxnewcake].Ingredients {
			newingredient[ingr.IngredientName] = i
		}

		for ingrname := range newingredient {
			_, ok := oldingredient[ingrname]
			if !ok {
				fmt.Printf("ADDED ingredient \\\"%s\\\" for cake \\\"%s\\\" \n", ingrname, new.Cake[idxnewcake].Name)
			}
		}
		for ingrname := range oldingredient {
			if _, ok := newingredient[ingrname]; !ok {
				fmt.Printf("REMOVED ingredient \\\"%s\\\" for cake \\\"%s\\\" \n", ingrname, old.Cake[idxoldcake].Name)
			}
		}

	}
}

func CheckChangeIngredients(old, new Recipes, idx map[int]int) {
	for idxoldcake, idxnewcake := range idx {
		for _, oldingr := range old.Cake[idxoldcake].Ingredients {
			for _, newingr := range new.Cake[idxnewcake].Ingredients {
				if oldingr.IngredientName == newingr.IngredientName {
					if oldingr.IngredientCount != newingr.IngredientCount {
						fmt.Printf(
							"CHANGED unit for ingredient \\\"%s\\\" for cake \\\"%s\\\" - \\\"%s\\\" instead of \\\"%s\\\" \n",
							newingr.IngredientName, new.Cake[idxnewcake].Name, newingr.IngredientCount,
							oldingr.IngredientCount,
						)
					}
					if oldingr.IngredientUnit != newingr.IngredientUnit {
						if oldingr.IngredientUnit == "" {
							fmt.Printf(
								"ADDED unit \\\"%s\\\" for ingredient \\\"%s\\\" for cake \\\"%s\\\"\n",
								newingr.IngredientUnit, oldingr.IngredientName, old.Cake[idxoldcake].Name,
							)
						} else if newingr.IngredientUnit == "" {
							fmt.Printf(
								"REMOVED unit \\\"%s\\\" for ingredient \\\"%s\\\" for cake \\\"%s\\\"\n",
								oldingr.IngredientUnit, oldingr.IngredientName, old.Cake[idxoldcake].Name,
							)
						} else {
							fmt.Printf(
								"CHANGED unit for ingredient \\\"%s\\\" for cake  \\\"%s\\\" - \\\"%s\\\" instead of \\\"%s\\\" \n",
								oldingr.IngredientName, old.Cake[idxoldcake].Name, newingr.IngredientUnit,
								oldingr.IngredientUnit,
							)
						}
					}
				}
			}
		}
	}
}
