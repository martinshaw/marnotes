import * as React from "react";
import { Card, Tag, Icon } from "@blueprintjs/core";

export default function StatusCard() {
  return (
    <Card elevation={1} style={{ marginBottom: "15px" }}>
      <div
        className="bp5-text-small bp5-text-muted"
        style={{ marginBottom: "10px" }}
      >
        SERVER STATUS
      </div>
      <Tag intent="success" large icon="tick-circle">
        Online & Running
      </Tag>
    </Card>
  );
}
