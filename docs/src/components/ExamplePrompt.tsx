import React from "react";
import styles from "./ExamplePrompt.module.css";

interface ExamplePromptProps {
  children: React.ReactNode;
}

// Example is used to render
const ExamplePrompt: React.FC<ExamplePromptProps> = (props) => {
  return (
    <div className={styles.container}>
      <div className={styles.top}>
        <span className={styles.dot} style={{background:"#ED594A"}}></span>
        <span className={styles.dot} style={{background:"#FDD800"}}></span>
        <span className={styles.dot} style={{background:"#5AC05A"}}></span>
      </div>
      <div className={styles.content}>{props.children}</div>
    </div>
  );
};

export default React.memo(ExamplePrompt);
