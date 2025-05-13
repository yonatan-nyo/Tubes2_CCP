import { Link } from "react-router";
import { useState, useEffect } from "react";
import { motion, useAnimationControls } from "framer-motion";

const FallingStar = ({ delay = 0 }) => {
  const controls = useAnimationControls();

  const starColors = [
    "rgba(79, 70, 229, 1)",
    "rgba(219, 39, 119, 1)",
    "rgba(147, 51, 234, 1)",
    "rgba(59, 130, 246, 1)",
    "rgba(16, 185, 129, 1)",
  ];

  const randomColor = starColors[Math.floor(Math.random() * starColors.length)];
  const glowColor = randomColor.replace("1)", "0.8)");

  useEffect(() => {
    const startAnimation = async () => {
      await controls.start({
        x: Math.random() * 200 - 100,
        y: window.innerHeight + 100,
        opacity: [1, 0.8, 0],
        transition: {
          duration: 2 + Math.random() * 3,
          ease: "easeIn",
          delay,
        },
      });

      controls.set({
        x: Math.random() * window.innerWidth,
        y: -20,
        opacity: 1,
      });

      startAnimation();
    };

    startAnimation();
  }, [controls, delay]);

  return (
    <motion.div
      animate={controls}
      initial={{
        x: Math.random() * window.innerWidth,
        y: -20,
        opacity: 1,
      }}
      className="absolute pointer-events-none z-10">
      <div className="relative">
        <div
          className="w-3 h-3 rounded-full"
          style={{
            background: randomColor,
            boxShadow: `0 0 15px 5px ${glowColor}, 0 0 30px 8px rgba(255,255,255,0.3)`,
          }}></div>

        <div
          className="absolute top-0 left-1/2 w-[2px] h-20 -z-10 transform -translate-x-1/2 origin-top"
          style={{ background: `linear-gradient(to bottom, ${randomColor}, transparent)` }}></div>
      </div>
    </motion.div>
  );
};

