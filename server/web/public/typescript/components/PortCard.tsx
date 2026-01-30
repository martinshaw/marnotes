import * as React from "react";
import { Card, H3 } from "@blueprintjs/core";

type PortCardProps = {
  port: string;
};

export default function PortCard({ port }: PortCardProps) {
  return (
    <Card elevation={1} style={{ marginBottom: "15px" }}>
      <div
        className="bp5-text-small bp5-text-muted"
        style={{ marginBottom: "10px" }}
      >
        JSON SERVER PORT
      </div>
      <H3 style={{ margin: 0, color: "#667eea" }}>{port}</H3>
    </Card>
  );
}
