import Convert from "ansi-to-html";
import { execSync } from "child_process";
import * as fs from "fs";
import { Code, Parent } from "mdast";
import tmp from "tmp";
import { Transformer } from "unified";
import visit from "unist-util-visit";
import escapeHtml from 'escape-html';

const convert = new Convert({
  fg: "#cccccc",
});

const examplePromptRegex = /^<ExamplePrompt>([^]*)<\/ExamplePrompt>$/;

function execCombined(command: string): string {
  const out = tmp.fileSync({ postfix: ".stdout" });
  execSync(command, { stdio: ["ignore", out.fd, out.fd] });
  return fs.readFileSync(out.name, { encoding: "utf-8" });
}

function runKitschPrompt(exmaple: string): string {
  const configParts = exmaple.split("---");

  let demo: string;
  let config: string;
  if (configParts.length === 2) {
    demo = configParts[0];
    config = configParts[1];
  } else {
    demo = "";
    config = configParts[0];
  }

  // Copy the example to a temporary file
  const demoFile = tmp.fileSync({ postfix: ".yaml" });
  fs.writeFileSync(demoFile.name, demo);

  let options = "";
  if(config.trim()) {
    const configFile = tmp.fileSync({ postfix: ".yaml" });
    fs.writeFileSync(configFile.name, config);
    options += ` --config "${configFile.name}"`;
  }

  // Run kitsch-prompt
  try {
    const output = execCombined(
      `kitsch-prompt prompt ${options} --demo "${demoFile.name}"`
    );


    let html = escapeHtml(output);

    // Replace "{" and "}" because otherwise they will be interpreted as jsx.
    html = html.replace(/([{}])/g, "{'$1'}");

    // Replace " " with "&nbsp;".
    html = html.replace(/ /g, "&nbsp;");

    // Convert ANSI to HTML.
    html = convert.toHtml(html);

    // Fix each style tag to be a JSX style tag.
    html = html.replace(/style="([^"]*)"/g, (match, style) => {
      const reactStyle = style
        .split(";")
        .map((s) => {
          const [key, value] = s.split(":");
          return `"${key}": "${value}"`;
        })
        .join(",");

      return `style={{${reactStyle}}}`;
    });

    html = html.replace(/\n/g, "<br/>");

    return html;
  } catch (err) {
    return "Error running example: " + err.stack;
  }
}

export default function kitschPromptPreview(): Transformer {
  const transformer: Transformer = async (ast) => {
    // visit(ast, "jsx", (node: Literal, index: number, parent: Parent) => {
    //   const match = examplePromptRegex.exec(node.value);
    //   if (!match) {
    //     return;
    //   }

    //   node.value = `<ExamplePrompt>\n${runKitschPrompt(
    //     match[1]
    //   )}</ExamplePrompt>`;
    // });

    // TODO: Add `import` for ExamplePrompt if it's not already there?

    visit(ast, "code", (node: Code, index: number, parent: Parent) => {
      if (node.lang !== "kitsch") {
        return;
      }

      parent.children[index] = {
        type: "jsx",
        value: `<ExamplePrompt>\n${
          runKitschPrompt(node.value) + "▌"
        }</ExamplePrompt>`,
      } as any;
    });
  };

  return transformer;
}
