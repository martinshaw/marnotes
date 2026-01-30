type DocumentsDataProps = {
  jsonPort: string;
};

type DocumentList = {
  documents: string[];
  count: number;
};

export default function DocumentsData({ jsonPort }: DocumentsDataProps) {
  const [documents, setDocuments] = React.useState<string[]>([]);
  const [selectedDoc, setSelectedDoc] = React.useState<string | null>(null);
  const [docContent, setDocContent] = React.useState<any>(null);
  const [loading, setLoading] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);

  const fetchDocuments = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`http://localhost:${jsonPort}/documents`);
      if (!response.ok)
        throw new Error(`HTTP error! status: ${response.status}`);
      const data: DocumentList = await response.json();
      setDocuments(data.documents || []);
      setDocContent(null);
      setSelectedDoc(null);
    } catch (err) {
      setError(`Failed to fetch documents: ${err}`);
    } finally {
      setLoading(false);
    }
  };

  const fetchDocument = async (docName: string) => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(
        `http://localhost:${jsonPort}/doc/${docName}`,
      );
      if (!response.ok)
        throw new Error(`HTTP error! status: ${response.status}`);
      const data = await response.json();
      setSelectedDoc(docName);
      setDocContent(data);
    } catch (err) {
      setError(`Failed to fetch document: ${err}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="info-card">
      <div className="info-label">Document Browser</div>
      <div style={{ marginTop: "12px" }}>
        <button
          onClick={fetchDocuments}
          disabled={loading}
          style={{
            padding: "8px 16px",
            background: "#667eea",
            color: "white",
            border: "none",
            borderRadius: "4px",
            cursor: "pointer",
            marginBottom: "12px",
            fontSize: "14px",
            fontWeight: "600",
          }}
        >
          {loading ? "Loading..." : "List Documents"}
        </button>

        {error && (
          <div
            style={{
              color: "#dc2626",
              padding: "8px",
              marginBottom: "12px",
              background: "#fee2e2",
              borderRadius: "4px",
              fontSize: "14px",
            }}
          >
            {error}
          </div>
        )}

        {documents.length > 0 && (
          <div style={{ marginBottom: "12px" }}>
            <div
              style={{ marginBottom: "8px", fontSize: "13px", color: "#555" }}
            >
              Found {documents.length} document
              {documents.length !== 1 ? "s" : ""}:
            </div>
            <div style={{ display: "flex", flexWrap: "wrap", gap: "8px" }}>
              {documents.map((doc) => (
                <button
                  key={doc}
                  onClick={() => fetchDocument(doc.replace(".json", ""))}
                  style={{
                    padding: "6px 12px",
                    background:
                      selectedDoc === doc.replace(".json", "")
                        ? "#10b981"
                        : "#e5e7eb",
                    color:
                      selectedDoc === doc.replace(".json", "")
                        ? "white"
                        : "#333",
                    border: "none",
                    borderRadius: "4px",
                    cursor: "pointer",
                    fontSize: "13px",
                    fontWeight: "500",
                  }}
                >
                  {doc}
                </button>
              ))}
            </div>
          </div>
        )}

        {docContent && (
          <div
            style={{
              marginTop: "12px",
              padding: "12px",
              background: "#f9fafb",
              borderRadius: "4px",
              border: "1px solid #e5e7eb",
            }}
          >
            <div
              style={{ fontSize: "12px", color: "#666", marginBottom: "8px" }}
            >
              Content of <strong>{selectedDoc}.json</strong>:
            </div>
            <pre
              style={{
                margin: 0,
                padding: "8px",
                background: "#fff",
                borderRadius: "4px",
                overflow: "auto",
                fontSize: "12px",
                color: "#333",
                fontFamily: "Courier New, monospace",
                border: "1px solid #e5e7eb",
              }}
            >
              {JSON.stringify(docContent, null, 2)}
            </pre>
          </div>
        )}
      </div>
    </div>
  );
}
