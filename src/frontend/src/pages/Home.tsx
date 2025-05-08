import { useState } from "react";
import axios from "axios";

const Home = () => {
  interface ApiResponse {
    tree: Record<string, string | number | boolean | null | ApiResponse>;
    searchTime: number;
    visitedNodes: number;
  }

  const [target, setTarget] = useState<string>("");
  const [method, setMethod] = useState<string>("dfs");
  const [multipleRecipes, setMultipleRecipes] = useState<boolean>(false);
  const [maxRecipes, setMaxRecipes] = useState<number>(1);
  const [data, setData] = useState<ApiResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  const fetchData = () => {
    setLoading(true);
    setError(null);

    const url = `http://localhost:4000/api/tree?target=${target}&method=${method}&multiple=${multipleRecipes}&maxRecipes=${maxRecipes}`;
    axios
      .get(url)
      .then((response) => {
        setData(response.data);
      })
      .catch((err) => {
        setError(err.message);
      })
      .finally(() => {
        setLoading(false);
      });
  };

  return (
    <div className="min-h-screen bg-gray-50 p-6 justify-center flex items-center">
      <div className="max-w-2xl mx-auto bg-white rounded-2xl shadow-md p-6 space-y-6">
        <h1 className="text-2xl font-bold text-center text-indigo-600">Little Alchemy 2 Recipe Finder</h1>

        <div className="space-y-4">
          {/* Target Element */}
          <div>
            <label className="block text-sm font-medium text-gray-700">Target Element</label>
            <input
              type="text"
              value={target}
              onChange={(e) => setTarget(e.target.value)}
              placeholder="Enter target element"
              className="mt-1 w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
            />
          </div>

          {/* Search Method */}
          <div>
            <label className="block text-sm font-medium text-gray-700">Search Method</label>
            <select
              value={method}
              onChange={(e) => setMethod(e.target.value)}
              className="mt-1 w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
            >
              <option value="dfs">DFS</option>
              <option value="bfs">BFS</option>
            </select>
          </div>

          {/* Multiple Recipes Toggle */}
          <div className="flex items-center space-x-2">
            <input
              type="checkbox"
              checked={multipleRecipes}
              onChange={(e) => setMultipleRecipes(e.target.checked)}
              className="h-4 w-4 text-indigo-600 border-gray-300 rounded"
            />
            <label className="text-sm text-gray-700">Find multiple recipes</label>
          </div>

          {/* Max Recipes Input */}
          {multipleRecipes && (
            <div>
              <label className="block text-sm font-medium text-gray-700">Max Recipes</label>
              <input
                type="number"
                value={maxRecipes}
                onChange={(e) => setMaxRecipes(Number(e.target.value))}
                min={1}
                className="mt-1 w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
              />
            </div>
          )}

          {/* Search Button */}
          <div>
            <button
              onClick={fetchData}
              disabled={loading || !target}
              className={`w-full py-2 px-4 rounded-md font-semibold text-white ${
                loading || !target ? "bg-gray-400 cursor-not-allowed" : "bg-indigo-600 hover:bg-indigo-700"
              }`}
            >
              {loading ? "Searching..." : "Find Recipe"}
            </button>
          </div>
        </div>

        {/* Error Display */}
        {error && <p className="text-red-500 text-sm mt-4">Error: {error}</p>}

        {/* Results */}
        {data && (
          <div className="bg-gray-100 p-4 rounded-md mt-4">
            <h2 className="text-lg font-semibold text-indigo-700">Search Results</h2>
            <p className="text-sm text-gray-700">
              Search Time: <span className="font-medium">{data.searchTime} ms</span>
            </p>
            <p className="text-sm text-gray-700">
              Visited Nodes: <span className="font-medium">{data.visitedNodes}</span>
            </p>
            <h3 className="mt-2 font-semibold text-gray-800">Recipe Tree:</h3>
            <pre className="bg-white p-3 rounded-md text-sm overflow-auto max-h-96">{JSON.stringify(data.tree, null, 2)}</pre>
          </div>
        )}
      </div>
    </div>
  );
};

export default Home;
