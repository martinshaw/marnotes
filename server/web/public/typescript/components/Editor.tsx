/*
All Rights Reserved, (c) 2026 marnotes

Author:      Martin Shaw (developer@martinshaw.co)
Created:     2026-01-31T03:30:49.947Z
Modified:    2026-01-31T03:30:49.947Z
File Name:   Editor.tsx

description

*/

import * as React from "react";
import {
  $getRoot,
  $getSelection,
  EditorThemeClasses,
  InitialEditorConfig,
  InitialEditorStateType,
} from "lexical";
import { LexicalComposer } from "@lexical/react/LexicalComposer";
import { PlainTextPlugin } from "@lexical/react/LexicalPlainTextPlugin";
import { ContentEditable } from "@lexical/react/LexicalContentEditable";
import { HistoryPlugin } from "@lexical/react/LexicalHistoryPlugin";
import { LexicalErrorBoundary } from "@lexical/react/LexicalErrorBoundary";
import { RichTextPlugin } from "@lexical/react/LexicalRichTextPlugin";
import { ListPlugin } from "@lexical/react/LexicalListPlugin";
import { ListNode, ListItemNode } from "@lexical/list";
import { CheckListPlugin } from "@lexical/react/LexicalCheckListPlugin";
import { EditorRefPlugin } from "@lexical/react/LexicalEditorRefPlugin";

export type EditorProps = {};

const Editor: React.FC<EditorProps> = (props: EditorProps) => {
  const initialConfig = {
    namespace: "MarnotesEditor",
    onError: (error) => console.error(error),
    nodes: [ListNode, ListItemNode],
  };

  const editorRef = React.useRef(null);

  return (
    <div style={{ width: "100%", height: "100%", padding: "7px 14px" }}>
      <LexicalComposer initialConfig={initialConfig}>
        <RichTextPlugin
          contentEditable={<ContentEditable />}
          ErrorBoundary={LexicalErrorBoundary}
        />
        <HistoryPlugin />
        <ListPlugin />
        <CheckListPlugin />
        <EditorRefPlugin editorRef={editorRef} />
      </LexicalComposer>
    </div>
  );
};

export default Editor;
