package tree_carbon

// TreesForCarbonKG returns the number of trees required to offset
// the given number of kilograms CO2. One tree, living 45 years, can
// "sequester" approximately 1000kg (one tonne) of CO2 in its lifetime
// (as per http://www.unm.edu/~jbrink/365/Documents/Calculating_tree_carbon.pdf)
func TreesForCarbonKG(carbon int) int {
	return carbon / 1000
}

func CarbonKGForTrees(trees int) int {
	return trees * 1000
}
