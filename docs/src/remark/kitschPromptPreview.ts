import Convert from "ansi-to-html";
import { execSync } from "child_process";
import * as fs from "fs";
import { Code, Parent } from "mdast";
import tmp from "tmp";
import { Transformer } from "unified";
import visit from "unist-util-visit";
import escapeHtml from "escape-html";
import path from "path";

const KITSCH = process.env.KITSCH || "kitsch";

const convert = new Convert({
  fg: "#cccccc",
});

const examplePromptRegex = /^<ExamplePrompt>([^]*)<\/ExamplePrompt>$/;

function escapeMdx(text: string): string {
  let result = escapeHtml(text);

  // Replace "{" and "}" because otherwise they will be interpreted as jsx.
  result = result.replace(/([{}])/g, "{'$1'}");

  // Replace " " with "&nbsp;".
  result = result.replace(/ /g, "&nbsp;");

  return result;
}

function execCombined(
  command: string,
  options: { env?: NodeJS.ProcessEnv | undefined }
): string {
  const out = tmp.fileSync({ postfix: ".stdout" });
  execSync(command, { ...options, stdio: ["ignore", out.fd, out.fd] });
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

  // Fix the CWD to be the root of the docs folder.
  config = config.replace(/\${CWD}/g, path.join(__dirname, "..", ".."));

  // Copy the example to a temporary file
  const demoFile = tmp.fileSync({ postfix: ".yaml" });
  fs.writeFileSync(demoFile.name, demo);

  let options = "";
  if (config.trim()) {
    const configFile = tmp.fileSync({ postfix: ".yaml" });
    fs.writeFileSync(configFile.name, config);
    options += ` --config "${configFile.name}"`;
  }

  // Run kitsch
  try {
    console.log(`Running ${KITSCH}...`);
    const output = execCombined(
      `${KITSCH} prompt ${options} --demo "${demoFile.name}"`,
      { env: { ...process.env, FORCE_COLOR: "3" } }
    );

    let html = escapeMdx(output);

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
    return escapeMdx("Error running example: " + err.stack);
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
