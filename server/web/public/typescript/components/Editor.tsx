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
  $createParagraphNode,
  $createTextNode,
  EditorThemeClasses,
  InitialEditorConfig,
  InitialEditorStateType,
  LexicalEditor,
} from "lexical";
import { LexicalComposer } from "@lexical/react/LexicalComposer";
import { PlainTextPlugin } from "@lexical/react/LexicalPlainTextPlugin";
import { ContentEditable } from "@lexical/react/LexicalContentEditable";
import { HistoryPlugin } from "@lexical/react/LexicalHistoryPlugin";
import { LexicalErrorBoundary } from "@lexical/react/LexicalErrorBoundary";
import { RichTextPlugin } from "@lexical/react/LexicalRichTextPlugin";
import { ListPlugin } from "@lexical/react/LexicalListPlugin";
import { ListNode, ListItemNode } from "@lexical/list";
import { HeadingNode, QuoteNode } from "@lexical/rich-text";
import { CheckListPlugin } from "@lexical/react/LexicalCheckListPlugin";
import { EditorRefPlugin } from "@lexical/react/LexicalEditorRefPlugin";
import { useLexicalComposerContext } from "@lexical/react/LexicalComposerContext";
import { AutoFocusPlugin } from "@lexical/react/LexicalAutoFocusPlugin";

export type EditorProps = {
  documentContent?: any;
};

// Plugin to update editor content when document changes
function DocumentLoaderPlugin({ documentContent }: { documentContent: any }) {
  const [editor] = useLexicalComposerContext();

  React.useEffect(() => {
    if (!documentContent) return;

    // Check if documentContent has a Lexical editor state structure
    if (documentContent.root && documentContent.root.children) {
      // It's a Lexical editor state, parse it directly
      try {
        const editorState = editor.parseEditorState(documentContent);
        editor.setEditorState(editorState);
      } catch (error) {
        console.error("Failed to parse Lexical editor state:", error);
        // Fall back to displaying as JSON
        editor.update(() => {
          const root = $getRoot();
          root.clear();
          const formattedContent = JSON.stringify(documentContent, null, 2);
          const paragraph = $createParagraphNode();
          const textNode = $createTextNode(formattedContent);
          paragraph.append(textNode);
          root.append(paragraph);
        });
      }
    } else {
      // It's regular JSON, display it as formatted text
      editor.update(() => {
        const root = $getRoot();
        root.clear();
        const formattedContent = JSON.stringify(documentContent, null, 2);
        const paragraph = $createParagraphNode();
        const textNode = $createTextNode(formattedContent);
        paragraph.append(textNode);
        root.append(paragraph);
      });
    }
  }, [documentContent, editor]);

  return null;
}

const Editor: React.FC<EditorProps> = ({ documentContent }: EditorProps) => {
  const initialConfig = {
    namespace: "MarnotesEditor",
    onError: (error) => console.error(error),
    nodes: [HeadingNode, ListNode, ListItemNode, QuoteNode],
  };

  const editorRef = React.useRef(null);

  return (
    <div style={{ width: "100%", height: "100%", padding: "10px 14px" }}>
      <LexicalComposer initialConfig={initialConfig}>
        <AutoFocusPlugin />
        <RichTextPlugin
          contentEditable={<ContentEditable />}
          ErrorBoundary={LexicalErrorBoundary}
        />
        <HistoryPlugin />
        <ListPlugin />
        <CheckListPlugin />
        <EditorRefPlugin editorRef={editorRef} />
        <DocumentLoaderPlugin documentContent={documentContent} />
      </LexicalComposer>
    </div>
  );
};

export default Editor;
