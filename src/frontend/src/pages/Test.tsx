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
    <div className="ml-1 border-l border-gray-300 pl-1 text-[10px] leading-tight">
      <div className="font-semibold">{node.name}</div>
      <img src={node.image_path} alt={node.name} className="w-6 h-6 my-0.5" />
      <div className="flex gap-1">
        {node.element_1 && <TreeNode node={node.element_1} />}
        {node.element_2 && <TreeNode node={node.element_2} />}
      </div>
    </div>
  );
}

export default function RecipeTreeVisualizer() {
  const [target, setTarget] = useState("Water");
  const [delayMs, setDelayMs] = useState(500);
  const [maxTreeCount, setMaxTreeCount] = useState(1);
  const [mode, setMode] = useState("bfs");

  const [exploringTree, setExploringTree] = useState<RecipeTreeNode | null>(null);
  const [finalTrees, setFinalTrees] = useState<RecipeTreeNode[]>([]);
  const [selectedTab, setSelectedTab] = useState(0);

  const wsRef = useRef<WebSocket | null>(null);

  const connectWebSocket = () => {
    setFinalTrees([]);
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
        } else if (Array.isArray(parsed)) {
          setFinalTrees(parsed);
          setExploringTree(null);
          setSelectedTab(0);
        } else if (parsed.error) {
          console.error("Backend error:", parsed.error);
        }
      } catch (err) {
        console.error("Invalid message:", err);
      }
    };

    ws.onerror = (err) => console.error("WebSocket error:", err);
    ws.onclose = () => console.log("WebSocket closed");
  };

  useEffect(() => {
    return () => {
      wsRef.current?.close();
    };
  }, []);

  return (
    <div className="p-4 max-w-7xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold">Recipe Tree Visualizer</h1>

      {/* Controls */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <input value={target} onChange={(e) => setTarget(e.target.value)} placeholder="Target" className="border p-2 rounded text-sm" />
        <select value={mode} onChange={(e) => setMode(e.target.value)} className="border p-2 rounded text-sm">
          <option value="bfs">BFS</option>
          <option value="dfs">DFS</option>
          <option value="bidirectional">Bidirectional</option>
        </select>
        <input type="number" value={maxTreeCount} onChange={(e) => setMaxTreeCount(Number(e.target.value))} className="border p-2 rounded text-sm" placeholder="Max Trees" />
        <input type="number" value={delayMs} onChange={(e) => setDelayMs(Number(e.target.value))} className="border p-2 rounded text-sm" placeholder="Delay (ms)" />
      </div>

      <button onClick={connectWebSocket} className="bg-blue-600 text-white px-4 py-2 rounded">
        Start
      </button>

      {/* Main Tree View Split */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Left: Final Trees with Tabs */}
        <div className="bg-green-100 p-4 rounded shadow">
          <h2 className="font-bold text-lg mb-2">Final Trees</h2>

          {finalTrees.length === 0 ? (
            <p className="text-gray-600 text-sm">No trees received yet.</p>
          ) : (
            <>
              <div className="flex flex-wrap gap-2 mb-2">
                {finalTrees.map((_, idx) => (
                  <button key={idx} onClick={() => setSelectedTab(idx)} className={`px-2 py-1 text-sm rounded border ${selectedTab === idx ? "bg-white font-bold" : "bg-gray-200"}`}>
                    Tree {idx + 1}
                  </button>
                ))}
              </div>
              <div className="bg-white border p-2 rounded overflow-x-auto">
                <TreeNode node={finalTrees[selectedTab]} />
              </div>
            </>
          )}
        </div>

        {/* Right: Exploring Tree */}
        <div className="bg-yellow-100 p-4 rounded shadow overflow-x-auto">
          <h2 className="font-bold text-lg mb-2">Exploring Tree</h2>
          {exploringTree ? <TreeNode node={exploringTree} /> : <p className="text-gray-600 text-sm">Waiting for updates...</p>}
        </div>
      </div>
    </div>
  );
}
