import { useEffect, useRef, useState } from "react";
import { BACKEND_BASE_URL } from "../lib/contant";

const Home = () => {
  const socketRef = useRef<WebSocket | null>(null);
  const [messages, setMessages] = useState<string[]>([]);

  useEffect(() => {
    // Connect to WebSocket
    socketRef.current = new WebSocket(`ws://${BACKEND_BASE_URL}/ws`);

    socketRef.current.onopen = () => {
      console.log("WebSocket connection opened");
      socketRef.current?.send("Hello from frontend!");
    };

    socketRef.current.onmessage = (event) => {
      console.log("Received from server:", event.data);
      setMessages((prev) => [...prev, event.data]);
    };

    socketRef.current.onerror = (err) => {
      console.error("WebSocket error:", err);
    };

    socketRef.current.onclose = () => {
      console.log("WebSocket connection closed");
    };

    // Cleanup on unmount
    return () => {
      socketRef.current?.close();
    };
  }, []);

  return (
    <div>
      <h1>WebSocket Client</h1>
      <ul>
        {messages.map((msg, idx) => (
          <li key={idx}>{msg}</li>
        ))}
      </ul>
    </div>
  );
};

export default Home;
