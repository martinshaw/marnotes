import * as React from "react";
import {
  Navbar,
  NavbarGroup,
  NavbarHeading,
  Button,
  Alignment,
  Menu,
  MenuItem,
  Popover,
  Spinner,
} from "@blueprintjs/core";

interface NavBarProps {
  port: string;
  onDocumentSelect: (document: any) => void;
}

export default function NavBar({ port, onDocumentSelect }: NavBarProps) {
  const [documents, setDocuments] = React.useState<string[]>([]);
  const [isLoading, setIsLoading] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);

  const buttonStyling: React.CSSProperties = {
    marginLeft: 2,
    marginRight: 2,
    marginTop: 5,
    color: "rgb(255, 255, 255)",
    height: 22,
    minHeight: "auto",
    userSelect: "none",
  };

  const titleStyling: React.CSSProperties = {
    flex: 1,
    height: 32,
    minHeight: "auto",
    color: "#fff",
    paddingLeft: 10,
    paddingRight: 10,
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    marginTop: 6,
    userSelect: "none",
  };

  // Fetch documents list
  React.useEffect(() => {
    if (port === "Loading...") return;

    const fetchDocuments = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const response = await fetch(`http://localhost:${port}/documents`);
        if (!response.ok) {
          throw new Error(`Failed to fetch documents: ${response.statusText}`);
        }
        const data = await response.json();
        setDocuments(data.documents || []);
      } catch (err) {
        console.error("Error fetching documents:", err);
        setError(err instanceof Error ? err.message : "Unknown error");
      } finally {
        setIsLoading(false);
      }
    };

    fetchDocuments();
  }, [port]);

  const handleDocumentClick = async (filename: string) => {
    try {
      const response = await fetch(
        `http://localhost:${port}/documents/${filename}`,
      );
      if (!response.ok) {
        throw new Error(`Failed to fetch document: ${response.statusText}`);
      }
      const documentData = await response.json();
      onDocumentSelect(documentData);
    } catch (err) {
      console.error("Error fetching document:", err);
    }
  };

  const documentsMenu = (
    <Menu>
      {isLoading && (
        <MenuItem text="Loading..." icon={<Spinner size={16} />} disabled />
      )}
      {error && <MenuItem text={`Error: ${error}`} disabled />}
      {!isLoading && !error && documents.length === 0 && (
        <MenuItem text="No documents available" disabled />
      )}
      {!isLoading &&
        !error &&
        documents.map((doc) => (
          <MenuItem
            key={doc}
            text={doc.replace(".json", "")}
            onClick={() => handleDocumentClick(doc)}
          />
        ))}
    </Menu>
  );

  return (
    <Navbar
      fixedToTop
      style={{
        display: "flex",
        flexDirection: "row",
        justifyContent: "space-between",
        marginBottom: 0,
        height: 34,
        background: "#000",
        boxShadow: "none",
        padding: "0 4px",
        borderBottom: "1px solid #222",
        userSelect: "none",
      }}
      tabIndex={-1}
    >
      <NavbarGroup style={{ height: 28 }} tabIndex={-1}>
        <Popover content={documentsMenu} placement="bottom-start">
          <Button
            variant="minimal"
            text={
              <span style={{ fontSize: 12, color: "#fff" }}>Documents</span>
            }
            style={buttonStyling}
            tabIndex={-1}
          />
        </Popover>
        <Button
          variant="minimal"
          text={<span style={{ fontSize: 12, color: "#fff" }}>Endpoints</span>}
          style={buttonStyling}
          tabIndex={-1}
        />
        <Button
          variant="minimal"
          text={<span style={{ fontSize: 12, color: "#fff" }}>Settings</span>}
          style={buttonStyling}
          tabIndex={-1}
        />
      </NavbarGroup>
      <NavbarGroup
        style={{ height: 28, flex: 1, display: "flex" }}
        tabIndex={-1}
      >
        <NavbarHeading style={titleStyling} tabIndex={-1}>
          Marnotes
        </NavbarHeading>
        <Button
          variant="minimal"
          text={<span style={{ fontSize: 12, color: "#fff" }}>Help</span>}
          style={buttonStyling}
          tabIndex={-1}
        />
      </NavbarGroup>
    </Navbar>
  );
}
