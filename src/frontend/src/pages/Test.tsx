import { useEffect, useRef, useState } from "react";

interface RecipeTreeNode {
  name: string;
  image_path: string;
  element_1?: RecipeTreeNode;
  element_2?: RecipeTreeNode;
  minimum_nodes_recipe_tree: number;
}

interface TreeUpdate {
  best_tree: RecipeTreeNode | null;
  exploring_tree: RecipeTreeNode | null;
}

const BACKEND_BASE_URL = import.meta.env.VITE_PUBLIC_BACKEND_BASE_URL || "localhost:4000";

export default function RecipeTreeVisualizer() {
  const [target, setTarget] = useState("Water");
  const [delayMs, setDelayMs] = useState(500);
  const [findBestTree, setFindBestTree] = useState(true);
  const [maxTreeCount, setMaxTreeCount] = useState(0);
  const [mode, setMode] = useState("bfs"); // New state for mode selection
  const [messages, setMessages] = useState<TreeUpdate[]>([]);
  const wsRef = useRef<WebSocket | null>(null);

  const connectWebSocket = () => {
    const effectiveMaxTreeCount = findBestTree ? 0 : maxTreeCount;

    const ws = new WebSocket(`ws://${BACKEND_BASE_URL}/ws`);
    wsRef.current = ws;

    ws.onopen = () => {
      console.log("WebSocket connected");

      const payload = {
        target,
        mode,
        find_best_tree: findBestTree,
        max_tree_count: effectiveMaxTreeCount,
        delay_ms: delayMs,
      };

      ws.send(JSON.stringify(payload));
    };

    ws.onmessage = (event: MessageEvent<string>) => {
      try {
        const data: TreeUpdate = JSON.parse(event.data);
        setMessages((prev) => [data, ...prev]);
      } catch (error) {
        console.error("Error parsing message:", error);
      }
    };

    ws.onerror = (error) => console.error("WebSocket error:", error);
    ws.onclose = () => console.log("WebSocket closed");
  };

  useEffect(() => {
    return () => {
      if (wsRef.current) wsRef.current.close();
    };
  }, []);

  return (
    <div className="p-4 max-w-3xl mx-auto">
      <h1 className="text-xl font-bold mb-4">Recipe Tree Visualizer</h1>

      <div className="flex flex-col gap-4 mb-4">
        <input type="text" value={target} onChange={(e) => setTarget(e.target.value)} className="border p-2 rounded" placeholder="Target element (e.g., Water)" />

        <div className="flex items-center gap-2">
          <input
            type="checkbox"
            checked={findBestTree}
            onChange={(e) => {
              const checked = e.target.checked;
              setFindBestTree(checked);
              setMaxTreeCount(checked ? 0 : 10); // default valid value
            }}
          />
          <label className="text-sm">Find Best Tree</label>
        </div>

        <input type="number" value={maxTreeCount} onChange={(e) => setMaxTreeCount(Number(e.target.value))} className="border p-2 rounded" disabled={findBestTree} placeholder="Max Tree Count" />

        <input type="number" value={delayMs} onChange={(e) => setDelayMs(Number(e.target.value))} className="border p-2 rounded" placeholder="Delay (ms)" />

        {/* Mode selection dropdown */}
        <select value={mode} onChange={(e) => setMode(e.target.value)} className="border p-2 rounded">
          <option value="bfs">BFS</option>
          <option value="dfs">DFS</option>
          <option value="bidirectional">Bidirectional</option>
        </select>

        <button onClick={connectWebSocket} className="bg-blue-600 text-white px-4 py-2 rounded">
          Start
        </button>
      </div>

      <div className="space-y-2 max-h-[400px] overflow-y-auto">
        {messages.map((msg, idx) => (
          <pre key={idx} className="bg-gray-100 p-2 rounded text-sm overflow-x-auto">
            {JSON.stringify(msg, null, 2)}
          </pre>
        ))}
      </div>
    </div>
  );
}
