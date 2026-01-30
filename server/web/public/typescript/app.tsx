import * as React from "react";
import * as ReactDOM from "react-dom/client";
import Header from "./components/Header";
import StatusCard from "./components/StatusCard";
import PortCard from "./components/PortCard";
import EndpointsCard from "./components/EndpointsCard";
import DocumentsData from "./components/DocumentsData";
import Footer from "./components/Footer";
import NavBar from "./components/NavBar";

function App() {
  const [port, setPort] = React.useState<string>("Loading...");

  React.useEffect(() => {
    const configuredPort = (window as any).__JSON_SERVER_PORT__ || "";
    const fallbackPort = window.location.port || "80";
    setPort(configuredPort || fallbackPort);
  }, []);

  return (
    <>
      <NavBar />
      <div className="container">
        <Header />
        <StatusCard />
        <PortCard port={port} />
        <DocumentsData jsonPort={port} />
        <EndpointsCard />
        <Footer />
      </div>
    </>
  );
}

const root = ReactDOM.createRoot(
  document.getElementById("root") as HTMLElement,
);
root.render(<App />);
