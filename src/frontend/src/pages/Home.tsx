import { Link } from "react-router";

export default function Home() {
  return (
    <div className="min-h-screen bg-gray-50">

      <div className="max-w-7xl mx-auto py-12 px-4 sm:px-6 lg:py-16 lg:px-8">
        <div className="bg-white shadow-lg rounded-lg overflow-hidden">
          <div className="px-6 py-8 sm:p-10">
            <h1 className="text-3xl font-extrabold text-gray-900 mb-6">Welcome to CCP Little Alchemy 2</h1>

            <p className="text-lg text-gray-700 mb-6">Discover the optimal crafting paths for all your favorite items!</p>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mt-10">
              <div className="bg-blue-50 p-6 rounded-lg">
                <h2 className="text-xl font-bold text-blue-800 mb-3">Recipe Visualizer</h2>
                <p className="text-gray-600 mb-4">
                  Explore different algorithms to find crafting recipes and visualize the crafting tree.
                </p>
                <Link to="/Visualizer" className="inline-block bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">
                  Try Visualizer
                </Link>
              </div>

              <div className="bg-green-50 p-6 rounded-lg">
                <h2 className="text-xl font-bold text-green-800 mb-3">Recipe Wiki</h2>
                <p className="text-gray-600 mb-4">
                  Browse all available recipes and learn how to craft any item with our comprehensive wiki.
                </p>
                <Link to="/Wiki" className="inline-block bg-green-600 text-white px-4 py-2 rounded hover:bg-green-700">
                  Browse Wiki
                </Link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