export default function Home() {
  const [isLoaded, setIsLoaded] = useState(false);

  useEffect(() => {
    setIsLoaded(true);
  }, []);

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-blue-50 overflow-hidden">
      <div className="fixed inset-0 overflow-hidden z-0 opacity-70 pointer-events-none">
        <div className="absolute top-0 left-0 w-full h-full bg-grid-pattern"></div>
        <div className="absolute top-20 right-10 w-64 h-64 rounded-full bg-gradient-to-br from-pink-400 to-purple-500 opacity-20 blur-3xl"></div>
        <div className="absolute bottom-10 left-20 w-80 h-80 rounded-full bg-gradient-to-tr from-blue-400 to-teal-300 opacity-20 blur-3xl"></div>

        {Array(15)
          .fill(0)
          .map((_, i) => (
            <FallingStar key={`star-${i}`} delay={i * 0.3} />
          ))}
      </div>

      <div className="relative overflow-hidden">
        <div className="absolute inset-0 bg-pattern opacity-10"></div>
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: isLoaded ? 1 : 0 }}
          transition={{ duration: 0.8 }}
          className="max-w-7xl mx-auto py-16 px-4 sm:px-6 lg:py-24 lg:px-8 text-center relative">
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-80 h-80 rounded-full bg-gradient-to-br from-blue-400/20 to-purple-500/20 blur-3xl -z-10"></div>

          <motion.h1
            initial={{ y: -30, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ delay: 0.2, duration: 0.7, type: "spring" }}
            className="text-5xl sm:text-6xl md:text-7xl font-extrabold text-gray-900 mb-6">
            <span className="bg-clip-text text-transparent bg-gradient-to-r from-blue-600 via-purple-500 to-blue-600 animate-text-shine">
              CCP Little Alchemy 2
            </span>
          </motion.h1>

          <motion.p
            initial={{ y: -20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ delay: 0.5, duration: 0.7 }}
            className="text-xl md:text-2xl text-gray-700 max-w-3xl mx-auto mb-12">
            Discover the optimal crafting paths and combine elements to create amazing new items!
          </motion.p>

          <motion.div
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ delay: 0.7, duration: 0.5 }}
            className="flex flex-wrap justify-center gap-6 mb-16">
            <Link
              to="/Visualizer"
              className="group relative inline-flex items-center justify-center px-8 py-4 font-bold text-white transition-all duration-300 ease-in-out bg-gradient-to-r from-blue-600 to-indigo-600 rounded-xl hover:from-indigo-600 hover:to-blue-600 shadow-lg hover:shadow-blue-500/40">
              <span>Start Crafting</span>
              <div className="absolute -top-1 -right-1 w-3 h-3 bg-yellow-300 rounded-full animate-ping opacity-75"></div>
            </Link>
            <Link
              to="/Wiki"
              className="group relative inline-flex items-center px-8 py-4 font-bold overflow-hidden rounded-xl bg-white border-2 border-blue-500 text-blue-600 shadow hover:shadow-lg transition-all duration-300 hover:bg-blue-50">
              <span className="absolute inset-0 translate-y-full bg-blue-100 transition-transform duration-300 ease-out group-hover:translate-y-0 -z-10"></span>
              Browse Elements
            </Link>
          </motion.div>

          <div className="relative h-64 w-full max-w-3xl mx-auto mb-8">
            {[1, 2, 3, 4, 5, 6, 7].map((i) => (
              <motion.div
                key={i}
                initial={{ x: Math.random() * 100 - 50, y: Math.random() * 100 - 50, opacity: 0, rotate: Math.random() * 180 }}
                animate={{
                  x: [Math.random() * 300 - 150, Math.random() * 300 - 150],
                  y: [Math.random() * 150, Math.random() * 150 - 75],
                  opacity: [0.6, 0.9, 0.6],
                  rotate: [Math.random() * 90, Math.random() * -90],
                }}
                transition={{
                  repeat: Infinity,
                  repeatType: "reverse",
                  duration: 8 + i * 2,
                  ease: "easeInOut",
                }}
                className={`absolute w-12 h-12 rounded-lg shadow-lg pointer-events-none backdrop-blur-sm ${
                  i % 3 === 0
                    ? "bg-gradient-to-br from-blue-400 to-purple-500"
                    : i % 3 === 1
                    ? "bg-gradient-to-tr from-green-400 to-blue-500"
                    : "bg-gradient-to-br from-amber-400 to-red-500"
                }`}
                style={{
                  left: `${15 + (i % 7) * 12}%`,
                  top: `${10 + (i % 4) * 20}%`,
                  zIndex: 1,
                  boxShadow: "0 8px 32px rgba(31, 38, 135, 0.15)",
                }}
              />
            ))}
          </div>
        </motion.div>
      </div>

      <div className="bg-gradient-to-br from-slate-900 to-blue-900 py-20 relative">
        <div className="absolute inset-0 overflow-hidden opacity-20">
          <div className="absolute top-1/3 right-1/4 w-96 h-96 rounded-full bg-blue-500 mix-blend-screen blur-3xl"></div>
          <div className="absolute bottom-1/3 left-1/4 w-64 h-64 rounded-full bg-purple-500 mix-blend-screen blur-3xl"></div>
        </div>

        <motion.div
          initial={{ opacity: 0 }}
          whileInView={{ opacity: 1 }}
          transition={{ duration: 1 }}
          viewport={{ once: true, margin: "-100px" }}
          className="max-w-7xl mx-auto px-4 sm:px-6 text-center relative z-10">
          <h2 className="text-3xl sm:text-4xl font-bold text-white mb-16">Search Algorithms</h2>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {[
              {
                name: "Breadth-First Search",
                description:
                  "Finds the shortest recipe paths between elements by exploring all possible combinations level by level, ideal for discovering minimal crafting steps",
                icon: (
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="h-10 w-10"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor">
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={1.5}
                      d="M9 17V7m0 10a2 2 0 01-2 2H5a2 2 0 01-2-2V7a2 2 0 012-2h2a2 2 0 012 2m0 10a2 2 0 002 2h2a2 2 0 002-2M9 7a2 2 0 012-2h2a2 2 0 012 2m0 10V7m0 10a2 2 0 002 2h2a2 2 0 002-2V7a2 2 0 00-2-2h-2a2 2 0 00-2 2"
                    />
                  </svg>
                ),
                color: "from-blue-500 to-cyan-400",
              },
              {
                name: "Depth-First Search",
                description:
                  "Explores crafting paths deeply before backtracking, allowing for efficient discovery of multiple recipe variations for any target element",
                icon: (
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="h-10 w-10"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M19 14l-7 7m0 0l-7-7m7 7V3" />
                  </svg>
                ),
                color: "from-purple-500 to-pink-400",
              },
              {
                name: "Bidirectional Search",
                description:
                  "Searches simultaneously from both basic elements and target element, dramatically speeding up discovery of complex recipes",
                icon: (
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="h-10 w-10"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor">
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={1.5}
                      d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4"
                    />
                  </svg>
                ),
                color: "from-amber-500 to-orange-400",
              },
            ].map((algo, index) => (
              <motion.div
                key={index}
                initial={{ y: 50, opacity: 0 }}
                whileInView={{ y: 0, opacity: 1 }}
                viewport={{ once: true }}
                transition={{ duration: 0.5, delay: index * 0.2 }}
                className="bg-white/10 backdrop-blur-sm rounded-xl p-6 text-white border border-white/20 hover:bg-white/20 transition-colors">
                <div
                  className={`mx-auto w-16 h-16 mb-4 rounded-full flex items-center justify-center bg-gradient-to-r ${algo.color}`}>
                  {algo.icon}
                </div>
                <h3 className="text-xl font-bold mb-3">{algo.name}</h3>
                <p className="text-blue-100 text-sm">{algo.description}</p>

                <div className="mt-6 pt-4 border-t border-white/10">
                  <div className="flex justify-center space-x-1">
                    {[...Array(5)].map((_, i) => (
                      <div
                        key={i}
                        className="w-1.5 h-6 bg-white/30 rounded-full"
                        style={{
                          animation: `equalizer 1.5s ${i * 0.1}s ease-in-out infinite alternate`,
                          height: `${1 + Math.random() * 1.5}rem`,
                        }}></div>
                    ))}
                  </div>
                </div>
              </motion.div>
            ))}
          </div>

          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: 0.5 }}
            className="mt-10 pt-6">
            <div className="inline-block px-6 py-3 rounded-lg bg-white/10 text-white text-sm border border-white/20">
              Toggle between algorithms to find <span className="font-bold text-yellow-300">single recipe</span> or{" "}
              <span className="font-bold text-yellow-300">multiple recipes</span> with real-time visualization
            </div>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: 0.7 }}
            className="mt-4">
            <div className="inline-block px-6 py-3 rounded-lg bg-white/10 text-white text-sm border border-white/20">
              <span className="font-bold text-green-300">Multithreaded optimization</span> for faster discovery of multiple recipe
              variations
            </div>
          </motion.div>
        </motion.div>

        <style>{`
          @keyframes equalizer {
            0% { height: 0.5rem; }
            100% { height: 2rem; }
          }
        `}</style>
      </div>

      <div className="max-w-7xl mx-auto py-16 px-4 sm:px-6 relative">
        <motion.div
          initial={{ opacity: 0 }}
          whileInView={{ opacity: 1 }}
          transition={{ duration: 1 }}
          viewport={{ once: true, margin: "-100px" }}
          className="text-center mb-16">
          <h2 className="text-3xl sm:text-4xl font-bold text-gray-900">Explore Our Features</h2>
          <div className="w-32 h-1.5 bg-gradient-to-r from-blue-400 to-purple-600 mx-auto mt-5 rounded-full"></div>
        </motion.div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-10 mb-20">
          <motion.div
            initial={{ x: -50, opacity: 0 }}
            whileInView={{ x: 0, opacity: 1 }}
            viewport={{ once: true, margin: "-100px" }}
            transition={{ duration: 0.6 }}
            whileHover={{ y: -8, boxShadow: "0 25px 50px -12px rgba(59, 130, 246, 0.25)" }}
            className="bg-white p-8 rounded-2xl shadow-lg transition-all duration-300 border-t-4 border-blue-600 hover:border-blue-500 group">
            <div className="bg-blue-100 rounded-full w-16 h-16 flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                className="h-8 w-8 text-blue-600"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </div>
            <h2 className="text-2xl font-bold text-gray-900 mb-4 group-hover:text-blue-600 transition-colors duration-300">
              Recipe Visualizer
            </h2>
            <p className="text-gray-600 mb-8 leading-relaxed">
              Find the most efficient crafting paths using our advanced algorithms. Visualize recipe trees and discover new
              combinations!
            </p>
            <Link to="/Visualizer" className="inline-flex items-center text-blue-600 font-medium hover:text-blue-800 group">
              <span className="border-b-2 border-transparent group-hover:border-blue-600 transition-all duration-300">
                Try Visualizer
              </span>
              <svg
                xmlns="http://www.w3.org/2000/svg"
                className="h-5 w-5 ml-2 transform transition-transform group-hover:translate-x-2"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" />
              </svg>
            </Link>
          </motion.div>

          <motion.div
            initial={{ x: 50, opacity: 0 }}
            whileInView={{ x: 0, opacity: 1 }}
            viewport={{ once: true, margin: "-100px" }}
            transition={{ duration: 0.6 }}
            whileHover={{ y: -8, boxShadow: "0 25px 50px -12px rgba(16, 185, 129, 0.25)" }}
            className="bg-white p-8 rounded-2xl shadow-lg transition-all duration-300 border-t-4 border-green-600 hover:border-green-500 group">
            <div className="bg-green-100 rounded-full w-16 h-16 flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                className="h-8 w-8 text-green-600"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
                />
              </svg>
            </div>
            <h2 className="text-2xl font-bold text-gray-900 mb-4 group-hover:text-green-600 transition-colors duration-300">
              Elements Encyclopedia
            </h2>
            <p className="text-gray-600 mb-8 leading-relaxed">
              Access our comprehensive database of elements, recipes, and combinations. Learn the secrets of crafting any item in
              Little Alchemy 2.
            </p>
            <Link to="/Wiki" className="inline-flex items-center text-green-600 font-medium hover:text-green-800 group">
              <span className="border-b-2 border-transparent group-hover:border-green-600 transition-all duration-300">
                Browse Encyclopedia
              </span>
              <svg
                xmlns="http://www.w3.org/2000/svg"
                className="h-5 w-5 ml-2 transform transition-transform group-hover:translate-x-2"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" />
              </svg>
            </Link>
          </motion.div>
        </div>
      </div>

      <div className="bg-gradient-to-r from-blue-900 to-indigo-900 text-white py-20 relative overflow-hidden">
        <div className="absolute inset-0 overflow-hidden opacity-20">
          <div className="absolute top-0 left-0 -translate-x-1/2 w-96 h-96 rounded-full bg-blue-400 mix-blend-multiply blur-3xl"></div>
          <div className="absolute bottom-0 right-0 translate-x-1/3 w-96 h-96 rounded-full bg-indigo-400 mix-blend-multiply blur-3xl"></div>
        </div>

        <motion.div
          initial={{ y: 50, opacity: 0 }}
          whileInView={{ y: 0, opacity: 1 }}
          viewport={{ once: true, margin: "-100px" }}
          transition={{ duration: 0.8 }}
          className="max-w-7xl mx-auto px-4 sm:px-6 text-center relative z-10">
          <h2 className="text-3xl sm:text-4xl font-bold mb-16">Discover the World of Alchemy</h2>

          <div className="grid grid-cols-2 md:grid-cols-4 gap-8">
            {[
              { value: "700+", label: "Total Elements", color: "from-blue-400 to-blue-200" },
              { value: "2500+", label: "Combinations", color: "from-purple-400 to-purple-200" },
              { value: "3", label: "Search Algorithms", color: "from-green-400 to-green-200" },
              { value: "∞", label: "Possibilities", color: "from-pink-400 to-pink-200" },
            ].map((stat, index) => (
              <motion.div
                key={index}
                initial={{ y: 20, opacity: 0 }}
                whileInView={{ y: 0, opacity: 1 }}
                viewport={{ once: true }}
                transition={{ duration: 0.5, delay: index * 0.1 }}
                className="p-4">
                <motion.div
                  initial={{ scale: 0 }}
                  whileInView={{ scale: 1 }}
                  viewport={{ once: true }}
                  transition={{
                    type: "spring",
                    stiffness: 260,
                    damping: 20,
                    delay: 0.1 + index * 0.1,
                  }}
                  className={`text-5xl font-bold bg-gradient-to-b ${stat.color} bg-clip-text text-transparent mb-3`}>
                  {stat.value}
                </motion.div>
                <div className="text-sm uppercase tracking-wider text-blue-100">{stat.label}</div>
              </motion.div>
            ))}
          </div>
        </motion.div>
      </div>

      <div className="max-w-7xl mx-auto py-20 px-4 sm:px-6">
        <motion.div
          initial={{ y: 50, opacity: 0 }}
          whileInView={{ y: 0, opacity: 1 }}
          viewport={{ once: true, margin: "-100px" }}
          transition={{ duration: 0.8 }}
          className="bg-gradient-to-r from-purple-600 via-indigo-600 to-blue-600 rounded-3xl shadow-2xl overflow-hidden relative">
          <div className="absolute top-0 right-0 w-32 h-32 bg-white opacity-10 rounded-full -translate-y-1/2 translate-x-1/2"></div>
          <div className="absolute bottom-0 left-0 w-24 h-24 bg-white opacity-10 rounded-full translate-y-1/2 -translate-x-1/2"></div>

          <div className="px-6 py-16 md:py-20 md:px-16 text-center text-white relative z-10">
            <motion.h2
              initial={{ y: -20, opacity: 0 }}
              whileInView={{ y: 0, opacity: 1 }}
              viewport={{ once: true }}
              transition={{ duration: 0.5 }}
              className="text-3xl sm:text-4xl font-extrabold mb-5 drop-shadow-md">
              Ready to start your alchemy journey?
            </motion.h2>
            <motion.p
              initial={{ y: -10, opacity: 0 }}
              whileInView={{ y: 0, opacity: 1 }}
              viewport={{ once: true }}
              transition={{ duration: 0.5, delay: 0.1 }}
              className="text-lg sm:text-xl opacity-90 mb-10 max-w-2xl mx-auto">
              Begin combining elements, discover new creations, and unlock the mysteries of Little Alchemy 2!
            </motion.p>
            <motion.div
              initial={{ y: 10, opacity: 0 }}
              whileInView={{ y: 0, opacity: 1 }}
              viewport={{ once: true }}
              transition={{ duration: 0.5, delay: 0.2 }}>
              <Link
                to="/Visualizer"
                className="inline-block bg-white text-indigo-600 px-8 py-4 rounded-xl font-bold shadow-lg hover:shadow-white/20 transform transition-all duration-300 hover:-translate-y-1 hover:scale-105">
                Start Now
              </Link>
            </motion.div>
          </div>
        </motion.div>
      </div>
    </div>
  );
}
