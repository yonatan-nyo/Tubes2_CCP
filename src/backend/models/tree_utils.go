package models

func (root *RecipeTreeNode) clone() *RecipeTreeNode {
	clone := &RecipeTreeNode{
		Name:                   root.Name,
		ImagePath:              root.ImagePath,
		MinimumNodesRecipeTree: root.MinimumNodesRecipeTree,
	}

	if clone.Element1 == nil && clone.Element2 == nil || IsBaseElement(root.Name) {
		return clone
	}

	clone.Element1 = root.Element1.clone()
	clone.Element2 = root.Element2.clone()

	return clone
}
