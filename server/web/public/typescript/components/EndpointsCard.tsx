import * as React from "react";
import { Card, Code } from "@blueprintjs/core";

export default function EndpointsCard() {
  return (
    <Card elevation={1} style={{ marginBottom: "15px" }}>
      <div
        className="bp5-text-small bp5-text-muted"
        style={{ marginBottom: "12px" }}
      >
        API ENDPOINTS
      </div>
      <div style={{ lineHeight: "2" }}>
        <div>
          <Code>GET /documents</Code> - List all documents
        </div>
        <div>
          <Code>GET /doc/&#123;name&#125;</Code> - Get specific document
        </div>
        <div>
          <Code>GET /health</Code> - Health check
        </div>
      </div>
    </Card>
  );
}
