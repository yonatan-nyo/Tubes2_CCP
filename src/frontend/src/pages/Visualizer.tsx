import { useEffect, useRef, useState } from "react";

interface RecipeTreeNode {
  name: string;
  image_path: string;
  element_1?: RecipeTreeNode;
  element_2?: RecipeTreeNode;
}

const BACKEND_BASE_URL = import.meta.env.VITE_PUBLIC_BACKEND_BASE_URL || "localhost:4000";

function TreeNode({ node }: { node: RecipeTreeNode }) {
  if (!node) return null;

  return (
    <div className="ml-1 border-l border-gray-200 text-[10px] leading-tight rounded-lg p-2">
      <div className="font-semibold">{node.name}</div>
      <img src={node.image_path} alt={node.name} className="w-6 h-6 my-0.5" />
      <div className="flex gap-1">
        {node.element_1 && <TreeNode node={node.element_1} />}
        {node.element_2 && <TreeNode node={node.element_2} />}
      </div>
    </div>
  );
}

export default function Visualizer() {
  const [target, setTarget] = useState("Water");
  const [delayMs, setDelayMs] = useState(500);
  const [maxTreeCount, setMaxTreeCount] = useState(1);
  const [mode, setMode] = useState("bfs");

  const [exploringTree, setExploringTree] = useState<RecipeTreeNode | null>(null);
  const [finalTrees, setFinalTrees] = useState<RecipeTreeNode[]>([]);
  const [selectedTab, setSelectedTab] = useState(0);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [searchStats, setSearchStats] = useState({ durationMs: 0, nodesExplored: 0 });

  const wsRef = useRef<WebSocket | null>(null);

  const connectWebSocket = () => {
    setFinalTrees([]);
    setSearchStats({ durationMs: 0, nodesExplored: 0 });
    setIsLoading(true);
    setError(null);

    const ws = new WebSocket(`ws://${BACKEND_BASE_URL}/ws`);
    wsRef.current = ws;

    ws.onopen = () => {
      ws.send(JSON.stringify({ target, mode, max_tree_count: maxTreeCount, delay_ms: delayMs }));
    };

    ws.onmessage = (event) => {
      try {
        const parsed = JSON.parse(event.data);

        if (parsed.exploring_tree) {
          setExploringTree(parsed.exploring_tree);
        }

        if (typeof parsed.duration_ms === "number" && typeof parsed.nodes_explored === "number") {
          setSearchStats({
            durationMs: parsed.duration_ms,
            nodesExplored: parsed.nodes_explored,
          });
        }

        if (Array.isArray(parsed)) {
          setFinalTrees(parsed);
          setExploringTree(null);
          setSelectedTab(0);
          setIsLoading(false);
        }

        if (parsed.error) {
          setError(parsed.error);
          console.error("Backend error:", parsed.error);
          setIsLoading(false);
        }
      } catch (err) {
        setError("Invalid response from server");
        console.error("Invalid message:", err);
        setIsLoading(false);
      }
    };

    ws.onerror = (err) => {
      setError("WebSocket connection error");
      console.error("WebSocket error:", err);
      setIsLoading(false);
    };

    ws.onclose = () => {
      console.log("WebSocket closed");
      setIsLoading(false);
    };
  };

  useEffect(() => {
    return () => {
      wsRef.current?.close();
    };
  }, []);

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="p-6 max-w-7xl mx-auto space-y-6">
        <div className="bg-white shadow-md rounded-lg p-6">
          <h1 className="text-2xl font-bold text-blue-800 mb-6">Recipe Tree Visualizer</h1>

          <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Target Item</label>
              <input
                value={target}
                onChange={(e) => setTarget(e.target.value)}
                placeholder="Target"
                className="w-full border border-gray-300 p-2 rounded text-sm focus:ring-blue-500 focus:border-blue-500"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Search Algorithm</label>
              <select
                value={mode}
                onChange={(e) => setMode(e.target.value)}
                className="w-full border border-gray-300 p-2 rounded text-sm focus:ring-blue-500 focus:border-blue-500"
              >
                <option value="bfs">BFS</option>
                <option value="dfs">DFS</option>
                <option value="bidirectional">Bidirectional</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Max Trees</label>
              <input
                type="number"
                value={maxTreeCount}
                onChange={(e) => setMaxTreeCount(Number(e.target.value))}
                className="w-full border border-gray-300 p-2 rounded text-sm focus:ring-blue-500 focus:border-blue-500"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Delay (ms)</label>
              <input
                type="number"
                value={delayMs}
                onChange={(e) => setDelayMs(Number(e.target.value))}
                className="w-full border border-gray-300 p-2 rounded text-sm focus:ring-blue-500 focus:border-blue-500"
              />
            </div>
          </div>

          <button
            onClick={connectWebSocket}
            disabled={isLoading}
            className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition-colors disabled:bg-blue-300 w-full"
          >
            {isLoading ? "Processing..." : "Start"}
          </button>

          {error && <div className="mt-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">{error}</div>}
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="bg-gray-50 shadow-md rounded-lg overflow-hidden">
            <div className="bg-green-100 p-4">
              <h2 className="font-bold text-lg text-green-800">Final Trees</h2>
            </div>
            <div className="p-4">
              {finalTrees.length === 0 ? (
                <div className="text-gray-500 text-center py-8">
                  <svg className="w-12 h-12 mx-auto mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={1.5}
                      d="M9 13h6m-3-3v6m-9 1V7a2 2 0 012-2h6l2 2h6a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2z"
                    />
                  </svg>
                  <p className="text-sm">No trees received yet. Start a search to generate results.</p>
                </div>
              ) : (
                <>
                  <div className="flex flex-wrap gap-2 mb-3">
                    {finalTrees.map((_, idx) => (
                      <button
                        key={idx}
                        onClick={() => setSelectedTab(idx)}
                        className={`px-3 py-1.5 text-sm rounded-full transition-colors ${
                          selectedTab === idx ? "bg-green-600 text-white" : "bg-gray-200 hover:bg-gray-300"
                        }`}
                      >
                        Tree {idx + 1}
                      </button>
                    ))}
                  </div>
                  <div className="overflow-x-auto bg-green-200/30 max-h-96">
                    <TreeNode node={finalTrees[selectedTab]} />
                  </div>
                </>
              )}
            </div>
          </div>

          <div className="bg-gray-50 shadow-md rounded-lg overflow-hidden">
            <div className="bg-yellow-100 p-4">
              <h2 className="font-bold text-lg text-yellow-800">Exploring Tree</h2>
            </div>
            <div className="p-4">
              {isLoading && !exploringTree ? (
                <div className="text-center py-8">
                  <div className="animate-spin rounded-full h-10 w-10 border-t-2 border-b-2 border-blue-500 mx-auto mb-2"></div>
                  <p className="text-gray-500 text-sm">Searching...</p>
                </div>
              ) : exploringTree ? (
                <div className="overflow-x-auto bg-yellow-200/30 max-h-96">
                  <TreeNode node={exploringTree} />
                </div>
              ) : (
                <div className="text-gray-500 text-center py-8">
                  <svg className="w-12 h-12 mx-auto mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={1.5}
                      d="M8 16l2.879-2.879m0 0a3 3 0 104.243-4.242 3 3 0 00-4.243 4.242zM21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                  <p className="text-sm">Waiting for exploration updates...</p>
                </div>
              )}
            </div>
          </div>
        </div>

        <div className="mt-6 bg-white shadow-md rounded-lg p-4 flex justify-center gap-8 text-sm text-gray-700">
          <div>
            ⏱ <strong>Time Elapsed:</strong> {searchStats.durationMs} ms
          </div>
          <div>
            🔍 <strong>Nodes Explored:</strong> {searchStats.nodesExplored}
          </div>
        </div>
      </div>
    </div>
  );
}
