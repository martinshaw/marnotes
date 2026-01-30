import * as React from "react";
import { H1, Text } from "@blueprintjs/core";

export default function Header() {
  return (
    <div
      className="header"
      style={{ textAlign: "center", marginBottom: "30px" }}
    >
      <H1>📝 MarNotes</H1>
      <Text className="bp5-text-large bp5-text-muted">
        JSON Document Server Dashboard
      </Text>
    </div>
  );
}
