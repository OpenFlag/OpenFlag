import Dashboard from "views/Dashboard.js";
import TableList from "views/TableList.js";
import Icons from "views/Icons.js";

import { App, Search, Eye } from "react-bootstrap-icons";

let iconSize = 25;

const dashboardRoutes = [
  {
    path: "/dashboard",
    name: "Dashboard",
    icon: <App size={iconSize} />,
    component: Dashboard,
    layout: "/admin",
  },
  {
    path: "/search",
    name: "Search",
    icon: <Search size={iconSize} />,
    component: Icons,
    layout: "/admin",
  },
  {
    path: "/evaluation",
    name: "Evaluation",
    icon: <Eye size={iconSize} />,
    component: TableList,
    layout: "/admin",
  },
];

export default dashboardRoutes;
