import * as React from "react";
import { Text } from "@blueprintjs/core";

export default function Footer() {
  return (
    <div style={{ marginTop: "30px", textAlign: "center" }}>
      <Text className="bp5-text-muted bp5-text-small">
        MarNotes Server v1.0 • React Dashboard
      </Text>
    </div>
  );
}
