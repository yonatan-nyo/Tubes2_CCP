import { lazy } from "react";
import { createBrowserRouter, Outlet } from "react-router";

const router = createBrowserRouter([
  {
    path: "/",
    element: <Outlet />,
    children: [
      {
        path: "/",
        Component: lazy(() => import("../pages/Home")),
      },
    ],
  },
]);

export default router;
