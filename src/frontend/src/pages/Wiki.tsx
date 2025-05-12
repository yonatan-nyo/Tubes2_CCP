import { useState, useEffect } from "react";
import Axios from "axios";
import { BACKEND_BASE_URL, NODE_ENV } from "../lib/constant";
import Select from "react-select";

interface Recipe {
  element_one: string;
  element_two: string;
  target_element_name: string;
}

interface Element {
  name: string;
  image_path: string;
  recipes_to_make_this_element: Recipe[];
  recipes_to_make_other_element: Recipe[];
  is_visited: boolean;
  tier: number;
  // Additional properties for UI use
  Name?: string;
  ImagePath?: string;
  Recipes?: string[][];
}

export default function Wiki() {
  const [elements, setElements] = useState<Element[]>([]);
  const [search, setSearch] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeCategory, setActiveCategory] = useState("all");
  const [viewMode, setViewMode] = useState<"grid" | "list">("grid");
  const [sortBy, setSortBy] = useState<"name" | "tier">("name");
  const [selectedElement, setSelectedElement] = useState<Element | null>(null);

  const isElementBasic = (elementName: string): boolean => {
    return ["Air", "Earth", "Fire", "Water"].includes(elementName);
  };

  const getBasicElementDescription = (): string => {
    return "Available from the start";
  };

  // Fetch elements data from backend
  useEffect(() => {
    const fetchElements = async () => {
      try {
        const response = await Axios(`${NODE_ENV === "production" ? "https" : "http"}://${BACKEND_BASE_URL}/api/elements`);
        if (response.status < 200 || response.status >= 300) {
          throw new Error(`Error: ${response.status}`);
        }
        const data = response.data;
        const processedElements = data.map((element: Element) => {
          const recipesAsArrays = element.recipes_to_make_this_element.map((recipe) => [recipe.element_one, recipe.element_two]);

          const isBasicElement = isElementBasic(element.name);

          let basicDescription = "";
          if (isBasicElement) {
            basicDescription = getBasicElementDescription();
          }

          const isTimeElement = element.name === "Time";

          return {
            ...element,
            Name: element.name,
            ImagePath: element.image_path,
            Recipes: recipesAsArrays,
            is_basic: isBasicElement,
            basic_description: basicDescription,
            is_special: isTimeElement,
            special_unlock_requirement: isTimeElement ? "Unlock 100 elements" : "",
          };
        });

        if (processedElements.some((e: Element) => e.name === "Time")) {
          const unlockedElementsCount: number = processedElements.filter((e: Element) => e.is_visited).length;
          const timeElement: Element | undefined = processedElements.find((e: Element) => e.name === "Time");
          if (timeElement) {
            timeElement.is_visited = unlockedElementsCount >= 100;
          }
        }

        setElements(processedElements);
        setLoading(false);
      } catch (error) {
        console.error("Failed to fetch elements:", error);
        setError("Failed to load elements. Please try again later.");
        setLoading(false);
      }
    };

    fetchElements();
  }, []);

  const elementOptions = elements
    .map((element) => ({
      value: element.name,
      label: element.name,
      image: element.image_path,
    }))
    .sort((a, b) => a.label.localeCompare(b.label));

  const isElementConsideredBasic = (element: Element): boolean => {
    if (isElementBasic(element.name)) return true;

    if (element.name === "Time") return true;

    return element.recipes_to_make_this_element.length === 0;
  };

  const filteredElements = elements.filter((element) => {
    const matchesSearch = element.name.toLowerCase().includes(search.toLowerCase());
    let matchesCategory = false;

    if (activeCategory === "all") {
      matchesCategory = true;
    } else if (activeCategory === "basic") {
      matchesCategory = isElementConsideredBasic(element);
    } else if (activeCategory === "advanced") {
      matchesCategory = !isElementConsideredBasic(element) && element.tier > 0;
    } else if (activeCategory.startsWith("tier-")) {
      matchesCategory = element.tier === parseInt(activeCategory.split("-")[1]);
    }

    return matchesSearch && matchesCategory;
  });

  const sortedElements = [...filteredElements].sort((a, b) => {
    if (sortBy === "name") {
      return a.name.localeCompare(b.name);
    } else {
      return a.tier - b.tier;
    }
  });
  const tiers = [...new Set(elements.map((element) => element.tier))].sort();

  const basicElementsCount = elements.filter(isElementConsideredBasic).length;
  const advancedElementsCount = elements.filter((element) => !isElementConsideredBasic(element)).length;

  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 20;

  const indexOfLastItem = currentPage * itemsPerPage;
  const indexOfFirstItem = indexOfLastItem - itemsPerPage;
  const currentItems = sortedElements.slice(indexOfFirstItem, indexOfLastItem);
  const totalPages = Math.ceil(sortedElements.length / itemsPerPage);

  const pageNumbers = [];
  for (let i = 1; i <= totalPages; i++) {
    pageNumbers.push(i);
  }

  return (
    <div className="flex min-h-screen bg-gray-100">
      {/* Sidebar */}
      <div className="w-64 bg-white shadow-md p-4 hidden md:block">
        <h2 className="text-xl font-bold mb-4">Elements Encyclopedia</h2>

        <div className="mb-6">
          <h3 className="font-medium text-gray-700 mb-2">Categories</h3>
          <ul>
            <li>
              <button
                className={`w-full text-left py-2 px-3 rounded ${activeCategory === "all" ? "bg-blue-100 text-blue-700" : "hover:bg-gray-100"}`}
                onClick={() => {
                  setActiveCategory("all");
                  setCurrentPage(1);
                }}>
                All Elements ({elements.length})
              </button>
            </li>
            <li>
              <button
                className={`w-full text-left py-2 px-3 rounded ${activeCategory === "basic" ? "bg-blue-100 text-blue-700" : "hover:bg-gray-100"}`}
                onClick={() => {
                  setActiveCategory("basic");
                  setCurrentPage(1);
                }}>
                Basic Elements ({basicElementsCount})
              </button>
            </li>
            <li>
              <button
                className={`w-full text-left py-2 px-3 rounded ${activeCategory === "advanced" ? "bg-blue-100 text-blue-700" : "hover:bg-gray-100"}`}
                onClick={() => {
                  setActiveCategory("advanced");
                  setCurrentPage(1);
                }}>
                Advanced Elements ({advancedElementsCount})
              </button>
            </li>
          </ul>
        </div>

        {tiers.length > 0 && (
          <div className="mb-6">
            <h3 className="font-medium text-gray-700 mb-2">Tiers</h3>
            <ul>
              {tiers
                .sort((a, b) => a - b)
                .map((tier) => (
                  <li key={tier}>
                    <button
                      className={`w-full text-left py-2 px-3 rounded ${activeCategory === `tier-${tier}` ? "bg-blue-100 text-blue-700" : "hover:bg-gray-100"}`}
                      onClick={() => {
                        setActiveCategory(`tier-${tier}`);
                        setCurrentPage(1);
                      }}>
                      Tier {tier} ({elements.filter((e) => e.tier === tier).length})
                    </button>
                  </li>
                ))}
            </ul>
          </div>
        )}
      </div>
      {/* Main content */}
      <div className="flex-1 p-6">
        {/* Desktop header */}
        <div className="hidden md:block mb-6">
          <h1 className="text-3xl font-bold text-gray-900">Elements Encyclopedia Wiki</h1>
          <p className="text-gray-600">Discover and learn about all the elements in your world</p>
        </div>

        {/* Controls */}
        <div className="mb-6 flex flex-col sm:flex-row items-stretch sm:items-center space-y-3 sm:space-y-0 sm:space-x-4">
          <div className="relative flex-grow">
            <Select
              options={elementOptions}
              placeholder="Search elements..."
              value={elementOptions.find((opt) => opt.value === search) || null}
              onChange={(selected) => {
                if (selected) {
                  setSearch(selected.value);
                } else {
                  setSearch("");
                }
                setCurrentPage(1);
              }}
              classNamePrefix="element-select"
              className="w-full"
              isSearchable={true}
              isClearable={true}
            />
          </div>

          <div className="flex items-center space-x-2">
            <select value={sortBy} onChange={(e) => setSortBy(e.target.value as "name" | "tier")} className="py-2 px-3 border border-gray-300 rounded-lg focus:ring-blue-500 focus:border-blue-500">
              <option value="name">Sort by Name</option>
              <option value="tier">Sort by Tier</option>
            </select>

            <div className="flex border border-gray-300 rounded-lg overflow-hidden">
              <button className={`px-3 py-2 ${viewMode === "grid" ? "bg-blue-100 text-blue-700" : "bg-white text-gray-700"}`} onClick={() => setViewMode("grid")}>
                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                  <path d="M5 3a2 2 0 00-2 2v2a2 2 0 002 2h2a2 2 0 002-2V5a2 2 0 00-2-2H5zM5 11a2 2 0 00-2 2v2a2 2 0 002 2h2a2 2 0 002-2v-2a2 2 0 00-2-2H5zM11 5a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V5zM11 13a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
                </svg>
              </button>
              <button className={`px-3 py-2 ${viewMode === "list" ? "bg-blue-100 text-blue-700" : "bg-white text-gray-700"}`} onClick={() => setViewMode("list")}>
                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                  <path
                    fillRule="evenodd"
                    d="M3 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z"
                    clipRule="evenodd"
                  />
                </svg>
              </button>
            </div>
          </div>
        </div>

        {/* Loading state */}
        {loading ? (
          <div className="flex justify-center py-20">
            <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
          </div>
        ) : error ? (
          <div className="bg-red-50 border-l-4 border-red-400 p-4 rounded shadow">
            <div className="flex">
              <div className="flex-shrink-0">
                <svg className="h-5 w-5 text-red-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path
                    fillRule="evenodd"
                    d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                    clipRule="evenodd"
                  />
                </svg>
              </div>
              <div className="ml-3">
                <p className="text-red-700">{error}</p>
                <button className="mt-2 text-sm text-red-600 hover:text-red-800 font-medium" onClick={() => window.location.reload()}>
                  Try again
                </button>
              </div>
            </div>
          </div>
        ) : sortedElements.length === 0 ? (
          <div className="bg-white rounded-lg shadow p-8 text-center">
            <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <h3 className="mt-2 text-lg font-medium text-gray-900">No elements found</h3>
            <p className="mt-1 text-gray-500">{search ? `No elements matching "${search}"` : "Try adjusting your filter criteria."}</p>
          </div>
        ) : (
          <>
            {/* Grid View */}
            {viewMode === "grid" && (
              <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
                {currentItems.map((element) => (
                  <div key={element.name} className="bg-white rounded-lg shadow overflow-hidden hover:shadow-md transition-shadow cursor-pointer" onClick={() => setSelectedElement(element)}>
                    <div className="h-32 p-4 flex items-center justify-center bg-gray-50">
                      <img src={element.image_path} alt={element.name} className="max-h-full max-w-full object-contain" />
                    </div>
                    <div className="p-4 border-t flex flex-row justify-between items-center">
                      <h3 className="font-medium text-gray-900 truncate" title={element.name}>
                        {element.name.length > 20 ? `${element.name.slice(0, 15)}...` : element.name}
                      </h3>
                      <span className="text-xs px-2 py-1 bg-blue-100 text-blue-800 rounded-full">Tier {element.tier}</span>
                    </div>
                  </div>
                ))}
              </div>
            )}

            {/* List View */}
            {viewMode === "list" && (
              <div className="bg-white shadow rounded-lg overflow-hidden">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Element
                      </th>
                      <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Tier
                      </th>

                      <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Recipes
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {currentItems.map((element) => (
                      <tr key={element.name} className="hover:bg-gray-50 cursor-pointer" onClick={() => setSelectedElement(element)}>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="flex items-center">
                            <div className="flex-shrink-0 h-10 w-10">
                              <img className="h-10 w-10 object-contain" src={element.image_path} alt={element.name} />
                            </div>
                            <div className="ml-4">
                              <div className="text-sm font-medium text-gray-900">{element.name}</div>
                            </div>
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm text-gray-900">Tier {element.tier}</div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm text-gray-900">
                            {element.recipes_to_make_this_element && element.recipes_to_make_this_element.length > 0 ? `${element.recipes_to_make_this_element.length} recipe(s)` : "Basic element"}
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
            {/* Pagination Controls */}
            {sortedElements.length > itemsPerPage && (
              <div className="mt-6 flex justify-center">
                <nav className="flex items-center space-x-1">
                  <button
                    onClick={() => setCurrentPage((prev) => Math.max(prev - 1, 1))}
                    disabled={currentPage === 1}
                    className={`px-3 py-1 rounded border ${currentPage === 1 ? "text-gray-400 border-gray-200 cursor-not-allowed" : "text-gray-700 border-gray-300 hover:bg-gray-50"}`}>
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                      <path fillRule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clipRule="evenodd" />
                    </svg>
                  </button>

                  {pageNumbers.map((number) => {
                    // Show limited page numbers with ellipsis for better UX
                    if (number === 1 || number === totalPages || (number >= currentPage - 1 && number <= currentPage + 1)) {
                      return (
                        <button
                          key={number}
                          onClick={() => setCurrentPage(number)}
                          className={`px-3 py-1 rounded ${currentPage === number ? "bg-blue-100 text-blue-700 border border-blue-300" : "text-gray-700 border border-gray-300 hover:bg-gray-50"}`}>
                          {number}
                        </button>
                      );
                    } else if (number === currentPage - 2 || number === currentPage + 2) {
                      return (
                        <span key={number} className="px-2">
                          ...
                        </span>
                      );
                    }
                    return null;
                  })}

                  <button
                    onClick={() => setCurrentPage((prev) => Math.min(prev + 1, totalPages))}
                    disabled={currentPage === totalPages || totalPages === 0}
                    className={`px-3 py-1 rounded border ${
                      currentPage === totalPages || totalPages === 0 ? "text-gray-400 border-gray-200 cursor-not-allowed" : "text-gray-700 border-gray-300 hover:bg-gray-50"
                    }`}>
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                      <path fillRule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clipRule="evenodd" />
                    </svg>
                  </button>
                </nav>
              </div>
            )}
            {/* Page Information */}
            {sortedElements.length > 0 && (
              <div className="mt-4 text-center text-sm text-gray-600">
                Showing {indexOfFirstItem + 1}-{Math.min(indexOfLastItem, sortedElements.length)} of {sortedElements.length} elements
              </div>
            )}
          </>
        )}
      </div>

      {/* Element Detail Modal */}
      {selectedElement && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-lg shadow-lg max-w-4xl w-full max-h-[90vh] overflow-y-auto">
            <div className="p-6">
              {/* Header */}
              <div className="flex justify-between items-start mb-6">
                <div className="flex items-center">
                  <img src={selectedElement.image_path} alt={selectedElement.name} className="w-24 h-24 object-contain bg-gray-200 p-3 rounded-lg mr-4 shadow-sm" />
                  <div>
                    <h2 className="text-2xl font-bold text-gray-900">{selectedElement.name}</h2>
                    <div className="flex flex-wrap gap-2 mt-2">
                      <span className="text-sm px-2 py-1 bg-blue-100 text-blue-800 rounded-full">Tier {selectedElement.tier}</span>
                      {isElementConsideredBasic(selectedElement) && <span className="text-sm px-2 py-1 bg-purple-100 text-purple-800 rounded-full">Basic Element</span>}
                    </div>
                  </div>
                </div>
                <button onClick={() => setSelectedElement(null)} className="text-gray-400 hover:text-gray-500">
                  <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>

              {/* Tabs for Recipes */}
              <div className="border-b border-gray-200 mb-6">
                <div className="flex -mb-px">
                  <button
                    className={`py-2 px-4 text-sm font-medium border-b-2 ${
                      selectedElement.recipes_to_make_this_element.length > 0 ? "border-blue-500 text-blue-600" : "border-transparent text-gray-500 hover:text-gray-700"
                    }`}>
                    How to Create ({selectedElement.recipes_to_make_this_element.length})
                  </button>
                  <button
                    className={`py-2 px-4 text-sm font-medium border-b-2 ${
                      selectedElement.recipes_to_make_other_element.length > 0 ? "border-transparent text-gray-500 hover:text-gray-700" : "border-transparent text-gray-400"
                    }`}>
                    Used In Other Recipes ({selectedElement.recipes_to_make_other_element.length})
                  </button>
                </div>
              </div>

              {/* Recipes Section */}
              {selectedElement.recipes_to_make_this_element.length > 0 ? (
                <div className="mt-6">
                  <h3 className="text-lg font-medium text-gray-900 mb-4">How to Create This Element</h3>

                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    {selectedElement.recipes_to_make_this_element.map((recipe, index) => {
                      const ingredient1 = elements.find((e) => e.name.toLowerCase() === recipe.element_one.toLowerCase());
                      const ingredient2 = elements.find((e) => e.name.toLowerCase() === recipe.element_two.toLowerCase());

                      return (
                        <div key={index} className="bg-gray-50 p-4 rounded-lg border border-gray-200 hover:shadow-md transition-shadow">
                          <div className="flex items-center justify-center">
                            {/* First ingredient */}
                            <div className="flex flex-col items-center">
                              <div className="w-16 h-16 bg-white rounded-lg shadow-sm p-2 flex items-center justify-center">
                                {ingredient1 ? (
                                  <img
                                    src={ingredient1.image_path}
                                    alt={recipe.element_one}
                                    className="max-h-full max-w-full object-contain"
                                    onClick={(e) => {
                                      e.stopPropagation();
                                      setSelectedElement(ingredient1);
                                    }}
                                    title={`View ${recipe.element_one} details`}
                                    style={{ cursor: "pointer" }}
                                  />
                                ) : (
                                  <div className="text-sm text-center text-gray-500">{recipe.element_one}</div>
                                )}
                              </div>
                              <span className="text-xs font-medium mt-1 text-center max-w-[80px] truncate" title={recipe.element_one}>
                                {recipe.element_one}
                              </span>
                            </div>

                            {/* Plus sign */}
                            <div className="mx-3 text-gray-400">
                              <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                              </svg>
                            </div>

                            {/* Second ingredient */}
                            <div className="flex flex-col items-center">
                              <div className="w-16 h-16 bg-white rounded-lg shadow-sm p-2 flex items-center justify-center">
                                {ingredient2 ? (
                                  <img
                                    src={ingredient2.image_path}
                                    alt={recipe.element_two}
                                    className="max-h-full max-w-full object-contain"
                                    onClick={(e) => {
                                      e.stopPropagation();
                                      setSelectedElement(ingredient2);
                                    }}
                                    title={`View ${recipe.element_two} details`}
                                    style={{ cursor: "pointer" }}
                                  />
                                ) : (
                                  <div className="text-sm text-center text-gray-500">{recipe.element_two}</div>
                                )}
                              </div>
                              <span className="text-xs font-medium mt-1 text-center max-w-[80px] truncate" title={recipe.element_two}>
                                {recipe.element_two}
                              </span>
                            </div>

                            {/* Equals sign */}
                            <div className="mx-3 text-gray-400">
                              <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 12H4" />
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 6H4" />
                              </svg>
                            </div>

                            {/* Result (current element) */}
                            <div className="flex flex-col items-center">
                              <div className="w-16 h-16 bg-blue-50 border-2 border-blue-200 rounded-lg p-2 flex items-center justify-center">
                                <img src={selectedElement.image_path} alt={selectedElement.name} className="max-h-full max-w-full object-contain" />
                              </div>
                              <span className="text-xs font-medium mt-1 text-blue-700 text-center">{selectedElement.name}</span>
                            </div>
                          </div>
                        </div>
                      );
                    })}
                  </div>
                </div>
              ) : (
                <div className="mt-4 bg-yellow-50 p-4 rounded-lg border border-yellow-200">
                  <div className="flex items-start">
                    <div className="flex-shrink-0">
                      <svg className="h-5 w-5 text-yellow-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                        <path
                          fillRule="evenodd"
                          d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
                          clipRule="evenodd"
                        />
                      </svg>
                    </div>
                    <div className="ml-3">
                      <h3 className="text-sm font-medium text-yellow-800">
                        {["Air", "Earth", "Fire", "Water"].includes(selectedElement.Name || "")
                          ? "Starting Element"
                          : ["Time"].includes(selectedElement.Name || "")
                          ? "Special Element"
                          : "No Recipe Available"}
                      </h3>
                      <p className="text-sm text-yellow-700 mt-1">This element cannot be created from other elements.</p>
                      {selectedElement.Name === "Time" && <p className="text-sm text-yellow-700 mt-1">Unlock 100 elements to discover the Time element.</p>}
                    </div>
                  </div>
                </div>
              )}

              {/* Used in other recipes section */}
              {selectedElement.recipes_to_make_other_element.length > 0 && (
                <div className="mt-8">
                  <h3 className="text-lg font-medium text-gray-900 mb-4">Used In Other Recipes</h3>
                  <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                      {selectedElement.recipes_to_make_other_element.map((recipe, index) => (
                        <div key={index} className="flex items-center justify-between p-2 border-b border-gray-200 last:border-b-0">
                          <div className="flex items-center">
                            <span className="text-sm font-medium">{selectedElement.name}</span>
                            <span className="mx-2 text-gray-400">+</span>
                            <span className="text-sm font-medium">{recipe.element_one === selectedElement.name ? recipe.element_two : recipe.element_one}</span>
                          </div>
                          <div className="flex items-center">
                            <span className="mr-2 text-gray-400">=</span>
                            <span className="text-sm font-medium text-blue-600">{recipe.target_element_name}</span>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                </div>
              )}

              <div className="mt-6 pt-4 border-t border-gray-200 flex justify-end">
                <button
                  onClick={() => setSelectedElement(null)}
                  className="inline-flex justify-center px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
                  Close
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
