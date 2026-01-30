export default function EndpointsCard() {
  return (
    <div className="info-card">
      <div className="info-label">API Endpoints</div>
      <div style={{ marginTop: "12px", lineHeight: "1.8" }}>
        <div>
          <code
            style={{
              background: "#e5e7eb",
              padding: "4px 8px",
              borderRadius: "4px",
            }}
          >
            GET /documents
          </code>{" "}
          - List all documents
        </div>
        <div>
          <code
            style={{
              background: "#e5e7eb",
              padding: "4px 8px",
              borderRadius: "4px",
            }}
          >
            GET /doc/&#123;name&#125;
          </code>{" "}
          - Get specific document
        </div>
        <div>
          <code
            style={{
              background: "#e5e7eb",
              padding: "4px 8px",
              borderRadius: "4px",
            }}
          >
            GET /health
          </code>{" "}
          - Health check
        </div>
      </div>
    </div>
  );
}
