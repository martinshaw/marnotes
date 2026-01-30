import * as React from "react";
import {
  Navbar,
  NavbarGroup,
  NavbarHeading,
  Button,
  Alignment,
} from "@blueprintjs/core";

export default function NavBar() {
  const buttonStyling = {
    marginLeft: 2,
    marginRight: 2,
    marginTop: 5,
    color: "rgb(255, 255, 255)",
    height: 22,
    minHeight: "auto",
  };

  return (
    <Navbar
      fixedToTop
      style={{
        marginBottom: 0,
        height: 34,
        background: "#000",
        boxShadow: "none",
        padding: "0 4px",
        borderBottom: "1px solid #222",
        userSelect: "none",
      }}
    >
      <NavbarGroup align={Alignment.START} style={{ height: 28 }}>
        <Button
          variant="minimal"
          text={<span style={{ fontSize: 12, color: "#fff" }}>Documents</span>}
          style={buttonStyling}
        />
        <Button
          variant="minimal"
          text={<span style={{ fontSize: 12, color: "#fff" }}>Endpoints</span>}
          style={buttonStyling}
        />
        <Button
          variant="minimal"
          text={<span style={{ fontSize: 12, color: "#fff" }}>Settings</span>}
          style={buttonStyling}
        />
      </NavbarGroup>
      <NavbarGroup align={Alignment.END} style={{ height: 28 }}>
        <NavbarHeading
          style={{
            height: 32,
            minHeight: "auto",
            color: "#fff",
            paddingLeft: 10,
            paddingRight: 10,
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            marginTop: 6,
            userSelect: "text",
          }}
        >
          MarNotes
        </NavbarHeading>
        <Button
          variant="minimal"
          text={<span style={{ fontSize: 12, color: "#fff" }}>Help</span>}
          style={buttonStyling}
        />
      </NavbarGroup>
    </Navbar>
  );
}
