import * as React from "react";
import * as ReactDOM from "react-dom/client";
import NavBar from "./components/NavBar";
import Editor from "./components/Editor";

function App() {
  const [port, setPort] = React.useState<string>("Loading...");
  const [currentDocument, setCurrentDocument] = React.useState<any>(null);

  React.useEffect(() => {
    const configuredPort = (window as any).__JSON_SERVER_PORT__ || "";
    const fallbackPort = window.location.port || "80";
    setPort(configuredPort || fallbackPort);
  }, []);

  return (
    <>
      <NavBar port={port} onDocumentSelect={setCurrentDocument} />
      <Editor documentContent={currentDocument} />
    </>
  );
}

const root = ReactDOM.createRoot(
  document.getElementById("root") as HTMLElement,
);
root.render(<App />);
