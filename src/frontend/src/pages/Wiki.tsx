import { useState, useEffect } from "react";

interface Recipe {
  name: string;
  image_path: string;
  ingredients?: {
    element_1: string;
    element_2: string;
  };
}

export default function Wiki() {
  const [recipes, setRecipes] = useState<Recipe[]>([]);
  const [search, setSearch] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(`http://${import.meta.env.VITE_PUBLIC_BACKEND_BASE_URL || "localhost:4000"}/recipes`)
      .then(res => res.json())
      .then(data => {
        setRecipes(data);
        setLoading(false);
      })
      .catch(err => {
        console.error("Error fetching recipes:", err);
        setLoading(false);
      });
  }, []);

  const filteredRecipes = recipes.filter(recipe => 
    recipe.name.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
      <div className="bg-white shadow rounded-lg overflow-hidden">
        <div className="p-6">
          <h1 className="text-2xl font-bold text-gray-900 mb-6">Recipe Wiki</h1>
          
          <div className="mb-6">
            <input
              type="text"
              placeholder="Search recipes..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-full p-3 border border-gray-300 rounded-lg focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
          
          {loading ? (
            <div className="text-center py-10">
              <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500 mx-auto mb-3"></div>
              <p className="text-gray-500">Loading recipes...</p>
            </div>
          ) : filteredRecipes.length > 0 ? (
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
              {filteredRecipes.map((recipe, index) => (
                <div key={index} className="border rounded-lg overflow-hidden shadow-sm hover:shadow-md transition-shadow">
                  <div className="p-4">
                    <div className="flex items-center mb-3">
                      <img src={recipe.image_path} alt={recipe.name} className="w-10 h-10 mr-3" />
                      <h3 className="font-medium text-gray-900">{recipe.name}</h3>
                    </div>
                    
                    {recipe.ingredients && (
                      <div className="mt-2 text-sm text-gray-500">
                        <div className="font-semibold mb-1">Recipe:</div>
                        <div className="flex items-center space-x-2">
                          <div>{recipe.ingredients.element_1}</div>
                          <div>+</div>
                          <div>{recipe.ingredients.element_2}</div>
                        </div>
                      </div>
                    )}
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-10">
              <svg className="w-16 h-16 mx-auto text-gray-400 mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <p className="text-gray-500">No recipes found matching your search.</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}