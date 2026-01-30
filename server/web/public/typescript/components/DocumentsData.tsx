import * as React from "react";
import { Card, Button, ButtonGroup, Callout, Pre } from "@blueprintjs/core";

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
  const [publicKey, setPublicKey] = React.useState<string | null>(null);

  React.useEffect(() => {
    fetchPublicKey();
  }, []);

  const fetchPublicKey = async () => {
    try {
      const response = await fetch(`http://localhost:${jsonPort}/publickey`);
      if (response.ok) {
        const keyText = await response.text();
        setPublicKey(keyText);
      }
    } catch (err) {
      console.log("Encryption not available, using unencrypted mode");
    }
  };

  const decryptResponse = async (encryptedData: any) => {
    if (typeof encryptedData === "object" && encryptedData.encrypted) {
      const workerCode = `
        self.onmessage = async (e) => {
          const { encrypted, publicKeyPem } = e.data;
          try {
            const publicKeyStr = publicKeyPem
              .replace('-----BEGIN PUBLIC KEY-----', '')
              .replace('-----END PUBLIC KEY-----', '')
              .replace(/\\n/g, '');
            
            const binaryString = atob(encrypted);
            const bytes = new Uint8Array(binaryString.length);
            for (let i = 0; i < binaryString.length; i++) {
              bytes[i] = binaryString.charCodeAt(i);
            }
            
            self.postMessage({ error: 'RSA decryption requires crypto module' });
          } catch (err) {
            self.postMessage({ error: err.message });
          }
        };
      `;

      return encryptedData;
    }
    return encryptedData;
  };

  const fetchDocuments = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`http://localhost:${jsonPort}/documents`);
      if (!response.ok)
        throw new Error(`HTTP error! status: ${response.status}`);
      const data = await response.json();
      const decrypted = await decryptResponse(data);
      setDocuments(decrypted.documents || []);
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
      const decrypted = await decryptResponse(data);
      setSelectedDoc(docName);
      setDocContent(decrypted);
    } catch (err) {
      setError(`Failed to fetch document: ${err}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card elevation={1} style={{ marginBottom: "15px" }}>
      <div
        className="bp5-text-small bp5-text-muted"
        style={{ marginBottom: "12px" }}
      >
        DOCUMENT BROWSER
      </div>
      <div>
        <Button
          onClick={fetchDocuments}
          disabled={loading}
          intent="primary"
          icon="document"
          loading={loading}
          style={{ marginBottom: "12px" }}
        >
          List Documents
        </Button>

        {error && (
          <Callout intent="danger" style={{ marginBottom: "12px" }}>
            {error}
          </Callout>
        )}

        {documents.length > 0 && (
          <div style={{ marginBottom: "12px" }}>
            <div
              className="bp5-text-small bp5-text-muted"
              style={{ marginBottom: "8px" }}
            >
              Found {documents.length} document
              {documents.length !== 1 ? "s" : ""}:
            </div>
            <ButtonGroup minimal>
              {documents.map((doc) => (
                <Button
                  key={doc}
                  onClick={() => fetchDocument(doc.replace(".json", ""))}
                  active={selectedDoc === doc.replace(".json", "")}
                  intent={
                    selectedDoc === doc.replace(".json", "")
                      ? "success"
                      : "none"
                  }
                  small
                >
                  {doc}
                </Button>
              ))}
            </ButtonGroup>
          </div>
        )}

        {docContent && (
          <Card
            elevation={0}
            style={{ marginTop: "12px", background: "#f5f8fa" }}
          >
            <div
              className="bp5-text-small bp5-text-muted"
              style={{ marginBottom: "8px" }}
            >
              Content of <strong>{selectedDoc}.json</strong>:
            </div>
            <Pre
              style={{
                margin: 0,
                fontSize: "12px",
                maxHeight: "400px",
                overflow: "auto",
              }}
            >
              {JSON.stringify(docContent, null, 2)}
            </Pre>
          </Card>
        )}
      </div>
    </Card>
  );
}
