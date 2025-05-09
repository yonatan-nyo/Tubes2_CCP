package models

func (node *RecipeTreeNode) clone() *RecipeTreeNode {
	if node == nil {
		return nil
	}

	clone := &RecipeTreeNode{
		Name:      node.Name,
		ImagePath: node.ImagePath,
	}

	// Only clone children if they exist
	if node.Element1 != nil {
		clone.Element1 = node.Element1.clone()
	}
	if node.Element2 != nil {
		clone.Element2 = node.Element2.clone()
	}

	return clone
}
