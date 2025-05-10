import { BrowserRouter as Router, Routes, Route } from "react-router";
import Visualizer from "./pages/Visualizer";
import Wiki from "./pages/Wiki";
import Home from "./pages/Home";
import Layout from "./pages/Layout";
import AboutUs from "./pages/AboutUs";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Layout><Home /></Layout>} />
        <Route path="/Wiki" element={<Layout><Wiki /></Layout>} />
        <Route path="/Visualizer" element={<Layout><Visualizer /></Layout>} />
        <Route path="/AboutUs" element={<Layout><AboutUs /></Layout>} />
      </Routes>
    </Router>
  );
}

export default App;
