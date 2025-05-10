import { FaGithub, FaLinkedin } from "react-icons/fa";

export default function AboutUs() {
  const teamMembers = [
    {
      name: "Varel Tiara",
      id: "13523008",
      image: "/Varel.png",
      github: "https://github.com/varel183",
      linkedin: "https://www.linkedin.com/in/varel-tiara/",
    },
    {
      name: "Nathaniel Jonathan Rusli",
      id: "13523013",
      image: "/Nathaniel.png",
      github: "https://github.com/NathanielJR-git",
      linkedin: "https://www.linkedin.com/in/nathanieljr/",
    },
    {
      name: "Yonatan Edward Njoto",
      id: "13523036",
      image: "/Yonatan.png",
      github: "https://github.com/yonatan-nyo",
      linkedin: "https://www.linkedin.com/in/yonatan-njoto/",
    },
  ];

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      {/* Hero Section */}
      <div className="bg-gradient-to-r from-blue-600 to-indigo-700 text-white">
        <div className="container mx-auto px-6 py-20 text-center">
          <h1 className="text-4xl md:text-5xl font-bold mb-4">Kelompok 26 - CCP</h1>
          <p className="text-xl md:text-2xl mb-8">Little Alchemy 2 Recipe Finder</p>
          <div className="inline-block bg-white text-indigo-700 font-semibold px-6 py-3 rounded-md shadow-md">
            Tugas Besar 2 IF2211 Strategi Algoritma
          </div>
        </div>
      </div>

      {/* Team Members Section */}
      <div className="container mx-auto px-6 py-16">
        <h2 className="text-3xl font-bold text-center text-gray-800 mb-12">Meet Our Team</h2>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-10">
          {teamMembers.map((member, index) => (
            <div key={index} className="bg-white rounded-xl shadow-lg overflow-hidden transform transition-all hover:scale-105">
              <div className="h-48 bg-gradient-to-r from-blue-400 to-indigo-500 flex items-center justify-center">
                <img src={member.image} alt={member.name} className="h-32 w-32 rounded-full border-4 border-white object-cover" />
              </div>
              <div className="p-6">
                <div className="flex flex-row justify-between items-center mb-4">
                  <h3 className="font-bold text-xl text-gray-800">{member.name}</h3>
                  <div className="flex space-x-4">
                    <p className="text-gray-500">{member.id}</p>
                    <a href={member.github} className="text-gray-600 hover:text-gray-900">
                      <FaGithub size={20} />
                    </a>
                    <a href={member.linkedin} className="text-gray-600 hover:text-gray-900">
                      <FaLinkedin size={20} />
                    </a>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Project Details Section */}
      <div className="bg-white">
        <div className="container mx-auto px-6 py-16">
          <div className="w-full mx-auto">
            <h2 className="text-3xl font-bold text-center text-gray-800 mb-8">About the Project</h2>

            <div className="bg-gray-50 rounded-xl p-8 shadow-md">
              <h3 className="text-xl font-semibold mb-4 text-gray-800">
                Utilizing BFS and DFS Algorithms for Recipe Finding in Little Alchemy 2
              </h3>

              <div className="space-y-4 text-gray-600">
                <p>
                  This project implements BFS (Breadth-First Search) and DFS (Depth-First Search) algorithms to find recipes in
                  the popular game Little Alchemy 2. By applying these graph traversal algorithms, we can efficiently discover the
                  optimal combination paths to create new elements.
                </p>

                <div className="py-2">
                  <h4 className="font-semibold text-gray-800 mb-2">Course Information:</h4>
                  <p>Major Assignment 2 - Algorithm Strategies (IF2211)</p>
                  <p>Second Semester, Academic Year 2024/2025</p>
                </div>

                <div className="py-2">
                  <h4 className="font-semibold text-gray-800 mb-2">Main Features:</h4>
                  <ul className="list-disc pl-5 space-y-1">
                    <li>Recipe finding using BFS algorithm</li>
                    <li>Recipe finding using DFS algorithm</li>
                    <li>Interactive element encyclopedia</li>
                    <li>Visual representation of element combinations</li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
