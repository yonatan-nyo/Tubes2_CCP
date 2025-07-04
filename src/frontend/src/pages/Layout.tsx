import { type ReactNode } from "react";
import { Link, useLocation } from "react-router";

function Navbar() {
  const location = useLocation();
  
  return (
    <nav className="bg-blue-700 text-white shadow-md">
      <div className="w-full mx-auto px-4">
        <div className="flex justify-between items-center h-16">
          <div className="flex-shrink-0 flex items-center">
            <span className="font-bold text-xl">CCP Little Alchemy 2</span>
          </div>
          <div className="flex space-x-4">
            <Link 
              to="/" 
              className={`px-3 py-2 rounded-md text-sm font-medium ${
                location.pathname === '/' 
                  ? 'bg-blue-900 text-white' 
                  : 'text-blue-100 hover:bg-blue-600'
              }`}
            >
              Home
            </Link>
            <Link 
              to="/Wiki" 
              className={`px-3 py-2 rounded-md text-sm font-medium ${
                location.pathname === '/Wiki' 
                  ? 'bg-blue-900 text-white' 
                  : 'text-blue-100 hover:bg-blue-600'
              }`}
            >
              Wiki
            </Link>
            <Link 
              to="/Visualizer" 
              className={`px-3 py-2 rounded-md text-sm font-medium ${
                location.pathname === '/Visualizer' 
                  ? 'bg-blue-900 text-white' 
                  : 'text-blue-100 hover:bg-blue-600'
              }`}
            >
              Visualizer
            </Link>
            <Link 
              to="/AboutUs" 
              className={`px-3 py-2 rounded-md text-sm font-medium ${
                location.pathname === '/AboutUs' 
                  ? 'bg-blue-900 text-white' 
                  : 'text-blue-100 hover:bg-blue-600'
              }`}
            >
              About Us
            </Link>
          </div>
        </div>
      </div>
    </nav>
  );
}

// New Footer component
function Footer() {
  return (
    <footer className="bg-gray-800 text-gray-300 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="md:flex md:justify-between">
          <div className="mb-6 md:mb-0">
            <h2 className="text-lg font-bold">CCP Little Alchemy 2</h2>
            <p className="mt-2 text-sm text-gray-400">
              Find the perfect crafting recipe for any item
            </p>
          </div>
          
          <div className="grid grid-cols-2 gap-8 sm:grid-cols-3">
            <div className="flex flex-col">
              <h3 className="text-sm font-semibold uppercase tracking-wider">Navigation</h3>
              <div className="flex flex-row gap-2">
              <p><Link to="/" className="text-sm hover:text-white">Home</Link></p>
                <p><Link to="/Wiki" className="text-sm hover:text-white">Wiki</Link></p>
                <p><Link to="/Visualizer" className="text-sm hover:text-white">Visualizer</Link></p>
                <p><Link to="/AboutUs" className="text-sm hover:text-white">About Us</Link></p>
              </div>
                
            </div>
          </div>
        </div>
        
        <div className="mt-8 border-t border-gray-700 pt-8 md:flex md:items-center md:justify-between">
          <p className="text-sm text-gray-400">
            &copy; {new Date().getFullYear()} CCP. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  );
}

interface LayoutProps {
  children: ReactNode;
}

export default function Layout({ children }: LayoutProps) {
  return (
    <div className="flex flex-col min-h-screen">
      <Navbar />
      <main className="flex-grow bg-gray-50">
        {children}
      </main>
      <Footer />
    </div>
  );
}